//
// Fichier     : sports.go
// Développeur : ?
//
// Permet de gérer toutes les interractions nécessaires à la
// récupération des informations d'un sport.
//

package api

import (
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// GetSports Gère la récupération de tous les sports
func (a *AcquisitionService) GetSports(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	s := []Sports{}
	db.Find(&s)

	Message(w, s, http.StatusOK)
}
