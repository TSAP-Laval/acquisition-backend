//
// Fichier     : gestion_test.go
// Développeur : ?
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

// TODO : Aucun commentaire sur les fonctions

// TODO : Ne test même pas les valeurs de retour ?
func TestGetSaison(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/saisons", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TODO : Ne test même pas les valeurs de retour ?
func TestCreerSaison(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Years": "2000"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/saisons", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var s api.Seasons
	err = json.Unmarshal(bodyBuffer, &s)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TODO : Ne test même pas les valeurs de retour ?
func TestGetSports(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/sports", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TODO : Ne test même pas les valeurs de retour ?
func TestGetNiveau(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/niveaux", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TODO : Ne test même pas les valeurs de retour ?
func TestGetJoueurs(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/joueurs", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TODO : Ne test que la création et la modification d'un joueurs ?
// TODO : Un joueur peut être dans plusieurs équipes à la fois, il ne
//        faudrait donc pas lui affecter un équipe de cette façon...
//        D'ailleurs, il faudrait mettre à jour ces tests
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
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

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

func TestModifJoueur(t *testing.T) {
	reader = strings.NewReader(
		`{
			"ID": "15",
			"Lname": "Test",
			"Fname": "Test",
			"Number": 8,
			"Email": "test@test.com",
			"PassHash": "test123",
			"EquipeID": "2"
		}`)

	// rmID est utilisé, ici pour permettre la modification du joueur tout juste créé
	request, err := http.NewRequest("PUT", baseURL+"/api/joueurs/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

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
