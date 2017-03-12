package api

import (
	"fmt"
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// FaireBD crée la base de données à partie du modèle de données (structures.go)
func (a *AcquisitionService) FaireBD(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		fmt.Print("ERROR : ")
		fmt.Println(err)
	}

	db.DropTableIfExists(&Admins{})
	db.DropTableIfExists(&Actions{})
	db.DropTableIfExists(&Games{})
	db.DropTableIfExists(&Teams{})
	db.DropTableIfExists(&PlayersTeam{})
	db.DropTableIfExists(&CoachTeam{})
	db.DropTableIfExists(&Zones{})
	db.DropTableIfExists(&Sports{})
	db.DropTableIfExists(&Players{})
	db.DropTableIfExists(&Locations{})
	db.DropTableIfExists(&Categories{})
	db.DropTableIfExists(&Coaches{})
	db.DropTableIfExists(&ActionsType{})
	db.DropTableIfExists(&Seasons{})
	db.DropTableIfExists(&Positions{})
	db.DropTableIfExists(&MovementsType{})
	db.DropTableIfExists(&PlayerPositionGameTeam{})
	db.DropTableIfExists(&Videos{})
	db.DropTableIfExists(&Metrics{})

	db.AutoMigrate(&Admins{})
	db.AutoMigrate(&Seasons{})
	db.AutoMigrate(&Sports{})
	db.AutoMigrate(&Categories{})
	db.AutoMigrate(&Teams{})
	db.AutoMigrate(&Players{})
	db.AutoMigrate(&Locations{})
	db.AutoMigrate(&Videos{})
	db.AutoMigrate(&Games{})
	db.AutoMigrate(&Positions{})
	db.AutoMigrate(&PlayerPositionGameTeam{})
	db.AutoMigrate(&Zones{})
	db.AutoMigrate(&MovementsType{})
	db.AutoMigrate(&ActionsType{})
	db.AutoMigrate(&Actions{})
	db.AutoMigrate(&PlayersTeam{})
	db.AutoMigrate(&Coaches{})
	db.AutoMigrate(&CoachTeam{})
	db.AutoMigrate(&Metrics{})
}

// RemplirBD permet de `seeder` la base de données avec des données `hard-codées`
func (a *AcquisitionService) RemplirBD(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		fmt.Println("ERROR : ")
		fmt.Println(err)
	}

	user := ActionsType{Name: "PO", Description: "Passe offensive"}
	if db.NewRecord(user) {
		db.Create(&user)
	}

	coach := Coaches{Fname: "alex", Lname: "Des", Email: "alex@hotmail.com", PassHash: "test"}
	if db.NewRecord(coach) {
		db.Create(&coach)
	}

	player := Players{Fname: "alex", Lname: "Des", Number: 1, Email: "alex@hotmail.com", PassHash: "test"}
	if db.NewRecord(player) {
		db.Create(&player)
	}

	season := Seasons{Years: "1997-1998"}
	if db.NewRecord(season) {
		db.Create(&season)
	}

	sport := Sports{Name: "soccer"}
	if db.NewRecord(sport) {
		db.Create(&sport)
	}

	category := Categories{Name: "AA"}
	if db.NewRecord(category) {
		db.Create(&category)
	}

	zone := Zones{Name: "off"}
	if db.NewRecord(zone) {
		db.Create(&zone)
	}
  
	location1 := Locations{Name: "SSF", City: "St-Augustin", Address: "1223 rue Truc"}
	if db.NewRecord(location1) {
		db.Create(&location1)
	}

	location2 := Locations{Name: "Stade Leclerc", City: "St-Augustin", Address: "1224 rue Leclerc"}
	if db.NewRecord(location2) {
		db.Create(&location2)
	}

	equipe1 := Teams{Name: "Lions", City: "Quebec", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe1) {
		db.Create(&equipe1)
	}

	equipe2 := Teams{Name: "Tigres", City: "Montreal", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe2) {
		db.Create(&equipe2)
	}

	equipe3 := Teams{Name: "Ligres", City: "Trois-Rivières", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe3) {
		db.Create(&equipe3)
	}

	action := Actions{ActionTypeID: 1, IsPositive: true, ZoneID: 1, GameID: 1, X1: 0, Y1: 0, X2: 0, Y2: 0, Time: 10, HomeScore: 0, GuestScore: 0, PlayerID: 1}
	if db.NewRecord(action) {
		db.Create(&action)
	}
}
