package api

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// UploadHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()
		a.ErrorHandler(w, err)

		err = r.ParseMultipartForm(8589934592) // 8Gb
		a.ErrorHandler(w, err)

		formdata := r.MultipartForm

		if _, err := os.Stat("./videos/"); os.IsNotExist(err) {
			os.MkdirAll("./videos/", 0777)
		}

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

				// Ajout de la partie
				v.Part = int(p)
			} else {
				v.Part = 1
			}

			// Le timestamp sera le nom du fichier
			timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			filename := timestamp + "." + ext

			v.Completed = 0
			v.Path = "home/tsap/api/videos/" + filename
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
					a.ErrorHandler(w, err)

					out, err := os.Create("./videos/" + filename)

					defer out.Close()
					a.ErrorHandler(w, err)

					_, err = io.Copy(out, file)

					a.ErrorHandler(w, err)
				}
			} else {
				msg := map[string]string{"error": "Une vidéo avec le même nom existe déjà. Veuillez renommer cette vidéo."}
				Message(w, msg, http.StatusInternalServerError)
				return
			}
		}
		msg := map[string]string{"succes": "Video(s) envoyé(s) avec succès!", "game_id": strconv.Itoa(int(g.ID))}
		Message(w, msg, http.StatusCreated)
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}
