//
// Fichier     : seeders.go
// Développeur : Alexandre Deschêne et Laurent Leclerc Poulin
//
// Permet de créer la base de données et d'y ajouter des information
//

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
	db.DropTableIfExists("actions_types")
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

	user := ActionsType{Description: "Passe offensive", TypeAction: "reception et action", Name: "Passe offensive"}
	if db.NewRecord(user) {
		db.Create(&user)
	}

	action := ActionsType{Description: "Reception suivi d'une perte de ballon", TypeAction: "balle perdu", Name: "Balle perdu"}
	if db.NewRecord(action) {
		db.Create(&action)
	}
	action3 := ActionsType{Description: "passe defensif", TypeAction: "reception et action", Name: "passe defensif"}
	if db.NewRecord(action3) {
		db.Create(&action3)
	}
	action4 := ActionsType{Description: "Dégagement gardien", TypeAction: "reception et action", Name: "Dégagement gardien"}
	if db.NewRecord(action4) {
		db.Create(&action4)
	}
	action5 := ActionsType{Description: "faute", TypeAction: "balle perdu", Name: "faute"}
	if db.NewRecord(action5) {
		db.Create(&action5)
	}
	action6 := ActionsType{Description: "tir arreter", TypeAction: "reception et action", Name: "tir arreter"}
	if db.NewRecord(action6) {
		db.Create(&action6)
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

	fieldType1 := FieldTypes{Type: "Synthétique", Description: "Terrain synthétique..."}
	if db.NewRecord(fieldType1) {
		db.Create(&fieldType1)
	}

	fieldType2 := FieldTypes{Type: "Gazon", Description: "Terrain extérieur..."}
	if db.NewRecord(fieldType2) {
		db.Create(&fieldType2)
	}

	location1 := Locations{Name: "SSF", City: "Saint-Augustin-de-Desmaures", Address: "4900 Rue Saint-Félix", IsInside: false, FieldTypesID: int(fieldType1.ID)}
	if db.NewRecord(location1) {
		db.Create(&location1)
	}

	location2 := Locations{Name: "Stade Leclerc", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	if db.NewRecord(location2) {
		db.Create(&location2)
	}
	location3 := Locations{Name: "Terrain univers", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	if db.NewRecord(location3) {
		db.Create(&location3)
	}

	location4 := Locations{Name: "Stade Bidule", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: false, FieldTypesID: int(fieldType2.ID)}
	if db.NewRecord(location4) {
		db.Create(&location4)
	}

	location5 := Locations{Name: "Stade Chnoubouc", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	if db.NewRecord(location5) {
		db.Create(&location5)
	}

	location6 := Locations{Name: "Terrains PGR", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	if db.NewRecord(location6) {
		db.Create(&location6)
	}

	equipe1 := Teams{Name: "Lions", City: "Quebec", SportID: 1, CategoryID: 1}
	if db.NewRecord(equipe1) {
		db.Create(&equipe1)
	}

	equipe2 := Teams{Name: "Loup", City: "Vancouver", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "M"}
	if db.NewRecord(equipe2) {
		db.Create(&equipe2)
	}

	equipe3 := Teams{Name: "Tigres", City: "Montreal", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "F"}
	if db.NewRecord(equipe3) {
		db.Create(&equipe3)
	}

	equipe4 := Teams{Name: "Ligres", City: "Trois-Rivières", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "M"}
	if db.NewRecord(equipe4) {
		db.Create(&equipe4)
	}

	equipe5 := Teams{Name: "Tatoo", City: "Rivière-du-loup", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "F"}
	if db.NewRecord(equipe5) {
		db.Create(&equipe5)
	}

	Uneaction := Actions{ActionTypeID: 1, ZoneID: 1, GameID: 1, X1: 0, Y1: 0, X2: 0, Y2: 0, X3: 0, Y3: 0, Time: 10, HomeScore: 0, GuestScore: 0, PlayerID: 1}
	if db.NewRecord(Uneaction) {
		db.Create(&Uneaction)
	}
}
