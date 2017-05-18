//
// Fichier     : coaches.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'un entraineur.
//

package api

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	//Import driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//GetCoachsHandler Gère la récupération de tous les entraineurs
func (a *AcquisitionService) GetCoachsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	coaches := []Coaches{}
	// Preload permet l'ajout des équipes si l'entraineur en possède
	db.Preload("Teams").Find(&coaches)

	Message(w, coaches, http.StatusOK)
}

//CreerCoachHandler Gère la création des entraineurs
func (a *AcquisitionService) CreerCoachHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		var coach Coaches
		if err = json.Unmarshal(body, &coach); err != nil {
			msg := map[string]string{"error": "Certaines informations entrées sont invalides!"}
			Message(w, msg, http.StatusBadRequest)
			return
		}

		var co Coaches
		db.Model(coach).Where("Email = ?", coach.Email).First(&co)

		if coach.Email == "" || coach.Fname == "" || coach.Lname == "" || coach.PassHash == "" {
			msg := map[string]string{"error": "Certaines informations sont manquantes!"}
			Message(w, msg, http.StatusBadRequest)
		} else {
			if co.ID < 1 {
				db.Create(&coach)
				Message(w, coach, http.StatusCreated)
			} else {
				msg := map[string]string{"error": "Un entraineur avec la même adresse courriel existe déjà"}
				Message(w, msg, http.StatusBadRequest)
			}
		}
	}
}

// ModifierCoachHandler Gère l'assignation d'une équipe à un entraineur
func (a *AcquisitionService) ModifierCoachHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	coachID, _ := strconv.Atoi(mux.Vars(r)["coach-id"])
	teamID, _ := strconv.Atoi(mux.Vars(r)["team-id"])

	var c Coaches
	db.Model(&c).Where("ID = ?", coachID).First(&c)

	if c.ID < 1 {
		msg := map[string]string{"error": "Aucun entraineur ne correspond"}
		Message(w, msg, http.StatusBadRequest)
	} else {

		var te Teams
		db.Model(&te).Where("ID = ?", teamID).First(&te)

		var ct CoachTeam
		db.Where("coach_id = ? AND team_id = ?", coachID, teamID).First(&ct)

		if te.ID < 1 {
			msg := map[string]string{"error": "Aucune équipe ne correspond"}
			Message(w, msg, http.StatusBadRequest)
		} else {
			if ct.ID == 0 {
				// Mise à jour des entraineurs de l'équipe
				ct.CoachID = coachID
				ct.TeamID = teamID
				db.Create(&ct)

				msg := map[string]string{"success": "L'entraineur fait maintenant parti de l'équipe : " + te.Name}
				Message(w, msg, http.StatusOK)
			} else {
				msg := map[string]string{"error": "L'entraineur fait déjà parti de l'équipe"}
				Message(w, msg, http.StatusBadRequest)
			}
		}
	}
}
