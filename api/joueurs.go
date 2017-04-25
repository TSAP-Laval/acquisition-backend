//
// Fichier     : niveaux.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonction...
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// TODO: Rendre le code beau et épuré!
// TODO: Linter le code... Aucun commentaire pour les fonctions
// TODO: Enlever tous ce qui est log, print...

// HandleJoueur gère la modification et l'ajout de joueur
func (a *AcquisitionService) HandleJoueur(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var t Players
	var dat map[string]interface{}
	// TODO: GÉRER LES ERREURS !?!?!?!?!?
	err = json.Unmarshal(body, &t)
	err = json.Unmarshal(body, &dat)
	switch r.Method {
	case "POST":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if db.NewRecord(t) {
			db.Create(&t)
			db.NewRecord(t)
			num := dat["EquipeID"]
			if num != "" {
				Team := Teams{}
				db.First(&Team, num)
				if err != nil {
					a.ErrorHandler(w, err)
					return
				}
				t.Teams = append(t.Teams, Team)
				db.Model(&Team).Association("Players").Append(t)
			}
			Message(w, t, http.StatusCreated)

		} else {
			// Très beau message !
			Message(w, "déjà créé", http.StatusBadRequest)
		}
	case "PUT":
		// Ca donne rien de mettre cette ligne là parce que le middleware le fait déjà !
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.Model(&t).Where("ID = ?", id).Updates(t)
		// Encore mieux !
		Message(w, t, http.StatusOK)
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}
