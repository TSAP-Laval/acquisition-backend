package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type dates struct {
	index int
	date  time.Time
}

// Size constants
const (
	GB = 1 << (10 * 3)
	MB = 1 << (10 * 2)
	KO = 1 << (10 * 1)
)

const videoPath = "../videos/"

type creationDates []dates

// UploadHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		gameID := mux.Vars(r)["game-id"]
		if _, err := strconv.Atoi(gameID); err == nil {
			db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
			defer db.Close()
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			// 10 Gb
			if r.ContentLength > 10*GB {
				msg := map[string]string{"error": "Fichier de trop grande taille. La taille maximal pour un fichier est de 10Gb"}
				Message(w, msg, http.StatusBadRequest)
				return
			}

			// Limit upload size
			r.Body = http.MaxBytesReader(w, r.Body, 10*GB) // 10 Gb

			if _, err := os.Stat(videoPath); os.IsNotExist(err) {
				os.MkdirAll(videoPath, 0777)
			}

			var form *multipart.Reader
			if form, err = r.MultipartReader(); err != nil {
				msg := map[string]string{"error": "Aucun fichier envoyé ! Veuillez réessayer !"}
				Message(w, msg, http.StatusBadRequest)
				return
			}

			var videos = make([]*Videos, 0)
			var partsDate = make(map[string]dates)
			// Infinite loop, but will stop when last part is read
			for i := 0; i >= 0; i++ {
				part, err := form.NextPart()
				if err == io.EOF {
					break
				}

				layout := "Mon Jan 02 2006 15:04:05 GMT-0400 (EDT)"
				str := part.FormName()
				t, _ := time.Parse(layout, str)

				// Forn name is use as file last modified date
				partsDate[strconv.Itoa(i)] = dates{index: i, date: t}

				// On crée un buffer qui contiendra l'en-tête du fichier
				// ** Cela permettra de déterminer le type du fichier.
				//    Ainsi, on valide que le fichier est bel et bien
				//    un fichier au format vidéo/* et non un fichier
				//    quelconque renommé en .mp4
				buffer := make([]byte, 512)

				var cBytes int
				if cBytes, err := part.Read(buffer); err != nil || cBytes == 0 {
					return
				}

				// On valide préalablement que le fichier est une vidéo au format
				// Quicktime (.mov). Dans le cas où ell ne l'est pas, on valide
				// alors son format avec la fonction native `DetectContentType`.
				// ** Cette fonction ne continent pas la définition pour les
				//    fichier .mov, c'est pourquoi j'ai ajouté la fonction
				//    qui permet de valider ce type de fichier.
				if !isMov(buffer) {
					contentType := http.DetectContentType(buffer)
					var validation *regexp.Regexp
					if validation, err = regexp.Compile("video/.*"); err != nil {
						a.ErrorHandler(w, err)
						return
					}

					if !validation.Match([]byte(contentType)) {
						msg := map[string]string{"error": "Le fichier \"" + part.FileName() + "\" n'est pas une vidéo de format valide ! Les format supportés sont : mp4, avi, mov, mpeg."}
						Message(w, msg, http.StatusBadRequest)
						return
					}
				}

				fileSplit := strings.Split(part.FileName(), ".")
				ext := fileSplit[len(fileSplit)-1]

				// Le timestamp sera le nom du fichier
				timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

				filename := timestamp + "." + ext

				var dst *os.File
				if dst, err = os.OpenFile(videoPath+filename, os.O_WRONLY|os.O_CREATE, 0777); err != nil {
					msg := map[string]string{"error": "Une erreur inconnue est survenue lors de l'écriture du fichier \"" + part.FileName() + "\". Veuillez réessayer !"}
					Message(w, msg, http.StatusBadRequest)
					return
				}
				defer dst.Close()

				// Écriture de l'en-tête venant d'être lue
				dst.Write(buffer)

				var v Videos
				v.Part = 1

				v.Completed = 0
				v.Path, err = filepath.Abs("../videos/" + filename)

				v.GameID, _ = strconv.Atoi(gameID)

				if db.NewRecord(v) {
					db.Create(&v)
					if db.NewRecord(v) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création de la video dans la base de données. Veuillez réessayer!"}
						Message(w, msg, http.StatusInternalServerError)
					}
					videos = append(videos, &v)
				}

				for {
					buffer = make([]byte, 4*KO)

					cBytes, err = part.Read(buffer)

					if cBytes != 0 {
						dst.Write(buffer[0:cBytes])
					} else {
						break
					}

					if err == io.EOF {
						break
					}
				}
			}

			//Sort the map by date
			creationDateSorted := make(creationDates, 0, len(partsDate))
			for _, d := range partsDate {
				creationDateSorted = append(creationDateSorted, d)
			}
			sort.Sort(creationDateSorted)

			for j, video := range videos {
				var index int
				if index = creationDateSorted.IndexOf(j); index == -1 {
					return
				}
				video.Part = index + 1
				db.Model(&video).Where("ID = ?", video.ID).Update("part", video.Part)
			}

			msg := map[string]string{"succes": "Video(s) envoyé(s) avec succès!"}
			Message(w, msg, http.StatusCreated)
		} else {
			msg := map[string]string{"error": "Une erreur est survenue lors de la création de la video dans la base de données. Veuillez réessayer!"}
			Message(w, msg, http.StatusInternalServerError)
		}
	case "DELETE":
		var g Games
		gameID := mux.Vars(r)["game-id"]
		fmt.Printf("\nGAMEID : %s", gameID)
		// Erreur, l'identifiant d'une partie ne peut être de 0
		if id, err := strconv.Atoi(gameID); id <= 0 || err != nil {
			fmt.Println(err)
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
			// On supprime les vidéos
			db.Where("game_id = ?", gameID).Find(&v)
			fmt.Println(v[0].Path)
			for _, video := range v {
				// delete file
				if err = os.Remove(video.Path); err != nil {
					msg := map[string]string{"error": "Impossible de supprimer la video ! Elle doit déjà avoir été supprimée!"}
					Message(w, msg, http.StatusNotFound)
					return
				}
			}

			// On supprime les vidéos
			db.Where("game_id = ?", gameID).Delete(&Videos{})
			// On supprime la partie
			db.Where("ID = ?", gameID).Delete(&g)

			msg := map[string]string{"succes": "L'équipe et les vidéos ont été supprimée avec succès!"}
			Message(w, msg, http.StatusNoContent)
		}
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}

func (d creationDates) Len() int {
	return len(d)
}

func (d creationDates) Less(i int, j int) bool {
	return d[i].date.Before(d[j].date)
}

func (d creationDates) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d creationDates) IndexOf(value int) int {
	for i, date := range d {
		if date.index == value {
			return i
		}
	}
	return -1
}

func isMov(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0x0 && buf[1] == 0x0 &&
		buf[2] == 0x0 && buf[3] == 0x14 &&
		buf[4] == 0x66 && buf[5] == 0x74 &&
		buf[6] == 0x79 && buf[7] == 0x70
}
