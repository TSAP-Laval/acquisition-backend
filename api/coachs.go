package api

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"encoding/json"

	"log"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	//Import driver

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//GetCoachsHandler :  fetch all created coachs
func (a *AcquisitionService) GetCoachsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/json")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	coachs := []Coaches{}
	db.Find(&coachs)

	for i := 0; i < len(coachs); i++ {
		var c Coaches
		c = coachs[i]
		coachs[i] = AjoutCoachInfo(db, c)
	}

	coachJSON, _ := json.Marshal(coachs)
	fmt.Println(string(coachJSON))

	w.Write(coachJSON)
	db.Close()
}

//PostCoachHandler : Create a new coach in the database
func (a *AcquisitionService) PostCoachHandler(w http.ResponseWriter, r *http.Request) {

	body, errorBody := ioutil.ReadAll(r.Body)

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
		var num = dat["TeamsIDs"]
		ids := strings.Split(num, ",")

		fmt.Println(ids)
		Team := Teams{}
		db.Find(&Team, ids)

		db.Create(&newCoach)
		newCoach = AjoutCoachInfo(db, newCoach)
		newCoach.Teams = append(newCoach.Teams, Team)

		db.Model(&Team).Association("Teams").Append(newCoach)
		w.WriteHeader(http.StatusCreated)
	} else if errorBody != nil {
		a.ErrorHandler(w, errorBody)
		return
	} else {
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
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	log.Println(string(body))
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
