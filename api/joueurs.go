//
// Fichier     : joueurs.go
// Développeur : ?
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'un joueur.
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GetJoueursHandler Gère la récupération de tous les joueurs
func (a *AcquisitionService) GetJoueursHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	users := []Players{}
	db.Find(&users)
	Message(w, users, http.StatusOK)
}

// HandleJoueur Gère la modification et l'ajout d'un joueur
func (a *AcquisitionService) HandleJoueur(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	id := mux.Vars(r)["id"]

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var p Players
	var dat map[string]interface{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	switch r.Method {
	case "POST":
		db.Create(&p)
		num := dat["EquipeID"]
		if num != "" {
			Team := Teams{}
			db.First(&Team, num)
			p.Teams = append(p.Teams, Team)
			db.Model(&Team).Association("Players").Append(p)
		}
		Message(w, p, http.StatusCreated)
	case "PUT":
		db.Model(&p).Where("ID = ?", id).Updates(p)
		Message(w, p, http.StatusOK)
	}
}
