//
// Fichier     : joueurs_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// de trop de choses ?.
//

package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/TSAP-Laval/acquisition-backend/api"
)

// TestGetJoueursErrBD test la récupération de tous les joueurs
// avec erreur de connexion à la base de données
func TestGetJoueursErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/joueurs", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetJoueurs test la récupération de tous les joueurs
func TestGetJoueurs(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/joueurs", reader)

	if err != nil {
		t.Error(err)
	}

	GetRequestHandler(request, t)
}

// TestCreerJoueur test la création d'un joueur
func TestCreerJoueur(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Lname": "Test",
			"Fname": "Test",
			"Number": 55,
			"Email" : "test@test.ca",
			"PassHash" : "test123" ,
			"EquipeID": "1"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/joueurs", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var p api.Players
	err = json.Unmarshal(bodyBuffer, &p)
	if err != nil {
		t.Error(err)
	}

	rmID = fmt.Sprintf("%d", p.ID)

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TestModifJoueurErrBD test la modification d'un joueur avec erreur
// de connexion à la base de données
func TestModifJoueurErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Lname": "Test",
			"Fname": "Test",
			"Number": 8,
			"Email": "test@test.com",
			"PassHash": "test123",
			"EquipeID": "2"
		}`)

	// rmID est utilisé, ici pour permettre la modification du joueur tout juste créé
	request, err := http.NewRequest("PUT", baseURL+"/api/joueurs/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestModifJoueurErr test la modification d'un joueur avec une erreur dans le JSON
func TestModifJoueurErr(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Lname" Test",
			"Fname": "Test",
			"Number": 8,
			"Email": "test@test.com",
			"PassHash": "test123",
			"EquipeID": "2"
		}`)

	// rmID est utilisé, ici pour permettre la modification du joueur tout juste créé
	request, err := http.NewRequest("PUT", baseURL+"/api/joueurs/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifJoueur test la modification d'un joueur
func TestModifJoueur(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Lname": "Test",
			"Fname": "Test",
			"Number": 8,
			"Email": "test@test.com",
			"PassHash": "test123",
			"EquipeID": "2"
		}`)

	// rmID est utilisé, ici pour permettre la modification du joueur tout juste créé
	request, err := http.NewRequest("PUT", baseURL+"/api/joueurs/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var p api.Players
	err = json.Unmarshal(bodyBuffer, &p)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

	// Test des valeurs de retour
	if p.Email != "test@test.com" {
		t.Error("Success expected: ", p.Email)
	}

	if p.Number != 8 {
		t.Errorf("Success expected: %d", p.Number)
	}
}

// TestJoueurOptions test une requête OPTIONS pour un joueur
func TestJoueurOptions(t *testing.T) {
	reader = strings.NewReader("")

	// rmID est utilisé, ici pour permettre la modification du joueur tout juste créé
	request, err := http.NewRequest("OPTIONS", baseURL+"/api/joueurs/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	GetRequestHandler(request, t)
}

// TestJoueurOptionsAgain test une requête OPTIONS pour un joueur
func TestJoueurOptionsAgain(t *testing.T) {
	reader = strings.NewReader("")

	request, err := http.NewRequest("OPTIONS", baseURL+"/api/joueurs", reader)

	if err != nil {
		t.Error(err)
	}

	GetRequestHandler(request, t)
}
