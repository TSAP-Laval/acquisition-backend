//
// Fichier     : coaches.go
// Développeur : Mehdi Laribi
//

package api

import (
	"io/ioutil"
	"net/http"

	"encoding/json"

	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	//Import driver

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//GetCoachesHandler :  fetch all created coaches
//Si l'identifiant est vide, la fonction retourne tous les coachs
func (a *AcquisitionService) GetCoachesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/json")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	vars := mux.Vars(r)

	idCoach := strings.ToLower(strings.TrimSpace(vars["coachID"]))

	if idCoach == "" {
		coaches := []Coaches{}
		db.Find(&coaches)

		for i := 0; i < len(coaches); i++ {
			var c Coaches
			c = coaches[i]
			coaches[i] = AjoutCoachInfo(db, c)
		}

		coachJSON, _ := json.Marshal(coaches)
		w.Write(coachJSON)
	} else {
		coach := Coaches{}
		var id = vars["ID"]
		db.First(coach, id)
		coach = AjoutCoachInfo(db, coach)

		coachJSON, _ := json.Marshal(coach)
		w.Write(coachJSON)
	}

	db.Close()
}

//PostCoachHandler : Create a new coach in the database
func (a *AcquisitionService) PostCoachHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	if len(body) > 0 {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		var newCoach Coaches
		var dat map[string]string
		err = json.Unmarshal(body, &newCoach)
		err = json.Unmarshal(body, &dat)

		var seaID = dat["SeasonID"]

		var num = dat["TeamsIDs"]
		ids := strings.Split(num, ",")

		Team := Teams{}
		db.Find(&Team, ids)

		Saison := Seasons{}
		db.Find(&Saison, seaID)

		db.Create(&newCoach)
		newCoach = AjoutCoachInfo(db, newCoach)
		newCoach.Teams = append(newCoach.Teams, Team)
		newCoach.Season = Saison

		db.Model(&Team).Association("Teams").Append(newCoach)
		w.WriteHeader(http.StatusCreated)
	} else if err != nil {
		a.ErrorHandler(w, err)
		return
	} else {
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
		msg := map[string]string{"error": "Veuillez remplir tous les champs."}
		Message(w, msg, http.StatusBadRequest)
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

//AssignerEquipeCoach : Assigne des equipes au coach
func (a *AcquisitionService) AssignerEquipeCoach(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	var c Coaches
	err = json.Unmarshal(body, &c)

	if err != nil {
		a.ErrorHandler(w, err)
	} else {

		id := strings.ToLower(strings.TrimSpace(vars["id"]))
		db.Model(&c).Where("ID = ?", id).Updates(c)

		Message(w, "Teams for this coach : OK", http.StatusCreated)
	}

}

// AjoutCoachInfo : Ajout des équipes au coach
func AjoutCoachInfo(db *gorm.DB, c Coaches) Coaches {

	var ts []Teams
	db.Where("ID in (?)", c.TeamsIDs).Find(&ts)
	if len(ts) > 0 {
		c.Teams = ts
	}
	return c

}
