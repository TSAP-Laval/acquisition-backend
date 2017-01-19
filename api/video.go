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

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("../video/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}
