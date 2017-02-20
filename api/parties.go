package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// PartiesHandler Gère la récupération et la création des parties
func (a *AcquisitionService) PartiesHandler(w http.ResponseWriter, r *http.Request) {
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
