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

// HandleJoueur gère la modification et l'ajout de joueur
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
		if num != "" {
			Team := Teams{}
			db.First(&Team, num)
			if err != nil {
				a.ErrorHandler(w, err)
				return
			}
			t.Teams = append(t.Teams, Team)
			db.Model(&Team).Association("Players").Append(t)
		}
		switch r.Method {
		case "POST":

			if db.NewRecord(t) {
				db.Create(&t)
				db.NewRecord(t)
				w.Header().Set("Content-Type", "application/text")
				w.WriteHeader(http.StatusCreated)

			} else {

				Message(w, "déjà créé", http.StatusBadRequest)
			}
		case "PUT":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			id := strings.ToLower(strings.TrimSpace(vars["id"]))
			db.Model(&t).Where("ID = ?", id).Updates(t)
			Message(w, "ok", http.StatusOK)
		}
	}

}
