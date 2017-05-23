//
// Fichier     : upload.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à l'upload,
// d'une vidéo sur le serveur et l'ajout des informations de base
// sur la partie à créer en lien avec la/les vidéo(s) `uploadé`.
//

package api

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Type dates permettant d'obtenir l'ordre des vidéos selon leur date de création
type dates struct {
	index int
	date  time.Time
}

// Constantes sur les tailles possibles
const (
	GB = 1 << (10 * 3)
	MB = 1 << (10 * 2)
	KO = 1 << (10 * 1)
)

// Le dossier de destination des vidéos uploadées
const videoPath = "../videos/"

// Type de données permettant d'obtenir l'ordre des vidéos selon leur date de création
type creationDates []dates

// UploadHandler Gère l'upload de videos sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		gameID := mux.Vars(r)["game-id"]
		if id, err := strconv.Atoi(gameID); id > 0 || err == nil {
			db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
			defer db.Close()
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			// Taille limite d'envoie de fichiers
			r.Body = http.MaxBytesReader(w, r.Body, 10*GB) // 10 Gb

			// Dans le cas où le dossier de destination des vidéos est inexistant, on le crée
			if _, err := os.Stat(videoPath); os.IsNotExist(err) {
				os.MkdirAll(videoPath, 0777)
			}

			var form *multipart.Reader
			if form, err = r.MultipartReader(); err != nil {
				msg := map[string]string{"error": "Aucun fichier n'a été envoyé ! Veuillez réessayer !"}
				Message(w, msg, http.StatusBadRequest)
				return
			}

			// Tableau des vidéos uploadées servant à l'insertion des informations dans la base de données
			var videos = make([]*Videos, 0)

			// variable servant à la récupération de l'itération actuelle de la
			// vidéo uploadée et de se date de création
			var partsDate = make(map[string]dates)

			// Boucle sur les "parties" du formulaire envoyée
			//  ** Boucle infinie, mais termine une fois que la dernière partie est lue
			for i := 0; i >= 0; i++ {
				part, err := form.NextPart()
				if err == io.EOF {
					break
				}

				// Convertion de la date de création de la vidéo selon le format donné
				layout := "Mon Jan 02 2006 15:04:05 GMT-0400 (EST)"
				date := part.FormName()

				t, err := time.Parse(layout, date)
				if err != nil {
					layout := "Mon Jan 02 2006 15:04:05 GMT-0400 (EDT)"
					t, err = time.Parse(layout, date)

					if err != nil {
						msg := map[string]string{"error": "Une erreur est survenue lors de l'envoie de la vidéo!"}
						Message(w, msg, http.StatusBadRequest)
						return
					}
				}

				// Forn name is use as file last modified date
				partsDate[strconv.Itoa(i)] = dates{index: i, date: t}

				// On crée un buffer qui contiendra l'en-tête du fichier
				//
				// ** Cela permettra de déterminer le type du fichier.
				//    Ainsi, on valide que le fichier est bel et bien
				//    un fichier au format vidéo/* et non un fichier
				//    quelconque renommé en .mp4
				buffer := make([]byte, 512)

				var cBytes int
				if cBytes, err := part.Read(buffer); err != nil || cBytes == 0 {
					msg := map[string]string{"error": "Le fichier envoyé ne possède aucun contenu!"}
					Message(w, msg, http.StatusBadRequest)
					return
				}

				var contentType string
				var ext []string

				// On valide préalablement que le fichier est une vidéo au format
				// Quicktime (.mov) ni MPEG (.mpg/mpeg). Dans le cas où ell ne l'est
				// pas, on valide alors son format avec la fonction native `DetectContentType`.
				//
				// ** Cette fonction ne continent pas la définition pour les
				//    fichier .mov et .mpg/mpeg, c'est pourquoi j'ai ajouté la fonction
				//    qui permet de valider ce type de fichier.
				if !isMov(buffer, &contentType, &ext) && !isMpg(buffer, &contentType, &ext) {
					// Utilisation de la fonction native permettant de déterminer le format du fichier
					contentType = http.DetectContentType(buffer)

					// Seule une erreur dans la REGEX peut causer le retour d'une erreur par la fonction
					// c'est pourquoi, en m'assurant qu'elle est conforme, j'ai pris la décision de retirer
					// la gestion de celle-ci.
					validation, _ := regexp.Compile("video/.*")

					if !validation.Match([]byte(contentType)) {
						msg := map[string]string{"error": "Le fichier \"" + part.FileName() + "\" n'est pas une vidéo de format valide ! Les format supportés sont : mp4, avi, mov, mpeg."}
						Message(w, msg, http.StatusBadRequest)
						return
					}

					// Récupération de l'extention du format de fichier
					ext, _ = mime.ExtensionsByType(contentType)
				}

				// Le timestamp sera le nom du fichier
				timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

				filename := timestamp + ext[0]

				dest, _ := os.OpenFile(videoPath+filename, os.O_WRONLY|os.O_CREATE, 0777)
				defer dest.Close()

				// Écriture de l'en-tête venant d'être lue
				dest.Write(buffer)

				var v Videos
				// Par défaut, la vidéo est la première partie.
				// De cette façon, dans le cas où il n'y a qu'une vidéo
				// uploadé, elle est la seule partie
				v.Part = 1

				v.Completed = 0
				v.Path, err = filepath.Abs("../videos/" + filename)

				v.GameID, _ = strconv.Atoi(gameID)

				// Création de la vidéo dans la base de données
				if db.NewRecord(v) {
					db.Create(&v)
					videos = append(videos, &v)
				}

				// Upload de la vidéo
				//
				// ** La boucle termine une fois que
				//    cette partie du formulaire a été
				//    envoyé entièrement au serveur
				for {
					buffer = make([]byte, 80*KO)

					cBytes, err = part.Read(buffer)

					if cBytes != 0 {
						dest.Write(buffer[0:cBytes])
					}

					if err == io.EOF {
						break
					}
				}
			}

			// Seulement dans le cas où il y a plus d'une partie uploadée
			if len(videos) > 1 {
				// Tri les dates en ordre croissant pour permettre de récupérer
				// l'ordre des vidéos en plusieurs parties
				creationDateSorted := make(creationDates, 0, len(partsDate))
				for _, d := range partsDate {
					creationDateSorted = append(creationDateSorted, d)
				}
				sort.Sort(creationDateSorted)

				// Ajout du noméro de la partie de la vidéo dans la base de données
				for j, video := range videos {
					index := creationDateSorted.IndexOf(j)
					video.Part = index + 1
					db.Model(&video).Where("ID = ?", video.ID).Update("part", video.Part)
				}
			}

			msg := map[string]string{"succes": "Video(s) envoyé(s) avec succès!"}
			Message(w, msg, http.StatusCreated)
		} else {
			Message(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	case "DELETE":
		var g Games
		gameID := mux.Vars(r)["game-id"]

		// Erreur, l'identifiant d'une partie ne peut être de 0
		if id, err := strconv.Atoi(gameID); id <= 0 || err != nil {
			msg := map[string]string{"error": "Impossible de mettre fin à la partie, car aucune partie ne correspond. Elle doit déjà avoir été supprimée!"}
			Message(w, msg, http.StatusNotFound)
		} else {
			db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
			defer db.Close()
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			v := []Videos{}
			// Récupération des vidéos à supprimer
			db.Where("game_id = ?", gameID).Find(&v)

			// Pour chacune des viséos contenues dans la partie
			for _, video := range v {
				// On supprime la vidéo du serveur
				if err = os.Remove(video.Path); err != nil {
					msg := map[string]string{"error": "Impossible de supprimer la video ! Elle doit déjà avoir été supprimée!"}
					Message(w, msg, http.StatusNotFound)
					return
				}
			}

			// On supprime la/les vidéo(s) de la base de donnée
			db.Where("game_id = ?", gameID).Delete(&Videos{})
			// On supprime la partie de la base de donnée
			db.Where("ID = ?", gameID).Delete(&g)

			Message(w, "", http.StatusNoContent)
		}
	}
}

