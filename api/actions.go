package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	"io/ioutil"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GetMovementTypeHandler Gestion du select des types de mouvements
func (a *AcquisitionService) GetMovementTypeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=tsaplaval")

	defer db.Close()
	fmt.Println(err)

	mvmType := []MovementType{}
	db.Find(&mvmType)

	mvmTypeJSON, _ := json.Marshal(mvmType)
	fmt.Println(string(mvmTypeJSON))

	w.Header().Set("Content-Type", "Application/json")
	w.Write(mvmTypeJSON)
}

//GetAllActionsTypes gestion du select des types d'actions
func (a *AcquisitionService) GetAllActionsTypes(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=tsaplaval")

	defer db.Close()
	fmt.Println(err)

	actionTypes := []TypeAction{}
	db.Find(&actionTypes)

	actionTypesJSON, _ := json.Marshal(actionTypes)
	fmt.Println(string(actionTypesJSON))

	w.Header().Set("Content-Type", "Application/json")
	w.Write(actionTypesJSON)
}

func (a *AcquisitionService) PostActionType(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=tsaplaval")

	defer db.Close()
	fmt.Println(err)

	body, err := ioutil.ReadAll(r.Body)
	fmt.Printf("-----------------------")
	fmt.Println(body)
	fmt.Printf("-----------------------")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	var newActionType TypeAction

	err = json.Unmarshal(body, &newActionType)

	fmt.Println(err)

	if err != nil {
		panic(err)
	}
	if db.NewRecord(newActionType) {
		db.Create(&newActionType)
		db.NewRecord(newActionType) // => return `false` after `user` created
	} else {
		fmt.Println("erreur")
	}
	w.Header().Set("Content-Type", "application/json")

}

/*func (a *AcquisitionService) DeleteActionType(w http.ResponseWriter, r *http.Request, pId *int) {

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=tsaplaval")

	defer db.Close()
	fmt.Println(err)

	if err != nil {
		db.Where("Id LIKE ?", pId).Delete(TypeAction{})
	} else {
		log.Fatal(err)
	}
}

func (a *AcquisitionService) UpdateActionType(w http.ResponseWriter, r *http.Request, pId *int, name *string, description *string) {

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=tsaplaval")
	defer db.Close()
	fmt.Println(err)

	actionUpdate := TypeAction{ID: pId, Nom: &name, Description: description}
}*/
