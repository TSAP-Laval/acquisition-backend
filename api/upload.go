package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// UploadHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		fmt.Print("\nERROR : ")
		fmt.Println(err)
		return
	}

	// On regarde si le dossier videos existe déjà.
	// Dans le cas contraire, on le crée
	if _, err := os.Stat("./videos/"); os.IsNotExist(err) {
		os.MkdirAll("./videos/", 0777)
	}

	f, err := os.OpenFile("./videos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		fmt.Print("ERROR : ")
		fmt.Println(err)
		return
	}

	io.Copy(f, file)

	msg := map[string]string{"succes": "Le video a été envoyé avec succès!"}
	succesJSON, _ := json.Marshal(msg)

	w.Header().Set("Content-Type", "application/json")
	w.Write(succesJSON)
	w.WriteHeader(201)
}
