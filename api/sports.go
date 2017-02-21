package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

func (a *AcquisitionService) GetSports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

	defer db.Close()
	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	strucSport := []Sports{}
	db.Find(&strucSport)

	SportJSON, _ := json.Marshal(strucSport)
	fmt.Println(string(SportJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SportJSON)
}
