package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// VideoHandler Gère l'upload de video sur le serveur
func (a *AcquisitionService) VideoHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\nmethod:", r.Method)
	//fmt.Println(r.Header)

	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Print("\nERROR :")
		fmt.Println(err)
		return
	}
	defer file.Close()

	//fmt.Fprintf(w, "\nOK %v", handler.Header)

	if _, err := os.Stat("./video/"); os.IsNotExist(err) {
		os.MkdirAll("./video/", 0777)
	}

	f, err := os.OpenFile("./video/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
}
