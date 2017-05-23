//
// Fichier     : parties.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'une partie à analyser.
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"

	// DarkSky api pour la météo historique
	forecast "github.com/mlbright/forecast/v2"
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
			games[i] = ajoutInfosPartie(db, games[i])
		}

		Message(w, games, http.StatusOK)
	case "POST":
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		body, _ := ioutil.ReadAll(r.Body)
		if len(body) > 0 {
			var g Games

			if err = json.Unmarshal(body, &g); err != nil {
				msg := map[string]string{"error": "Une erreur est survenue lors de la création de la partie. Les données entrées sont invalides"}
				Message(w, msg, http.StatusBadRequest)
				return
			}

			g.OpposingTeam = strings.TrimSpace(g.OpposingTeam)

			// On vérifie que la partie n'existe pas déjà
			game := []Games{}
			db.Where("team_id = ? AND opposing_team = ? AND Date = ?",
				g.TeamID, g.OpposingTeam, g.Date).Find(&game)

			if len(game) > 0 {
				msg := map[string]string{"error": "Une partie de même date avec les mêmes equipes existe déjà!"}
				Message(w, msg, http.StatusBadRequest)
			} else {
				if g.Date == "" || g.TeamID < 1 || g.OpposingTeam == "" || g.FieldCondition == "" ||
					g.Status == "" || g.LocationID < 1 || g.SeasonID < 1 {
					msg := map[string]string{"error": "Veuillez remplir tous les champs!"}
					Message(w, msg, http.StatusBadRequest)
				} else {
					db.Create(&g)
					g = ajoutInfosPartie(db, g)
					Message(w, g, http.StatusCreated)
				}
			}
		} else {
			var g Games
			db.Create(&g)
			msg := map[string]string{"game_id": strconv.Itoa(int(g.ID))}
			Message(w, msg, http.StatusCreated)
		}
	}
}

// PartieHandler Gère la récupération et la création des parties
func (a *AcquisitionService) PartieHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := mux.Vars(r)["id"]
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		var game Games
		db.Where("ID = ?", id).Find(&game)

		if game.ID == 0 {
			msg := map[string]string{"error": "Aucune partie ne correspond"}
			Message(w, msg, http.StatusNotFound)
			return
		}
		ajoutInfosPartie(db, game)
		Message(w, game, http.StatusOK)
	case "PUT":
		id := mux.Vars(r)["id"]
		body, _ := ioutil.ReadAll(r.Body)
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

			// On vérifie que la partie existe bien
			game := Games{}
			db.First(&game, "ID = ?", id)

			if game.ID == 0 {
				db.Create(&g)
				g = ajoutInfosPartie(db, g)
				Message(w, g, http.StatusCreated)
			} else {
				var l Locations
				var temp string
				var summary string

				db.First(&l, "ID = ?", g.LocationID)

				if !l.IsInside {
					geocoder.SetAPIKey(a.keys.Geodecoder)
					query := l.City + " Canada" // l.Address TODO

					latitude, longitude, err := geocoder.Geocode(query)
					if err != nil {
						msg := map[string]string{"error": "Une erreur inconnue est survenue lors de la création de la partie. Veuillez réessayer"}
						Message(w, msg, http.StatusBadRequest)
						return
					}
					lat := strconv.FormatFloat(latitude, 'f', 10, 64)
					lng := strconv.FormatFloat(longitude, 'f', 10, 64)

					layout := "2006-01-02 15:04"
					time, err := time.Parse(layout, g.Date)
					if err != nil {
						msg := map[string]string{"error": err.Error()}
						Message(w, msg, http.StatusBadRequest)
						return
					}
					date := time.Unix()
					d := strconv.Itoa(int(date))

					f, err := forecast.Get(a.keys.Weather, lat, lng, d, forecast.CA, forecast.French)
					if err != nil {
						msg := map[string]string{"error": err.Error()}
						Message(w, msg, http.StatusBadRequest)
						return
					}

					temp = strconv.FormatFloat(f.Currently.Temperature, 'f', 2, 64)
					summary = f.Currently.Summary
				}

				// Modification de la partie
				game.TeamID = g.TeamID
				game.OpposingTeam = g.OpposingTeam
				game.Status = g.Status
				game.SeasonID = g.SeasonID
				game.LocationID = g.LocationID
				game.Date = g.Date
				game.FieldCondition = g.FieldCondition
				game.Degree = temp
				game.Temperature = summary

				db.Model(&game).Where("ID = ?", id).Updates(game)

				Message(w, game, http.StatusOK)
			}
		} else {
			msg := map[string]string{"error": "Veuillez remplir tous les champs."}
			Message(w, msg, http.StatusBadRequest)
		}
	}
}

// SupprimerPartiesHandler Gère la suppression des parties
func (a *AcquisitionService) SupprimerPartiesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	g := Games{}
	id := strings.ToLower(strings.TrimSpace(vars["id"]))
	db.First(&g, "ID = ?", id)

	// Erreur
	if g.ID == 0 {
		msg := map[string]string{"error": "Aucune partie ne correspond. Elle doit déjà avoir été supprimée!"}
		Message(w, msg, http.StatusBadRequest)
	} else {
		// On supprime l'équipe
		db.Where("ID = ?", id).Delete(&g)
		Message(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	}
}

// ajoutInfosPartie ajout des informations sur une parties
// à la structure de celle-ci
func ajoutInfosPartie(db *gorm.DB, g Games) Games {
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
