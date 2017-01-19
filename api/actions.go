package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//MovementType represent Movement type entity
//1: Offensive
//2: Defensive
//Neutral
type MovementType struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetMovementTypeHandler Gestion du select des types de mouvements
func (a *AcquisitionService) GetMovementTypeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=DB_TSAP sslmode=disable password=tsaplaval")

	defer db.Close()
	fmt.Println(err)

	mvmType := []MovementType{}
	db.Find(&mvmType)

	mvmTypeJSON, _ := json.Marshal(mvmType)
	fmt.Println(string(mvmTypeJSON))

	w.Header().Set("Content-Type", "Application/json")
	w.Write(mvmTypeJSON)
}
