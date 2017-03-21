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
		fmt.Println("ERROR : ")
		fmt.Println(err)
	}

	db.DropTableIfExists("admins")
	db.DropTableIfExists("actions")
	db.DropTableIfExists("videos")
	db.DropTableIfExists("player_position_game_team")
	db.DropTableIfExists("games")
	db.DropTableIfExists("player_team")
	db.DropTableIfExists("players")
	db.DropTableIfExists("coach_team")
	db.DropTableIfExists("metrics")
	db.DropTableIfExists("teams")
	db.DropTableIfExists("sports")
	db.DropTableIfExists("categories")
	db.DropTableIfExists("zones")
	db.DropTableIfExists("locations")
	db.DropTableIfExists("field_types")
	db.DropTableIfExists("coaches")
	db.DropTableIfExists("actions_type")
	db.DropTableIfExists("seasons")
	db.DropTableIfExists("positions")
	db.DropTableIfExists("movements_type")

	db.AutoMigrate(&Admins{})
	db.AutoMigrate(&Seasons{})
	db.AutoMigrate(&Sports{})
	db.AutoMigrate(&Categories{})
	db.AutoMigrate(&Teams{})
	db.AutoMigrate(&Players{})
	db.AutoMigrate(&Locations{})
	db.AutoMigrate(&FieldTypes{})
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

	user := ActionsType{Description: "Passe offensive", Acquisition: "Positive", Separation: "Positive"}
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

	season := Seasons{Years: "2017-2018"}
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

	location1 := Locations{Name: "SSF", City: "Saint-Augustin-de-Desmaures", Address: "4900 Rue Saint-Félix", IsInside: false}
	if db.NewRecord(location1) {
		db.Create(&location1)
	}

	location2 := Locations{Name: "Stade Leclerc", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true}
	if db.NewRecord(location2) {
		db.Create(&location2)
	}

	equipe1 := Teams{Name: "Lions", City: "Quebec", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe1) {
		db.Create(&equipe1)
	}

	equipe2 := Teams{Name: "Loup", City: "Vancouver", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe2) {
		db.Create(&equipe2)
	}

	equipe3 := Teams{Name: "Tigres", City: "Montreal", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe3) {
		db.Create(&equipe3)
	}

	equipe4 := Teams{Name: "Ligres", City: "Trois-Rivières", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe4) {
		db.Create(&equipe4)
	}

	equipe5 := Teams{Name: "Tatoo", City: "Rivière-du-loup", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe5) {
		db.Create(&equipe5)
	}

	action := Actions{ActionTypeID: 1, IsPositive: true, ZoneID: 1, GameID: 1, X1: 0, Y1: 0, X2: 0, Y2: 0, Time: 10, HomeScore: 0, GuestScore: 0, PlayerID: 1}
	if db.NewRecord(action) {
		db.Create(&action)
	}
}
