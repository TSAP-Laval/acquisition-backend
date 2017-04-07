package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Size constants
const (
	GB = 1 << (10 * 3)
	MB = 1 << (10 * 2)
	KO = 1 << (10 * 1)
)

const videoPath = "../videos/"

// UploadHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()
		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		// 10 GB
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

		var g Games

		for {
			var part *multipart.Part
			if part, err = form.NextPart(); err == io.EOF {
				break
			}
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

			contentType := http.DetectContentType(buffer)
			var validation *regexp.Regexp
			if validation, err = regexp.Compile("video/.*"); err != nil {
				a.ErrorHandler(w, err)
				return
			}

			if !validation.Match([]byte(contentType)) {
				msg := map[string]string{"error": "Le fichier \"" + part.FileName() + "\" n'est pas une vidéo de format valide ! Les format supportés sont : mp4, avi, mov."}
				Message(w, msg, http.StatusBadRequest)
				return
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

			for {
				buffer = make([]byte, 4*KO)

				cBytes, err = part.Read(buffer)

				if cBytes != 0 {
					dst.Write(buffer[0:cBytes])
				} else {
					break
				}

				if err == io.EOF {
					var v Videos
					v.Part = 1

					v.Completed = 0
					v.Path, err = filepath.Abs("../videos/" + filename)
					db.Create(&g)
					v.Game = g

					if db.NewRecord(v) {
						db.Create(&v)
						if db.NewRecord(v) {
							// Dans le cas où il y a une erreur, on supprime la partie
							// venant d'être créée
							db.Delete(&g)
							msg := map[string]string{"error": "Une erreur est survenue lors de la création de la video dans la base de données. Veuillez réessayer!"}
							Message(w, msg, http.StatusInternalServerError)
						}
					}
					break
				}
			}
		}
		msg := map[string]string{"succes": "Video(s) envoyé(s) avec succès!", "game_id": strconv.Itoa(int(g.ID))}
		Message(w, msg, http.StatusCreated)
	case "DELETE":
		var g Games
		gameID := mux.Vars(r)["game-id"]
		// Erreur, l'identifiant d'une partie ne peut être de 0
		if id, err := strconv.Atoi(gameID); id <= 0 || err != nil {
			fmt.Printf("ID : %d", id)
			fmt.Print(err)
			msg := map[string]string{"error": "Aucune partie ne correspond. Elle doit déjà avoir été supprimée!"}
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
