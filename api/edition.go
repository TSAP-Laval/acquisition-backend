package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

func (a *AcquisitionService) GetJoueurs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	user := []Players{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) GetActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	user := []ActionsType{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) PostAction(w http.ResponseWriter, r *http.Request) {
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
	var t Actions
	err = json.Unmarshal(body, &t)
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	log.Println(t.ZoneID)
	if db.NewRecord(t) {
		db.Create(&t)
		db.NewRecord(t)
		w.Header().Set("Content-Type", "application/text")

	} else {

		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}

}
