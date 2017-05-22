//
// Fichier     : niveaux.go
// Développeur : ?
//
// Permet de gérer toutes les interractions nécessaires à la
// récupération des informations d'un niveau.
//

package api

import (
	"net/http"

	"github.com/jinzhu/gorm"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GetNiveauHandler Gère la récupération des niveaux
func (a *AcquisitionService) GetNiveauHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	c := []Categories{}
	db.Find(&c)

	Message(w, c, http.StatusOK)
}
