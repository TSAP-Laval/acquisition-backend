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

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	fmt.Println(err)

	coach := []Coaches{}

	db.Find(&coach)

	coachJSON, _ := json.Marshal(coach)
	fmt.Println(string(coachJSON))

	w.Header().Set("Content-Type", "Application/json")
	w.Write(coachJSON)
	db.Close()
}

//PostCoachHandler : Create a new coach in the database
func (a *AcquisitionService) PostCoachHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	fmt.Println(err)

	body, errorBody := ioutil.ReadAll(r.Body)

	fmt.Printf("-----------------------")
	fmt.Println(body)
	fmt.Printf("-----------------------")
	if err != nil {
		defer db.Close()
		panic(errorBody)
	}

	fmt.Println(string(body))

	var newCoach Coaches
	var dat map[string]interface{}
	err = json.Unmarshal(body, &newCoach)
	err = json.Unmarshal(body, &dat)
	num := dat["Teams"]

	Team := Teams{}

	db.First(&Team, num)

	newCoach.Teams = append(newCoach.Teams, Team)
	fmt.Println(err)

	if err != nil {
		panic(err)
	}

	if db.NewRecord(newCoach) {
		db.Create(&newCoach)
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusCreated)

	} else {
		fmt.Println("Not created")
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
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

// GetEquipeByIDHandler Gère la récupération des équipes correspondant au nom entré
func (a *AcquisitionService) GetEquipeByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars != nil {
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		team := []Teams{}
		ID := strings.ToLower(strings.TrimSpace(vars["ID"]))
		db.Where("ID LIKE ?", ID).Find(&team)

		for i := 0; i < len(team); i++ {
			team[i] = AjoutNiveauSport(db, team[i])
		}

		Message(w, team, http.StatusOK)
	} else {
		msg := map[string]string{"error": "Veuillez entrer un nom d'équipe ou en créer une préalablement"}
		Message(w, msg, http.StatusNotFound)
	}
}
