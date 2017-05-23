//
// Fichier     : terrains.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'un terrain.
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetTerrainHandler Gère la récupération des terrains correspondant au nom entré
func (a *AcquisitionService) GetTerrainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	locations := []Locations{}
	name := strings.ToLower(strings.TrimSpace(vars["nom"]))
	db.Where("LOWER(Name) LIKE LOWER(?)", name+"%").Find(&locations)

	locations = addFieldTypes(db, locations)

	Message(w, locations, http.StatusOK)
}

// TerrainsHandler gère la modification et la suppression d'un terrains
func (a *AcquisitionService) TerrainsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var location Locations
	id := strings.ToLower(strings.TrimSpace(vars["id"]))
	db.First(&location, "ID = ?", id)

	switch r.Method {
	case "PUT":
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if len(body) > 0 {
			var l Locations
			err = json.Unmarshal(body, &l)
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			l.Name = strings.TrimSpace(l.Name)
			l.City = strings.TrimSpace(l.City)
			l.Address = strings.TrimSpace(l.Address)

			var o string

			if l.Name == "" {
				o += "Name, "
			}
			if l.City == "" {
				o += "City, "
			}
			if l.Address == "" {
				o += "Address, "
			}

			db.Model(&location).Where("ID = ?", id).Omit(o).Updates(l)

			// Le lieu modifié
			var nl Locations
			db.Where("ID = ?", id).Find(&nl)
			Message(w, nl, http.StatusCreated)
		} else {
			msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
			Message(w, msg, http.StatusBadRequest)
		}
	case "DELETE":
		if location.ID == 0 {
			msg := map[string]string{"error": "Aucun terrain ne correspond. Il doit déjà avoir été supprimé!"}
			Message(w, msg, http.StatusNotFound)
		} else {
			db.Where("ID = ?", id).Delete(&location)
			Message(w, "", http.StatusNoContent)
		}
	}
}

// GetTerrainsHandler Gère la récupération de tous les terrains de la base de données
func (a *AcquisitionService) GetTerrainsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	locations := []Locations{}
	db.Find(&locations)

	locations = addFieldTypes(db, locations)

	Message(w, locations, http.StatusOK)
}

// CreerTerrainHandler Gère la création d'un terrain dans la base de données
func (a *AcquisitionService) CreerTerrainHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		var l Locations
		err = json.Unmarshal(body, &l)
		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		// On enlève les espaces superflues
		l.Name = strings.TrimSpace(l.Name)
		l.Address = strings.TrimSpace(l.Address)
		l.City = strings.TrimSpace(l.City)

		if l.Name == "" || l.Address == "" || l.City == "" {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			Message(w, msg, http.StatusBadRequest)
		} else {

			locations := []Locations{}
			name := strings.ToLower(strings.TrimSpace(l.Name))
			db.Where("LOWER(Name) = LOWER(?)", name).Find(&locations)

			if len(locations) > 0 {
				msg := map[string]string{"error": "Un terrain de même nom existe déjà. Veuillez choisir un autre nom."}
				Message(w, msg, http.StatusBadRequest)
			} else {
				db.Create(&l)
				Message(w, l, http.StatusCreated)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez remplir tous les champs."}
		Message(w, msg, http.StatusBadRequest)
	}
}

func addFieldTypes(db *gorm.DB, locations []Locations) []Locations {
	for i := range locations {
		var ft FieldTypes
		db.Where("ID = ?", locations[i].FieldTypesID).Find(&ft)
		locations[i].FieldType = ft
	}

	return locations
}
