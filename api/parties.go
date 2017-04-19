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

	// DarkSky api for historical weather
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
			db.Where("team_id = ? AND opposing_team = ? AND Date = ?",
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
	case "PUT":
		id := mux.Vars(r)["id"]
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

			// On vérifie que la partie existe bien
			game := Games{}
			db.First(&game, "ID = ?", id)

			if game.ID == 0 {
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
						a.ErrorHandler(w, err)
						return
					}
					date := time.Unix()
					d := strconv.Itoa(int(date))

					f, err := forecast.Get(a.keys.Weather, lat, lng, d, forecast.CA, forecast.French)
					if err != nil {
						a.ErrorHandler(w, err)
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
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
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
