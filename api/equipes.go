package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetEquipeHandler Gère la récupération des équipes correspondant au nom entré
func (a *AcquisitionService) GetEquipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		team := []Teams{}
		name := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Where("LOWER(Nom) LIKE LOWER(?)", "%"+name+"%").Find(&team)

		for i := 0; i < len(team); i++ {
			team[i] = AjoutNiveauSport(db, team[i])
		}

		teamJSON, _ := json.Marshal(team)

		Message(w, teamJSON, 200)
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom d'équipe ou en créer une préalablement"}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, 404)
	}
}

// EquipesHandler gère la modification et la suppression des équipes
func (a *AcquisitionService) EquipesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		team := []Teams{}
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.First(&team, "ID = ?", id)

		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				var t Teams
				err = json.Unmarshal(body, &t)
				ErrorHandler(w, err)

				t.Name = strings.TrimSpace(t.Name)
				t.City = strings.TrimSpace(t.City)

				var o string

				if t.Name == "" {
					o += "Nom, "
				}
				if t.City == "" {
					o += "Ville, "
				}
				db.Model(&team).Where("ID = ?", id).Omit(o).Updates(t)

				// L'équipe modifiée
				var nt Teams
				db.Where("ID = ?", id).Find(&nt)

				nt = AjoutNiveauSport(db, t)

				teamJSON, _ := json.Marshal(nt)
				Message(w, teamJSON, 201)

			} else if err != nil {
				ErrorHandler(w, err)
			} else {
				msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 400)
			}
		case "DELETE":
			if len(team) == 0 {
				msg := map[string]string{"error": "Aucune equipe ne correspond. Elle doit déjà avoir été supprimée!"}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 204)
			} else {
				db.Where("ID = ?", id).Delete(&team)
				msg := map[string]string{"succes": "L'équipe a été supprimée avec succès!"}
				succesJSON, _ := json.Marshal(msg)
				Message(w, succesJSON, 204)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom d'équipe ou en créer une préalablement."}
		errorJSON, _ := json.Marshal(msg)
		Message(w, errorJSON, 404)
	}
}

// GetEquipesHandler Gère la récupération de toutes les équipes de la base de donnée
func (a *AcquisitionService) GetEquipesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	ErrorHandler(w, err)

	team := []Teams{}
	db.Find(&team)

	for i := 0; i < len(team); i++ {
		team[i] = AjoutNiveauSport(db, team[i])
	}

	teamJSON, _ := json.Marshal(team)

	Message(w, teamJSON, 200)
}

// CreerEquipeHandler Gère la création d'une équipe dans la base de donnée
func (a *AcquisitionService) CreerEquipeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		var t Teams
		err = json.Unmarshal(body, &t)
		ErrorHandler(w, err)

		// On enlève les espaces superflues
		t.Name = strings.TrimSpace(t.Name)
		t.City = strings.TrimSpace(t.City)

		if t.Name == "" || t.City == "" {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			errorJSON, _ := json.Marshal(msg)
			Message(w, errorJSON, 400)
		} else {

			team := []Teams{}
			name := strings.ToLower(strings.TrimSpace(t.Name))
			db.Where("LOWER(Nom) = LOWER(?)", name).Find(&team)

			if len(team) > 0 {
				msg := map[string]string{"error": "Une équipe de même nom existe déjà. Veuillez choisir une autre nom."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 401)
			} else {
				if db.NewRecord(t) {
					db.Create(&t)
					if db.NewRecord(t) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création de l'équipe. Veuillez réessayer!"}
						errorJSON, _ := json.Marshal(msg)
						Message(w, errorJSON, 500)
					} else {
						t = AjoutNiveauSport(db, t)
						succesJSON, _ := json.Marshal(t)
						Message(w, succesJSON, 201)
					}
				} else {
					msg := map[string]string{"error": "L'équipe existe déjà dans la base de donnée!"}
					errorJSON, _ := json.Marshal(msg)
					Message(w, errorJSON, 400)
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

// Message Gère les messages (erreurs, messages de succès) à envoyer au client
func Message(w http.ResponseWriter, msg []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msg)
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
