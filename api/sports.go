//
// Fichier     : sports.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonctions...
//

package api

import (
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// TODO: Linter le code... Aucun commentaire pour les fonctions
// TODO: Enlever tous ce qui est log, print...

func (a *AcquisitionService) GetSports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
