//
// Fichier     : actions.go
// Développeur : Laurent Leclerc-Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'un type d'action.
//

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"io/ioutil"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//GetActionsTypeHandler Gère la récupération de tous les types d'actions
func (a *AcquisitionService) GetActionsTypeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	id := mux.Vars(r)["id"]

	if id != "" {
		acType := ActionsType{}
		db.Where("ID = ?", id).First(&acType)

		if acType.ID != 0 {
			Message(w, acType, http.StatusOK)
		} else {
			msg := map[string]string{"error": "Aucun type d'action ne correspond à celui entré"}
			Message(w, msg, http.StatusNotFound)
		}
	} else {
		acType := []ActionsType{}
		db.Find(&acType)

		Message(w, acType, http.StatusOK)
	}
}

// GetAllReceptionTypes gestion du select pour les types de reception
func (a *AcquisitionService) GetAllReceptionTypes(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	receptionType := []ReceptionType{}
	db.Find(&receptionType)
	Message(w, receptionType, http.StatusOK)
}

//CreerActionsType Gère la création d'un type d'action
func (a *AcquisitionService) CreerActionsType(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var acType ActionsType
	if err = json.Unmarshal(body, &acType); err != nil {
		msg := map[string]string{"error": "Certaines informations entrées sont invalides!"}
		Message(w, msg, http.StatusBadRequest)
		return
	}

	var at ActionsType
	db.Model(&at).Where("name = ?", acType.Name).Find(&at)

	if at.ID == 0 {
		if db.NewRecord(acType) {
			db.Create(&acType)
			Message(w, acType, http.StatusCreated)
		}
	} else {
		msg := map[string]string{"error": "Un type d'action avec le même nom existe déjà"}
		Message(w, msg, http.StatusBadRequest)
	}
}
