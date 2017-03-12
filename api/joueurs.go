package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func (a *AcquisitionService) HandleJoueur(w http.ResponseWriter, r *http.Request) {
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
	var t Players
	var dat map[string]interface{}
	err = json.Unmarshal(body, &t)
	err = json.Unmarshal(body, &dat)
	num := dat["EquipeID"]
	if err != nil {
		a.ErrorHandler(w, err)

	} else {

		switch r.Method {
		case "POST":

			Team := Teams{}
			db.First(&Team, num)
			t.Teams = append(t.Teams, Team)
			db.Model(&Team).Association("Players").Append(t)
			fmt.Println(num)
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}

			log.Println(t.ID)
			if db.NewRecord(t) {
				db.Create(&t)
				db.NewRecord(t)
				w.Header().Set("Content-Type", "application/text")
				w.WriteHeader(http.StatusCreated)

			} else {
				fmt.Println("Test22")
				w.Header().Set("Content-Type", "application/text")
				w.Write([]byte("erreur"))
			}
		case "PUT":
			id := strings.ToLower(strings.TrimSpace(vars["id"]))
			db.Model(&t).Where("ID = ?", id).Updates(t)

			Message(w, "ok", http.StatusCreated)
		}
	}

}
