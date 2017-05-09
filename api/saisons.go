//
// Fichier     : saisons.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonction...
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	// Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// TODO: Linter le code... Aucun commentaire pour les fonctions
// TODO: Enlever tous ce qui est log, print...

func (a *AcquisitionService) PostSaison(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var t Seasons
	err = json.Unmarshal(body, &t)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	if db.NewRecord(t) {
		db.Create(&t)
		if db.NewRecord(t) {
			msg := map[string]string{"error": "Une erreur est survenue lors de la création de la saison. Veuillez réessayer!"}
			Message(w, msg, http.StatusInternalServerError)
		} else {
			Message(w, t, http.StatusCreated)
		}

	} else {
		msg := map[string]string{"error": "Une erreur est survenue lors de la création de la saison. Veuillez réessayer!"}
		Message(w, msg, http.StatusInternalServerError)
	}

}

func (a *AcquisitionService) GetSeasons(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	s := []Seasons{}
	db.Find(&s)

	Message(w, s, http.StatusOK)
}
