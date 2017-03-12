package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	"io/ioutil"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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
	fmt.Println(string(mvmTypeJSON))

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
	fmt.Println(string(actionTypesJSON))

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
	fmt.Printf("-----------------------")
	fmt.Println(body)
	fmt.Printf("-----------------------")
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	fmt.Println(string(body))

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
		fmt.Println("erreur")
	}

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
