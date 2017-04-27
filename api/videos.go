//
// Fichier     : videos.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de récupérer une vidéo à afficher et analyser
//

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// VideoHandler Gère l'envoie d'une vidéo au client
func (a *AcquisitionService) VideoHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}
	GameID := mux.Vars(r)["id"]
	part := mux.Vars(r)["part"]

	var v Videos
	db.Where("game_id = ? AND Part = ?", GameID, part).First(&v)

	if v.ID != 0 {
		http.ServeFile(w, r, v.Path)
	} else {
		msg := map[string]string{"error": "Fichier inexistant"}
		Message(w, msg, http.StatusNotFound)
	}
}
