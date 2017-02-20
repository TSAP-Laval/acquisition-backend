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

		equipe := []Equipe{}
		nom := strings.ToLower(strings.TrimSpace(vars["nom"]))
		db.Where("LOWER(Nom) LIKE LOWER(?)", "%"+nom+"%").Find(&equipe)

		for i := 0; i < len(equipe); i++ {
			equipe[i] = AjoutNiveauSport(db, equipe[i])
		}

		equipeJSON, _ := json.Marshal(equipe)

		Message(w, equipeJSON, 200)
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

		equipe := []Equipe{}
		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.First(&equipe, "ID = ?", id)

		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				var e Equipe
				err = json.Unmarshal(body, &e)
				ErrorHandler(w, err)

				e.Nom = strings.TrimSpace(e.Nom)
				e.Ville = strings.TrimSpace(e.Ville)

				var o string

				if e.Nom == "" {
					o += "Nom, "
				}
				if e.Ville == "" {
					o += "Ville, "
				}
				db.Model(&equipe).Where("ID = ?", id).Omit(o).Updates(e)

				// L'équipe modifiée
				var ne Equipe
				db.Where("ID = ?", id).Find(&ne)

				ne = AjoutNiveauSport(db, e)

				equipeJSON, _ := json.Marshal(ne)
				Message(w, equipeJSON, 201)

			} else if err != nil {
				ErrorHandler(w, err)
			} else {
				msg := map[string]string{"error": "Veuillez choisir au moins un champs à modifier."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 400)
			}
		case "DELETE":
			if len(equipe) == 0 {
				msg := map[string]string{"error": "Aucune equipe ne correspond. Elle doit déjà avoir été supprimée!"}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 204)
			} else {
				db.Where("ID = ?", id).Delete(&equipe)
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

	equipe := []Equipe{}
	db.Find(&equipe)

	for i := 0; i < len(equipe); i++ {
		equipe[i] = AjoutNiveauSport(db, equipe[i])
	}

	equipeJSON, _ := json.Marshal(equipe)

	Message(w, equipeJSON, 200)
}

// CreerEquipeHandler Gère la création d'une équipe dans la base de donnée
func (a *AcquisitionService) CreerEquipeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		ErrorHandler(w, err)

		var e Equipe
		err = json.Unmarshal(body, &e)
		ErrorHandler(w, err)

		// On enlève les espaces superflues
		e.Nom = strings.TrimSpace(e.Nom)
		e.Ville = strings.TrimSpace(e.Ville)

		if e.Nom == "" || e.Ville == "" {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			errorJSON, _ := json.Marshal(msg)
			Message(w, errorJSON, 400)
		} else {

			equipe := []Equipe{}
			nom := strings.ToLower(strings.TrimSpace(e.Nom))
			db.Where("LOWER(Nom) = LOWER(?)", nom).Find(&equipe)

			if len(equipe) > 0 {
				msg := map[string]string{"error": "Une équipe de même nom existe déjà. Veuillez choisir une autre nom."}
				errorJSON, _ := json.Marshal(msg)
				Message(w, errorJSON, 401)
			} else {
				if db.NewRecord(e) {
					db.Create(&e)
					if db.NewRecord(e) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création de l'équipe. Veuillez réessayer!"}
						errorJSON, _ := json.Marshal(msg)
						Message(w, errorJSON, 500)
					} else {
						e = AjoutNiveauSport(db, e)
						succesJSON, _ := json.Marshal(e)
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
func AjoutNiveauSport(db *gorm.DB, e Equipe) Equipe {
	// Ajout du sport pour l'affichage
	var s Sport
	db.Where("ID = ?", e.SportID).Find(&s)
	if s.Nom != "" {
		e.Sport = s
	}

	// Ajout du niveau pour l'affichage
	var n Niveau
	db.Where("ID = ?", e.NiveauID).Find(&n)
	if n.Nom != "" {
		e.Niveau = n
	}

	return e
}
