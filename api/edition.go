package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

func (a *AcquisitionService) GetJoueurs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
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
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
	defer db.Close()
	fmt.Println(err)
	user := []TypeAction{}
	db.Find(&user)

	userJSON, _ := json.Marshal(user)
	fmt.Println(string(userJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}
