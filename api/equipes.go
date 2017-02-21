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

// GetEquipeHandler Gère la récupération des équipes correspondant au nom entré
func (a *AcquisitionService) GetEquipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		a.ErrorHandler(w, err)

		team := []Teams{}
		name := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Where("LOWER(Name) LIKE LOWER(?)", name+"%").Find(&team)

		for i := 0; i < len(team); i++ {
			team[i] = AjoutNiveauSport(db, team[i])
		}

		Message(w, team, http.StatusOK)
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom d'équipe ou en créer une préalablement"}
		Message(w, msg, http.StatusNotFound)
	}
}

// EquipesHandler gère la modification et la suppression des équipes
func (a *AcquisitionService) EquipesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		a.ErrorHandler(w, err)

		team := []Teams{}
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.First(&team, "ID = ?", id)

		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				var t Teams
				err = json.Unmarshal(body, &t)
				a.ErrorHandler(w, err)

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

				nt = AjoutNiveauSport(db, t)

				Message(w, nt, http.StatusCreated)

			} else if err != nil {
				a.ErrorHandler(w, err)
			} else {
				msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
				Message(w, msg, http.StatusBadRequest)
			}
		case "DELETE":
			// Erreur
			if len(team) == 0 {
				msg := map[string]string{"error": "Aucune equipe ne correspond. Elle doit déjà avoir été supprimée!"}
				Message(w, msg, http.StatusNoContent)
			} else {
				// On supprime l'équipe
				db.Where("ID = ?", id).Delete(&team)
				msg := map[string]string{"succes": "L'équipe a été supprimée avec succès!"}
				Message(w, msg, http.StatusNoContent)
			}
		}
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom d'équipe ou en créer une préalablement."}
		Message(w, msg, http.StatusBadRequest)
	}
}

// GetEquipesHandler Gère la récupération de toutes les équipes de la base de donnée
func (a *AcquisitionService) GetEquipesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	a.ErrorHandler(w, err)

	team := []Teams{}
	db.Find(&team)

	for i := 0; i < len(team); i++ {
		team[i] = AjoutNiveauSport(db, team[i])
	}

	Message(w, team, http.StatusOK)
}

// CreerEquipeHandler Gère la création d'une équipe dans la base de donnée
func (a *AcquisitionService) CreerEquipeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		a.ErrorHandler(w, err)

		var t Teams
		err = json.Unmarshal(body, &t)
		a.ErrorHandler(w, err)

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
				Message(w, msg, http.StatusUnauthorized)
			} else {
				if db.NewRecord(t) {
					db.Create(&t)
					if db.NewRecord(t) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création de l'équipe. Veuillez réessayer!"}
						Message(w, msg, http.StatusInternalServerError)
					} else {
						t = AjoutNiveauSport(db, t)
						Message(w, t, http.StatusCreated)
					}
				} else {
					msg := map[string]string{"error": "L'équipe existe déjà dans la base de donnée!"}
					Message(w, msg, http.StatusBadRequest)
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

// Message Gère les messages (erreurs, messages de succès) à envoyer au client
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
		w.WriteHeader(http.StatusNotFound)
		a.ErrorWrite(err.Error(), w)
		return
	}
}
