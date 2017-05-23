//
// Fichier     : seeders.go
// Développeur : Alexandre Deschêne et Laurent Leclerc Poulin
//
// Permet de créer la base de données et d'y ajouter des information
//

package api

import (
	"net/http"

	//Import DB driver
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// FaireBD crée la base de données à partie du modèle de données (structures.go)
func (a *AcquisitionService) FaireBD(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()

	if err != nil {
		a.ErrorHandler(w, err)
		return
	}

	db.DropTableIfExists("admins")
	db.DropTableIfExists("reception_types")
	db.DropTableIfExists("actions")
	db.DropTableIfExists("videos")
	db.DropTableIfExists("player_position_game_team")
	db.DropTableIfExists("games")
	db.DropTableIfExists("player_team")
	db.DropTableIfExists("players")
	db.DropTableIfExists("coach_teams")
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

	db.AutoMigrate(&ReceptionType{})
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
		a.ErrorHandler(w, err)
		return
	}

	reception := ReceptionType{Name: "Ballon conquis aérien"}
	db.Create(&reception)

	reception2 := ReceptionType{Name: "Ballon conquis sol"}
	db.Create(&reception2)

	reception3 := ReceptionType{Name: "Ballon reçu sol"}
	db.Create(&reception3)

	reception4 := ReceptionType{Name: "Ballon reçu aérien"}
	db.Create(&reception4)

	reception5 := ReceptionType{Name: "Ballon reçu sur faute de l’adversaire"}
	db.Create(&reception5)

	reception6 := ReceptionType{Name: "Ballon reçu sur remise en jeu"}
	db.Create(&reception6)

	reception7 := ReceptionType{Name: "Passe reçue au sol"}
	db.Create(&reception7)

	reception8 := ReceptionType{Name: "Passe aérienne reçue "}
	db.Create(&reception8)

	user := ActionsType{Description: "Tir au but cadré", TypeAction: "reception et action", Name: "Tir au but cadré"}
	db.Create(&user)

	userTirAuBut := ActionsType{Description: "Tir au but non-cadré", TypeAction: "reception et action", Name: "Tir au but non-cadré"}
	db.Create(&userTirAuBut)

	action := ActionsType{Description: "Perte directe sur contrôle", TypeAction: "balle perdu", Name: "Perte directe sur contrôle"}
	db.Create(&action)

	action3 := ActionsType{Description: "Perte directe sur passe tentée", TypeAction: "passe incomplete", Name: "Perte directe sur passe tentée"}
	db.Create(&action3)

	action4 := ActionsType{Description: "Perte directe autres(faute, etc)", TypeAction: "balle perdu", Name: "Perte directe autres(faute, etc)"}
	db.Create(&action4)

	action5 := ActionsType{Description: "Passe offensive positive (dans la course du joueur)", TypeAction: "reception et action", Name: "Passe offensive positive (dans la course du joueur)"}
	db.Create(&action5)

	action6 := ActionsType{Description: "Passe offensive négative (joueur doit modifier sa course)", TypeAction: "reception et action", Name: "Passe offensive négative (joueur doit modifier sa course)"}
	db.Create(&action6)

	action7 := ActionsType{Description: "Dégagement réussi", TypeAction: "passe incomplete", Name: "Dégagement réussi"}
	db.Create(&action7)

	action8 := ActionsType{Description: "Passe neutre", TypeAction: "reception et action", Name: "Passe neutre"}
	db.Create(&action8)

	coach := Coaches{Fname: "alex", Lname: "Des", Email: "alex@hotmail.com", PassHash: "test"}
	db.Create(&coach)

	player := Players{Fname: "alex", Lname: "Des", Number: 1, Email: "alex@hotmail.com", PassHash: "test"}
	db.Create(&player)

	season := Seasons{Years: "2017-2018"}
	db.Create(&season)

	sport := Sports{Name: "soccer"}
	db.Create(&sport)

	category := Categories{Name: "AA"}
	db.Create(&category)

	zone := Zones{Name: "off"}
	db.Create(&zone)

	fieldType1 := FieldTypes{Type: "Synthétique", Description: "Terrain synthétique..."}
	db.Create(&fieldType1)

	fieldType2 := FieldTypes{Type: "Gazon", Description: "Terrain extérieur..."}
	db.Create(&fieldType2)

	location1 := Locations{Name: "SSF", City: "Saint-Augustin-de-Desmaures", Address: "4900 Rue Saint-Félix", IsInside: false, FieldTypesID: int(fieldType1.ID)}
	db.Create(&location1)

	location2 := Locations{Name: "Stade Leclerc", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	db.Create(&location2)

	location3 := Locations{Name: "Terrain univers", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	db.Create(&location3)

	location4 := Locations{Name: "Stade Bidule", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: false, FieldTypesID: int(fieldType2.ID)}
	db.Create(&location4)

	location5 := Locations{Name: "Stade Chnoubouc", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	db.Create(&location5)

	location6 := Locations{Name: "Terrains PGR", City: "Saint-Augustin-de-Desmaures", Address: "QC G3A 0C3", IsInside: true, FieldTypesID: int(fieldType1.ID)}
	db.Create(&location6)

	equipe1 := Teams{Name: "Lions", City: "Quebec", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "M"}
	db.Create(&equipe1)

	equipe2 := Teams{Name: "Loup", City: "Vancouver", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "M"}
	db.Create(&equipe2)

	equipe3 := Teams{Name: "Tigres", City: "Montreal", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "F"}
	db.Create(&equipe3)

	equipe4 := Teams{Name: "Ligres", City: "Trois-Rivières", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "M"}
	db.Create(&equipe4)

	equipe5 := Teams{Name: "Tattoo", City: "Rivière-du-loup", SportID: 1, CategoryID: 1, SeasonID: 1, Sexe: "F"}
	db.Create(&equipe5)

	pass, _ := bcrypt.GenerateFromPassword([]byte("aaaaa"), bcrypt.DefaultCost)
	admin := Admins{Email: "admin@admin.ca", PassHash: string(pass)}
	db.Create(&admin)

	// Admin avec un token expiré (pour les tests seulement)
	pass, _ = bcrypt.GenerateFromPassword([]byte("aaaaa"), bcrypt.DefaultCost)
	badAdmin := Admins{Email: "mauvais@mauvais.ca", PassHash: string(pass),
		TokenLogin: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNDk0Njk0OTU0fQ.TBukRueijLUla7hejpR064CERMXJy3CRbWWhPQPQ5fY"}
	db.Create(&badAdmin)

	Uneaction := Actions{ActionTypeID: 1, ZoneID: 1, GameID: 1, X1: 0, Y1: 0, X2: 0, Y2: 0, X3: 0, Y3: 0, Time: 10.5, HomeScore: 0, GuestScore: 0, PlayerID: 1}
	db.Create(&Uneaction)

	Message(w, "", http.StatusCreated)
}
