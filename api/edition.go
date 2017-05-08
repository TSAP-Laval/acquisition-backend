//
// Fichier     : edition.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonctions...
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// TODO: Changer le nom du fichier edition avec un s
// TODO: Changer le nom du fichier et ses références pour catégorie...
// TODO: Linter le code... Aucun commentaire pour les fonctions
// TODO: Enlever tous ce qui est log, print...

// TODO: Pourquoi tous les endpoints pour la gestion d'un joueur sont dans joueurs.go,
//       mais celui-ci (GET) dans edition.go ?

func (a *AcquisitionService) GetJoueurs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	user := []Players{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) GetActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	user := []ActionsType{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) PostAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var t Actions
	err = json.Unmarshal(body, &t)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}
	switch r.Method {
	case "POST":

	if db.NewRecord(t) {
		db.Create(&t)
		db.NewRecord(t)
		w.Header().Set("Content-Type", "application/text")

	} else {

		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}	

}
