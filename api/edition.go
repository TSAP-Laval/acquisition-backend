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

<<<<<<< HEAD
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
=======
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
>>>>>>> 9a0cab33f0bdf108302a9dc72005a8fd4abbfb44
	defer db.Close()
	fmt.Println(err)

	user := []Joueur{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) GetActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
<<<<<<< HEAD
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
=======
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
>>>>>>> 9a0cab33f0bdf108302a9dc72005a8fd4abbfb44

	defer db.Close()
	fmt.Println(err)
	user := []TypeAction{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) PostJoueur(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
<<<<<<< HEAD
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
=======
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
>>>>>>> 9a0cab33f0bdf108302a9dc72005a8fd4abbfb44

	defer db.Close()
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t Action
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t.ZoneID)
	if db.NewRecord(t) {
		db.Create(&t)
		db.NewRecord(t)
		w.Header().Set("Content-Type", "application/text")

		w.Write([]byte("ok"))
	} else {
		fmt.Println("Test22")
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}

}
