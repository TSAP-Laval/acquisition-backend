package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
)

// PartiesHandler Gère la récupération et la création des parties
func (a *AcquisitionService) PartiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		games := []Games{}
		db.Find(&games)

		for i := 0; i < len(games); i++ {
			games[i] = AjoutInfosPartie(db, games[i])
		}

		Message(w, games, http.StatusOK)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if len(body) > 0 {
			db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
			defer db.Close()

			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			var g Games
			err = json.Unmarshal(body, &g)
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			g.OpposingTeam = strings.TrimSpace(g.OpposingTeam)

			// On vérifie que la partie n'existe pas déjà
			game := []Games{}
			db.Where("home_team_id = ? AND opposing_team = ? AND Date = ?",
				g.TeamID, g.OpposingTeam, g.Date).Find(&game)

			if len(game) > 0 {
				msg := map[string]string{"error": "Une partie de même date avec les mêmes equipes existe déjà!"}
				Message(w, msg, http.StatusUnauthorized)
			} else {
				if db.NewRecord(g) {
					db.Create(&g)
					if db.NewRecord(g) {
						msg := map[string]string{"error": "Une erreur est survenue lors de la création de la partie. Veuillez réessayer!"}
						Message(w, msg, http.StatusInternalServerError)
					} else {
						g = AjoutInfosPartie(db, g)
						Message(w, g, http.StatusCreated)
					}
				} else {
					msg := map[string]string{"error": "La partie existe déjà dans la base de donnée!"}
					Message(w, msg, http.StatusBadRequest)
				}
			}
		} else if err != nil {
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}
		} else {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			Message(w, msg, http.StatusBadRequest)
		}
	}
}

// SupprimerPartiesHandler Gère la suppression des parties
func (a *AcquisitionService) SupprimerPartiesHandler(w http.ResponseWriter, r *http.Request) {
	return
}

// AjoutInfosPartie ajout des informations sur une parties
// à la structure de celle-ci
func AjoutInfosPartie(db *gorm.DB, g Games) Games {
	// Home team
	var ht Teams
	db.Where("ID = ?", g.TeamID).Find(&ht)
	if ht.Name != "" {
		ht = AjoutNiveauSport(db, ht)
		g.Team = ht
	}

	// Ajout du lieu pour l'affichage
	var l Locations
	db.Where("ID = ?", g.LocationID).Find(&l)
	if l.Name != "" {
		g.Location = l
	}

	return g
}