// Retourne le nombre de dates contenus dans le tableau
func (d creationDates) Len() int {
	return len(d)
}

// Less fonction utilisée pour trier les dates
func (d creationDates) Less(i int, j int) bool {
	return d[i].date.Before(d[j].date)
}

// Swap fonction utilisée pour trier les dates
func (d creationDates) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// IndexOf fonction utilisée pour trier les dates
func (d creationDates) IndexOf(value int) int {
	var index int
	for i, date := range d {
		if date.index == value {
			index = i
		}
	}
	return index
}

// isMov permet de déterminer si le fichier
// envoyé vers le serveur est bel est bien
// dans un format pris en charge.
//
// ** Dans ce cas ci, on n'utilise pas la
//    fonction DetectContentType, car elle
//    ne valide pas les vidéos au format mov
func isMov(buf []byte, contentType *string, ext *[]string) bool {
	isMov := len(buf) > 7 &&
		(buf[0] == 0x0 && buf[1] == 0x0 &&
			buf[2] == 0x0 && buf[3] == 0x14 &&
			buf[4] == 0x66 && buf[5] == 0x74 &&
			buf[6] == 0x79 && buf[7] == 0x70) ||
		(buf[0] == 0x6D && buf[1] == 0x6F &&
			buf[2] == 0x6F && buf[3] == 0x76)

	if isMov {
		*contentType = "video/quicktime"
		*ext = []string{".mov"}
	}
	return isMov
}

// isMpg permet de déterminer si le fichier
// envoyé vers le serveur est bel est bien
// dans un format pris en charge.
//
// ** Dans ce cas ci, on n'utilise pas la
//    fonction DetectContentType, car elle
//    ne valide pas les vidéos au format mpg
func isMpg(buf []byte, contentType *string, ext *[]string) bool {
	isMpg := len(buf) > 4 &&
		buf[0] == 0x0 && buf[1] == 0x0 &&
		buf[2] == 0x1 && buf[3] == 0xBA

	if isMpg {
		*contentType = "video/mpeg"
		*ext = []string{".mpg"}
	}
	return isMpg
}
