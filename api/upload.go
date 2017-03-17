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

		err = r.ParseMultipartForm(2000000000) // grab the multipart form
		a.ErrorHandler(w, err)

		formdata := r.MultipartForm // ok, no problem so far, read the Form data

		if _, err := os.Stat("./videos/"); os.IsNotExist(err) {
			os.MkdirAll("./videos/", 0777)
		}

		//get the *fileheaders
		files := formdata.File["file"] // grab the filenames

		var g Games
		db.Create(&g)
		for i := range files { // loop through the files one by one
			fileSplit := strings.Split(files[i].Filename, ".")
			ext := fileSplit[len(fileSplit)-1]

			re := regexp.MustCompile("^.*?([^d]*(d+)[^d]*).*$")
			part := re.FindString(files[i].Filename)
			part = strings.Replace(part, "(", " ", 1)
			part = strings.Replace(part, ")", " ", 1)

			timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			filename := timestamp + "." + ext

			var v Videos
			v.Completed = false
			v.Path = "home/tsap/api/videos/" + filename
			p, _ := strconv.ParseInt(part, 64, 10)
			v.Part = int(p)
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
