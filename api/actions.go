//
// Fichier     : actions.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonctions...
//

package api

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"

	"io/ioutil"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// TODO: Linter le code...
// TODO: Gérer les erreurs comme du monde
// TODO: Enlever tous ce qui est log, print...

// GetMovementTypeHandler Gestion du select des types de mouvements
func (a *AcquisitionService) GetMovementTypeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	mvmType := []MovementsType{}
	db.Find(&mvmType)

	mvmTypeJSON, _ := json.Marshal(mvmType)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(mvmTypeJSON)
}

//GetAllActionsTypes gestion du select des types d'actions
func (a *AcquisitionService) GetAllActionsTypes(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	actionTypes := []ActionsType{}
	db.Find(&actionTypes)

	actionTypesJSON, _ := json.Marshal(actionTypes)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(actionTypesJSON)

	defer db.Close()
}

//PostActionType : Create new action type
func (a *AcquisitionService) PostActionType(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var newActionType ActionsType

	err = json.Unmarshal(body, &newActionType)

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	if db.NewRecord(newActionType) {
		db.Create(&newActionType)
		db.NewRecord(newActionType) // => return `false` after `user` created
	} else {
		// TODO: Gérer l'erreur
		return
	}

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
