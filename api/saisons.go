//
// Fichier     : saisons.go
// Développeur : ?
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'une saison.
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

// CreerSaisonHandler Gère la création d'une nouvelle saison
func (a *AcquisitionService) CreerSaisonHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var s Seasons
	err = json.Unmarshal(body, &s)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var dbSaison Seasons
	db.Model(&dbSaison).Where("Years = ?", s.Years).First(&dbSaison)

	if dbSaison.ID == 0 {
		db.Create(&s)
		Message(w, s, http.StatusCreated)
	} else {
		msg := map[string]string{"error": "La saison entrée existe déjà !"}
		Message(w, msg, http.StatusBadRequest)
	}
}

// GetSeasonsHandler Gère la récupération de toutes les saisons
func (a *AcquisitionService) GetSeasonsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	s := []Seasons{}
	db.Find(&s)

	Message(w, s, http.StatusOK)
}
