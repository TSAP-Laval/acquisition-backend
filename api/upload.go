package api

import (
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
)

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

		// Limit upload size
		r.Body = http.MaxBytesReader(w, r.Body, 3*GB) // 3 Gb

		var form *multipart.Reader
		if form, err = r.MultipartReader(); err != nil {
			a.ErrorHandler(w, err)
			return
		}

		var part *multipart.Part
		if part, err = form.NextPart(); err != nil {
			a.ErrorHandler(w, err)
			return
		}
		// Create a buffer to store the header of the file in
		fileHeader := make([]byte, 512)
		if _, err := part.Read(fileHeader); err != nil {
			return
		}
		contentType := http.DetectContentType(fileHeader)
		var validation *regexp.Regexp
		if validation, err = regexp.Compile("video/.*"); err != nil {
			a.ErrorHandler(w, err)
			return
		}
		if !validation.Match([]byte(contentType)) {
			msg := map[string]string{"error": "Le fichier n'est pas une vidéo de format valide ! Les format supportés sont : mp4, avi, mov."}
			Message(w, msg, http.StatusBadRequest)
			return
		}

		// Taille max de 3Gb pour le fichier
		if err = r.ParseMultipartForm(3 * GB); err != nil {
			a.ErrorHandler(w, err)
			return
		}

		if _, err := os.Stat("../videos/"); os.IsNotExist(err) {
			os.MkdirAll("../videos/", 0777)
		}

		formdata := r.MultipartForm

		files := formdata.File["file"]

		var g Games
		db.Create(&g)

		for i := range files {
			fileSplit := strings.Split(files[i].Filename, ".")
			ext := fileSplit[len(fileSplit)-1]

			var v Videos
			if len(files) > 1 {
				re := regexp.MustCompile(`\((.*?)\)`)
				part := re.FindStringSubmatch(files[i].Filename)[1]
				p, err := strconv.ParseInt(part, 10, 0)

				if err != nil {
					msg := map[string]string{"error": "La nomenlature des fichiers est incorrecte! Veuillez vous assurer qu'ils contiennent un (#)!"}
					a.Error(msg["error"])
					Message(w, msg, http.StatusInternalServerError)
				}

				// Ajout de la partie de la vidéo
				v.Part = int(p)
			} else {
				v.Part = 1
			}

			// Le timestamp sera le nom du fichier
			timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			filename := timestamp + "." + ext

			v.Completed = 0
			v.Path, err = filepath.Abs("../videos/" + filename)
			v.Game = g

			if db.NewRecord(v) {
				db.Create(&v)
				if db.NewRecord(v) {
					msg := map[string]string{"error": "Une erreur est survenue lors de la création de la video dans la base de données. Veuillez réessayer!"}
					a.Error(msg["error"])
					Message(w, msg, http.StatusInternalServerError)
				} else {
					// On regarde si le dossier videos existe déjà.
					// Dans le cas contraire, on le crée
					file, err := files[i].Open()
					defer file.Close()
					if err != nil {
						a.ErrorHandler(w, err)
						return
					}

					out, err := os.Create("../videos/" + filename)

					defer out.Close()
					if err != nil {
						a.ErrorHandler(w, err)
						return
					}

					if written, err := io.Copy(out, file); err != nil || written == 0 {
						a.ErrorHandler(w, err)
						return
					}
				}
			} else {
				msg := map[string]string{"error": "Une vidéo avec le même nom existe déjà. Veuillez renommer cette vidéo.", "exist": "true"}
				Message(w, msg, http.StatusInternalServerError)
				return
			}
		}
		msg := map[string]string{"succes": "Video(s) envoyé(s) avec succès!", "game_id": strconv.Itoa(int(g.ID))}
		Message(w, msg, http.StatusCreated)
	case "DELETE":
		var g Games
		gameID := mux.Vars(r)["game-id"]
		// Erreur
		if gameID == "0" {
			msg := map[string]string{"error": "Aucune partie ne correspond. Elle doit déjà avoir été supprimée!"}
			Message(w, msg, http.StatusNoContent)
		} else {
			db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
			defer db.Close()
			if err != nil {
				a.ErrorHandler(w, err)
				return
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
