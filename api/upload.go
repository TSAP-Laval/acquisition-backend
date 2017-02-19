package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// UploadHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) UploadHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\nmethod:", r.Method)
	//fmt.Println(r.Header)

	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Print("\nERROR : ")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// On regarde si le dossier videos existe déjà.
	// Dans le cas contraire, on le crée
	if _, err := os.Stat("./videos/"); os.IsNotExist(err) {
		os.MkdirAll("./videos/", 0777)
	}

	f, err := os.OpenFile("./videos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
}
