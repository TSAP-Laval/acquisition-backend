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

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t Saison
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t.ID)
	if db.NewRecord(t) {
		db.Create(&t)
		db.NewRecord(t)
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	} else {
		fmt.Println("Test22")
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}

}
func (a *AcquisitionService) PostTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t Equipe
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t.ID)
	if db.NewRecord(t) {
		x := Sport{}
		db.First(&x, t.SportID)
		t.Sport = x
		Niv := Niveau{}
		db.First(&Niv, t.NiveauID)
		t.Niveau = Niv

		SportJSON, _ := json.Marshal(t)
		fmt.Println(string(SportJSON))
		db.Create(&t)
		SportJSON2, _ := json.Marshal(t)
		fmt.Println(string(SportJSON2))

		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	} else {
		fmt.Println("Test22")
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}

}
func (a *AcquisitionService) PostJoueur(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t Joueur
	var dat map[string]interface{}
	err = json.Unmarshal(body, &t)
	err = json.Unmarshal(body, &dat)
	num := dat["EquipeID"]
	Team := Equipe{}
	db.First(&Team, num)
	t.Equipes = append(t.Equipes, Team)
	db.Model(&Team).Association("Equipes").Append(t)
	fmt.Println(num)
	if err != nil {
		panic(err)
	}
	log.Println(t.ID)
	if db.NewRecord(t) {
		db.Create(&t)
		db.NewRecord(t)
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	} else {
		fmt.Println("Test22")
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}

}
func (a *AcquisitionService) PostEquipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t Equipe
	err = json.Unmarshal(body, &t)

	if err != nil {
		panic(err)
	}
	log.Println(t.ID)
	if db.NewRecord(t) {

		db.Create(&t)
		db.NewRecord(t)
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	} else {
		fmt.Println("Test22")
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("erreur"))
	}

}
func (a *AcquisitionService) GetSeasons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)
	strucSaison := []Saison{}
	db.Find(&strucSaison)

	SaisonJSON, _ := json.Marshal(strucSaison)
	fmt.Println(string(SaisonJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SaisonJSON)
}
func (a *AcquisitionService) GetSports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)
	strucSport := []Sport{}
	db.Find(&strucSport)

	SportJSON, _ := json.Marshal(strucSport)
	fmt.Println(string(SportJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SportJSON)
}

func (a *AcquisitionService) GetUnNiveauTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)
	strucSport := Equipe{}
	db.Last(&strucSport)

	SportJSON, _ := json.Marshal(strucSport)
	fmt.Println(string(SportJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SportJSON)
}
func (a *AcquisitionService) GetNiveau(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)
	strucNiveau := []Niveau{}
	db.Find(&strucNiveau)

	NiveauJSON, _ := json.Marshal(strucNiveau)
	fmt.Println(string(NiveauJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(NiveauJSON)
}
func (a *AcquisitionService) GetEquipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")

	defer db.Close()
	fmt.Println(err)
	strucEquipe := []Equipe{}
	db.Find(&strucEquipe)

	EquipeJSON, _ := json.Marshal(strucEquipe)
	fmt.Println(string(EquipeJSON))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(EquipeJSON)
}
