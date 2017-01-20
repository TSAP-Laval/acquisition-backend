package api

import (
	"fmt"
	"net/http"
  
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

func (a *AcquisitionService) FaireBD(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsap_dev sslmode=disable password=")
=======
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)

>>>>>>> 0b02044ad79ea84a15a4314b264bdf8701d704f9
	defer db.Close()
	fmt.Println(err)

	db.AutoMigrate(&TypeAction{})
	db.AutoMigrate(&Sport{})
	db.AutoMigrate(&Niveau{})
	db.AutoMigrate(&Entraineur{})
	db.AutoMigrate(&Joueur{})
	db.AutoMigrate(&Equipe{})
	db.AutoMigrate(&Zone{})
	db.AutoMigrate(&Saison{})
	db.AutoMigrate(&Lieu{})
	db.AutoMigrate(&Video{})
	db.AutoMigrate(&Partie{})
	db.AutoMigrate(&Action{})
}
func (a *AcquisitionService) Remplir(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsap_dev sslmode=disable password=")
=======
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
>>>>>>> 0b02044ad79ea84a15a4314b264bdf8701d704f9
	defer db.Close()
	fmt.Println(err)
	user := TypeAction{Nom: "PO", Description: "Passe offensive"}
	if db.NewRecord(user) {
		fmt.Println("Action")
		db.Create(&user)
		db.NewRecord(user) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}
	entraineur := Entraineur{Prenom: "alex", Nom: "Des", Email: "alex@hotmail.com", PassHash: "test", Token: "test"}
	if db.NewRecord(entraineur) {
		fmt.Println("Test")
		db.Create(&entraineur)
		db.NewRecord(entraineur) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}
	joueur := Joueur{Prenom: "alex", Nom: "Des", Numero: 1, Email: "alex@hotmail.com", PassHash: "test", TokenInvitation: "test", TokenReinitialisation: "test", TokenConnexion: "test"}
	if db.NewRecord(joueur) {
		fmt.Println("Test")
		db.Create(&joueur)
		db.NewRecord(joueur) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}
	Sport := Sport{Nom: "soccer"}
	if db.NewRecord(Sport) {
		fmt.Println("Test")
		db.Create(&Sport)
		db.NewRecord(Sport) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}
	Niveau := Niveau{Nom: "AA"}
	if db.NewRecord(Niveau) {
		fmt.Println("Test")
		db.Create(&Niveau)
		db.NewRecord(Niveau) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}
	Zone := Zone{Nom: "off"}
	if db.NewRecord(Zone) {
		fmt.Println("Test")
		db.Create(&Zone)
		db.NewRecord(Zone) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}

	Action := Action{TypeActionID: 1, ActionPositive: true, ZoneID: 1, PartieID: 1, X1: 0, Y1: 0, X2: 0, Y2: 0, Temps: 10, PointageMaison: 0, PointageAdverse: 0, JoueurID: 1}
	if db.NewRecord(Action) {
		fmt.Println("Test")
		db.Create(&Action)
		db.NewRecord(Action) // => return `false` after `user` created
	} else {
		fmt.Println("Test22")
	}
}
