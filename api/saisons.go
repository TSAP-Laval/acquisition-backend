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

func (a *AcquisitionService) PostSaison(w http.ResponseWriter, r *http.Request) {
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
	var t Seasons
	err = json.Unmarshal(body, &t)
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

}

func (a *AcquisitionService) GetSeasons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	strucSaison := []Seasons{}
	db.Find(&strucSaison)

	SaisonJSON, _ := json.Marshal(strucSaison)
	fmt.Println(string(SaisonJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SaisonJSON)
}
