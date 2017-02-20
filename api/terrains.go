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

	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		a.ErrorHandler(w, err)

		location := []Locations{}
		name := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Where("LOWER(Name) LIKE LOWER(?)", name+"%").Find(&location)

		Message(w, location, http.StatusOK)
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom de terrain ou en créer un préalablement"}
		Message(w, msg, http.StatusBadRequest)
	}
}

// TerrainsHandler gère la modification et la suppression de terrains
func (a *AcquisitionService) TerrainsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		a.ErrorHandler(w, err)

		location := []Locations{}
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.First(&location, "ID = ?", id)

		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				var l Locations
				err = json.Unmarshal(body, &l)
				a.ErrorHandler(w, err)

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
			} else if err != nil {
				a.ErrorHandler(w, err)
			} else {
				msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
				Message(w, msg, http.StatusBadRequest)
			}
		case "DELETE":
			if len(location) == 0 {
				msg := map[string]string{"error": "Aucun terrain ne correspond. Il doit déjà avoir été supprimé!"}
				Message(w, msg, http.StatusNoContent)
			} else {
				db.Where("ID = ?", id).Delete(&location)
				msg := map[string]string{"succes": "Le terrain a été supprimé avec succès!"}
				Message(w, msg, http.StatusNoContent)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom de terrain ou en créer un préalablement."}
		Message(w, msg, http.StatusNotFound)
	}
}

// GetTerrainsHandler Gère la récupération de tous les terrains de la base de donnée
func (a *AcquisitionService) GetTerrainsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	a.ErrorHandler(w, err)

	locations := []Locations{}
	db.Find(&locations)

	Message(w, locations, http.StatusOK)
}

// CreerTerrainHandler Gère la création de terrain dans la base de donnée
func (a *AcquisitionService) CreerTerrainHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		a.ErrorHandler(w, err)

		var l Locations
		err = json.Unmarshal(body, &l)
		a.ErrorHandler(w, err)

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
				Message(w, msg, http.StatusUnauthorized)
			} else {
				if db.NewRecord(l) {
					db.Create(&l)
					if db.NewRecord(l) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création du terrain. Veuillez réessayer!"}
						Message(w, msg, http.StatusInternalServerError)
					} else {
						Message(w, l, http.StatusCreated)
					}
				} else {
					msg := map[string]string{"error": "Le terrain existe déjà dans la base de donnée!"}
					Message(w, msg, http.StatusUnauthorized)
				}
			}
		}
	} else if err != nil {
		a.ErrorHandler(w, err)
	} else {
		msg := map[string]string{"error": "Veuillez remplir tous les champs."}
		Message(w, msg, http.StatusBadRequest)
	}
}
