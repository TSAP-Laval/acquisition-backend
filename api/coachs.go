package api

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"encoding/json"

	"github.com/jinzhu/gorm"
)

//GetCoachsHandler :  fetch all created coachs
func (a *AcquisitionService) GetCoachsHandler(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	fmt.Println(err)

	coach := []Entraineur{}

	db.Find(&coach)

	coachJSON, _ := json.Marshal(coach)
	fmt.Println(string(coachJSON))

	w.Header().Set("Content-Type", "Application/json")
	w.Write(coachJSON)
	db.Close()
}

//PostCoach : Create a new coach in the database
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

	var newCoach Entraineur

	err = json.Unmarshal(body, &newCoach)

	fmt.Println(err)

	if err != nil {
		panic(err)
	}

	if db.NewRecord(newCoach) {
		db.Create(&newCoach)
		db.NewRecord(&newCoach)
	} else {
		fmt.Println("Not created")
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

//UpdateCoachHandler : update coach's teams
/*func (a *AcquisitionService) UpdateCoachHandler(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	fmt.Println(err)

	var coachUpdate Entraineur

	paramMail := r.URL.Query().Get("email")
	paramTeams := r.URL.Query().Get("email")

	db.Where("email LIKE ?", paramMail).First(&coachUpdate)

	if &coachUpdate != nil {
		//db.Model(&coachUpdate).upda
	}

}*/
