//
// Fichier     : equipes.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'une équipes.
//

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetEquipeHandler gère la récupération des équipes correspondant au nom entré
func (a *AcquisitionService) GetEquipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}
	switch r.Method {
	case "GET":
		teams := []Teams{}
		name := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Model(&teams).Preload("Coaches").Preload("Players").Where("LOWER(Name) LIKE LOWER(?)", name+"%").Find(&teams)

		for i := range teams {
			db.Model(&teams[i]).Related(&teams[i].Season, "SeasonID")
			db.Model(&teams[i]).Related(&teams[i].Category, "CategoryID")
			db.Model(&teams[i]).Related(&teams[i].Sport, "SportID")
		}

		Message(w, teams, http.StatusOK)
	}
}

// EquipesHandler gère la modification et la suppression des équipes
func (a *AcquisitionService) EquipesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	team := []Teams{}
	id := strings.ToLower(strings.TrimSpace(vars["id"]))
	db.First(&team, "ID = ?", id)

	switch r.Method {
	case "PUT":
		body, err := ioutil.ReadAll(r.Body)
		if len(body) > 0 {
			var t Teams
			err = json.Unmarshal(body, &t)
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			t.Name = strings.TrimSpace(t.Name)
			t.City = strings.TrimSpace(t.City)

			// Omit
			var o string
			if t.Name == "" {
				o += "Name, "
			}
			if t.City == "" {
				o += "City, "
			}
			db.Model(&team).Where("ID = ?", id).Omit(o).Updates(t)

			// L'équipe modifiée (new team)
			var nt Teams
			db.Where("ID = ?", id).Find(&nt)

			nt = AjoutNiveauSport(db, nt)
			db.Model(&nt).Related(&nt.Season, "SeasonID")
			db.Model(&nt).Related(&nt.Category, "CategoryID")
			db.Model(&nt).Related(&nt.Sport, "SportID")

			Message(w, nt, http.StatusCreated)

		} else {
			msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
			Message(w, msg, http.StatusBadRequest)
		}
	case "DELETE":
		// Erreur
		if len(team) == 0 {
			msg := map[string]string{"error": "Aucune equipe ne correspond. Elle doit déjà avoir été supprimée!"}
			Message(w, msg, http.StatusBadRequest)
		} else {
			// On supprime l'équipe
			db.Where("ID = ?", id).Delete(&team)
			Message(w, "", http.StatusNoContent)
		}

	}
}

// GetEquipesHandler gère la récupération de toutes les équipes de la base de donnée
func (a *AcquisitionService) GetEquipesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	teams := []Teams{}
	db.Model(&teams).Preload("Coaches").Preload("Players").Find(&teams)

	for i := range teams {
		db.Model(&teams[i]).Related(&teams[i].Season, "SeasonID")
		db.Model(&teams[i]).Related(&teams[i].Category, "CategoryID")
		db.Model(&teams[i]).Related(&teams[i].Sport, "SportID")
	}

	Message(w, teams, http.StatusOK)
}

// CreerEquipeHandler gère la création d'une équipe dans la base de donnée
func (a *AcquisitionService) CreerEquipeHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		var t Teams
		err = json.Unmarshal(body, &t)
		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		// On enlève les espaces superflues
		t.Name = strings.TrimSpace(t.Name)
		t.City = strings.TrimSpace(t.City)

		if t.Name == "" || t.City == "" {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			Message(w, msg, http.StatusBadRequest)
		} else {

			team := []Teams{}
			name := strings.ToLower(strings.TrimSpace(t.Name))
			db.Where("LOWER(Name) = LOWER(?)", name).Find(&team)

			if len(team) > 0 {
				msg := map[string]string{"error": "Une équipe de même nom existe déjà. Veuillez choisir une autre nom."}
				Message(w, msg, http.StatusBadRequest)
			} else {
				db.Create(&t)
				db.Model(&t).Related(&t.Season, "SeasonID")
				db.Model(&t).Related(&t.Category, "CategoryID")
				db.Model(&t).Related(&t.Sport, "SportID")
				Message(w, t, http.StatusCreated)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez remplir tous les champs."}
		Message(w, msg, http.StatusBadRequest)
	}
}

// Message gère les messages (erreurs, messages de succès) à envoyer au client
func Message(w http.ResponseWriter, msg interface{}, code int) {
	message, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(message)

	return
}

// AjoutNiveauSport permet d'ajouter les informations sur le sport et le niveau lors de l'affchage des infos
func AjoutNiveauSport(db *gorm.DB, t Teams) Teams {
	// Ajout du sport pour l'affichage
	var s Sports
	db.Where("ID = ?", t.SportID).Find(&s)
	if s.Name != "" {
		t.Sport = s
	}

	// Ajout du niveau pour l'affichage
	var c Categories
	db.Where("ID = ?", t.CategoryID).Find(&c)
	if c.Name != "" {
		t.Category = c
	}

	return t
}

// ErrorHandler gère les erreurs côté serveur
func (a *AcquisitionService) ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		a.Error(fmt.Sprintf("ERROR : %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		a.ErrorWrite(err.Error(), w)
	}
	return
}
