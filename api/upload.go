package api

import (
	"io"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
)

// UploadHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	defer file.Close()

	a.ErrorHandler(w, err)

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	a.ErrorHandler(w, err)

	var v Videos
	v.Completed = false
	v.Path = "home/tsap/api/videos/" + handler.Filename

	if db.NewRecord(v) {
		db.Create(&v)
		if db.NewRecord(v) {
			msg := map[string]string{"error": "Une erreur est survenue lors de la création de la video dans la base de données. Veuillez réessayer!"}
			Message(w, msg, http.StatusInternalServerError)
		} else {

			// On regarde si le dossier videos existe déjà.
			// Dans le cas contraire, on le crée
			if _, err := os.Stat("./videos/"); os.IsNotExist(err) {
				os.MkdirAll("./videos/", 0777)
			}

			f, err := os.OpenFile("./videos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()

			a.ErrorHandler(w, err)

			io.Copy(f, file)

			msg := map[string]string{"succes": "Le video a été envoyé avec succès!"}
			Message(w, msg, http.StatusCreated)
		}
	} else {
		msg := map[string]string{"error": "Une vidéo avec le même nom existe déjà. Veuillez renommer cette vidéo."}
		Message(w, msg, http.StatusInternalServerError)
		return
	}
}
