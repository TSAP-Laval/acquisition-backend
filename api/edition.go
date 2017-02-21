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

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)

	user := []Players{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) GetActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)
	user := []ActionsType{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
func (a *AcquisitionService) PostAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t Actions
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
