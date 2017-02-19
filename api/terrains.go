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

		lieu := []Lieu{}
		nom := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Where("LOWER(Nom) LIKE LOWER(?)", "%"+nom+"%").Find(&lieu)

		lieuJSON, _ := json.Marshal(lieu)

		Message(w, lieuJSON, false)
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom de terrain ou en créer un préalablement"}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, true)
	}
}

// TerrainsHandler gère la modification et la suppression de terrains
func (a *AcquisitionService) TerrainsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		lieu := []Lieu{}
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.First(&lieu, "ID = ?", id)

		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				var l Lieu
				err = json.Unmarshal(body, &l)
				ErrorHandler(w, err)

				l.Nom = strings.TrimSpace(l.Nom)
				l.Ville = strings.TrimSpace(l.Ville)
				l.Adresse = strings.TrimSpace(l.Adresse)

				var o string

				if l.Nom == "" {
					o += "Nom, "
				}
				if l.Ville == "" {
					o += "Ville, "
				}
				if l.Adresse == "" {
					o += "Adresse, "
				}

				db.Model(&lieu).Omit(o).Updates(l)

				// Le lieu modifié
				var nl Lieu
				db.Where("ID = ?", id).Find(&nl)
				equipeJSON, _ := json.Marshal(nl)
				Message(w, equipeJSON, false)
			} else if err != nil {
				ErrorHandler(w, err)
			} else {
				msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, true)
			}
		case "DELETE":
			if len(lieu) == 0 {
				msg := map[string]string{"error": "Aucun terrain ne correspond. Il doit déjà avoir été supprimé!"}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, true)
			} else {
				db.Where("ID = ?", id).Delete(&lieu)
				msg := map[string]string{"succes": "Le terrain a été supprimé avec succès!"}
				succesJSON, _ := json.Marshal(msg)
				Message(w, succesJSON, false)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom de terrain ou en créer un préalablement."}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, true)
	}
}

// GetTerrainsHandler Gère la récupération de tous les terrains de la base de donnée
func (a *AcquisitionService) GetTerrainsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	ErrorHandler(w, err)

	lieu := []Lieu{}
	db.Find(&lieu)

	lieuJSON, _ := json.Marshal(lieu)

	Message(w, lieuJSON, false)
}

// CreerTerrainHandler Gère la création de terrain dans la base de donnée
func (a *AcquisitionService) CreerTerrainHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		var l Lieu
		err = json.Unmarshal(body, &l)

		// On enlève les espaces superflues
		l.Nom = strings.TrimSpace(l.Nom)
		l.Adresse = strings.TrimSpace(l.Adresse)
		l.Ville = strings.TrimSpace(l.Ville)

		if l.Nom == "" || l.Adresse == "" || l.Ville == "" {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			errorJSON, _ := json.Marshal(msg)
			Message(w, errorJSON, true)
		} else {

			lieu := []Lieu{}
			nom := strings.ToLower(strings.TrimSpace(l.Nom))
			db.Where("LOWER(Nom) = LOWER(?)", nom).Find(&lieu)

			if len(lieu) > 0 {
				msg := map[string]string{"error": "Un terrain de même nom existe déjà. Veuillez choisir un autre nom."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, true)
			} else {
				if db.NewRecord(l) {
					db.Create(&l)
					if db.NewRecord(l) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création du terrain. Veuillez réessayer!"}
						errorJSON, _ := json.Marshal(msg)
						Message(w, errorJSON, true)
					} else {

						succesJSON, _ := json.Marshal(l)
						Message(w, succesJSON, false)
					}
				} else {
					msg := map[string]string{"error": "Le terrain existe déjà dans la base de donnée!"}
					errorJSON, _ := json.Marshal(msg)
					Message(w, errorJSON, true)
				}
			}
		}
	} else if err != nil {
		ErrorHandler(w, err)
	} else {
		msg := map[string]string{"error": "Veuillez remplir tous les champs."}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, true)
	}
}

// ErrorHandler gère les erreurs côté serveur
func ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Print("\nERROR : ")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
