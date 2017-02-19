package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

// EquipesHandler Gère la création d'une partie avec les informations sur les équipes
func (a *AcquisitionService) EquipesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			fmt.Print("\nERROR : ")
			fmt.Println(err)
		}

		mvmType := []MovementType{}
		db.Find(&mvmType)

		mvmTypeJSON, _ := json.Marshal(mvmType)

		w.Header().Set("Content-Type", "Application/json")
		w.Write(mvmTypeJSON)
	case "POST":
	case "PUT":
	}
}
