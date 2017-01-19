package api

import (
	"fmt"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

func (a *AcquisitionService) FaireBD(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=tsaplaval")
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
	fmt.Println(a.config.DatabaseDriver)
	fmt.Println(a.config.ConnectionString)
	db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
	defer db.Close()
	fmt.Println(err)
	user := TypeAction{Nom: "passe", Description: "une passe"}
	if db.NewRecord(user) {
		fmt.Println("Test")
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
}
