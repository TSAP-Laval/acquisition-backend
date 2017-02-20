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

// GetTerrainHandler Gère la récupération des terrains correspondant au nom entré
func (a *AcquisitionService) GetTerrainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		location := []Locations{}
		name := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Where("LOWER(Name) LIKE LOWER(?)", name+"%").Find(&location)

		locationJSON, _ := json.Marshal(location)

		Message(w, locationJSON, 200)
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom de terrain ou en créer un préalablement"}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, 400)
	}
}

// TerrainsHandler gère la modification et la suppression de terrains
func (a *AcquisitionService) TerrainsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		location := []Locations{}
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.First(&location, "ID = ?", id)

		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				var l Locations
				err = json.Unmarshal(body, &l)
				ErrorHandler(w, err)

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
				equipeJSON, _ := json.Marshal(nl)
				Message(w, equipeJSON, 201)
			} else if err != nil {
				ErrorHandler(w, err)
			} else {
				msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 400)
			}
		case "DELETE":
			if len(location) == 0 {
				msg := map[string]string{"error": "Aucun terrain ne correspond. Il doit déjà avoir été supprimé!"}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 204)
			} else {
				db.Where("ID = ?", id).Delete(&location)
				msg := map[string]string{"succes": "Le terrain a été supprimé avec succès!"}
				succesJSON, _ := json.Marshal(msg)
				Message(w, succesJSON, 204)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom de terrain ou en créer un préalablement."}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, 404)
	}
}

// GetTerrainsHandler Gère la récupération de tous les terrains de la base de donnée
func (a *AcquisitionService) GetTerrainsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	ErrorHandler(w, err)

	locations := []Locations{}
	db.Find(&locations)

	locationsJSON, _ := json.Marshal(locations)

	Message(w, locationsJSON, 200)
}

// CreerTerrainHandler Gère la création de terrain dans la base de donnée
func (a *AcquisitionService) CreerTerrainHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		var l Locations
		err = json.Unmarshal(body, &l)

		// On enlève les espaces superflues
		l.Name = strings.TrimSpace(l.Name)
		l.Address = strings.TrimSpace(l.Address)
		l.City = strings.TrimSpace(l.City)

		if l.Name == "" || l.Address == "" || l.City == "" {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			errorJSON, _ := json.Marshal(msg)
			Message(w, errorJSON, 400)
		} else {

			locations := []Locations{}
			name := strings.ToLower(strings.TrimSpace(l.Name))
			db.Where("LOWER(Name) = LOWER(?)", name).Find(&locations)

			if len(locations) > 0 {
				msg := map[string]string{"error": "Un terrain de même nom existe déjà. Veuillez choisir un autre nom."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 401)
			} else {
				if db.NewRecord(l) {
					db.Create(&l)
					if db.NewRecord(l) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création du terrain. Veuillez réessayer!"}
						errorJSON, _ := json.Marshal(msg)
						Message(w, errorJSON, 500)
					} else {

						succesJSON, _ := json.Marshal(l)
						Message(w, succesJSON, 201)
					}
				} else {
					msg := map[string]string{"error": "Le terrain existe déjà dans la base de donnée!"}
					errorJSON, _ := json.Marshal(msg)
					Message(w, errorJSON, 401)
				}
			}
		}
	} else if err != nil {
		ErrorHandler(w, err)
	} else {
		msg := map[string]string{"error": "Veuillez remplir tous les champs."}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, 400)
	}
}

// ErrorHandler gère les erreurs côté serveur
func ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Print("\nERROR : ")
		fmt.Println(err)
		//w.WriteHeader(http.StatusInternalServerError)
	}
}
