//
// Fichier     : niveaux.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonction...
//

package api

import (
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// TODO: Changer le nom du fichier et ses références pour catégorie...
// TODO: Linter le code... Aucun commentaire pour les fonctions
// TODO: Enlever tous ce qui est log, print...

func (a *AcquisitionService) GetNiveau(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
