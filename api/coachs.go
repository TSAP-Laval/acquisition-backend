//
// Fichier     : coaches.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonctions...
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

// TODO: Changer le nom du fichier et ses références pour coach au pluriel...
//		http://www.wordhippo.com/what-is/the-plural-of/coach.html
// TODO: Linter le code...
// TODO: Gérer les erreurs comme du monde
// TODO: Enlever tous ce qui est log, print...

//GetCoachsHandler :  fetch all created coachs
func (a *AcquisitionService) GetCoachsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	coach := []Coaches{}

	db.Find(&coach)

	coachJSON, _ := json.Marshal(coach)

	w.Header().Set("Content-Type", "Application/json")
	w.Write(coachJSON)
	db.Close()
}

//PostCoachHandler : Create a new coach in the database
func (a *AcquisitionService) PostCoachHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	var newCoach Coaches
	var dat map[string]interface{}
	err = json.Unmarshal(body, &newCoach)
	err = json.Unmarshal(body, &dat)
	num := dat["Teams"]

	Team := Teams{}

	db.First(&Team, num)

	newCoach.Teams = append(newCoach.Teams, Team)

	if err != nil {
		panic(err)
	}

	if db.NewRecord(newCoach) {
		db.Create(&newCoach)
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusCreated)

	} else {
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
