//
// Fichier     : edition.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'une action.
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// ActionsPartieHandler Gère la récupération des actions d'une partie
func (a *AcquisitionService) ActionsPartieHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["id"]

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	g := Games{}
	db.First(&g, "ID = ?", gameID)

	if g.ID == 0 {
		msg := map[string]string{"error": "Aucune partie ne correspond."}
		Message(w, msg, http.StatusBadRequest)
	} else {
		var acts []Actions
		db.Model(&acts).Where("game_id = ?", gameID).Find(&acts)
		Message(w, acts, http.StatusOK)
	}
}

// CreerActionHandler Gère la création d'une nouvelle action
func (a *AcquisitionService) CreerActionHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var ac Actions
	err = json.Unmarshal(body, &ac)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var dbAction Actions
	db.Model(&dbAction).Where("time = ? AND game_id = ? ", ac.Time, ac.GameID).First(&dbAction)

	if dbAction.ID == 0 {
		db.Create(&ac)
		Message(w, ac, http.StatusCreated)
	} else {
		msg := map[string]string{"error": "Une action existe déjà à ce moment précis de la partie !"}
		Message(w, msg, http.StatusBadRequest)
	}

}

// SupprimerActionHandler Gère la suppression d'une action
func (a *AcquisitionService) SupprimerActionHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	id := mux.Vars(r)["id"]

	var ac Actions
	db.Model(&ac).Where("ID = ?", id).First(&ac)

	if ac.ID == 0 {
		msg := map[string]string{"error": "Aucune action ne correspond. Elle doit déjà avoir été supprimée!"}
		Message(w, msg, http.StatusNotFound)
	} else {
		db.Where("ID = ?", ac.ID).Delete(&ac)
		Message(w, "", http.StatusNoContent)
	}

}
