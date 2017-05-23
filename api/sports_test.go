//
// TEST
//
// Fichier     : sports_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de tester les interractions sur un sport.
//

package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/TSAP-Laval/acquisition-backend/api"
)

// TestGetSportErrBD permet de récupérer les sports
// avec erreur de connexion à la base de données
func TestGetSportsErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/sports", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetSports permet de récupérer tous les sports
func TestGetSports(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/sports", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var sp []api.Sports
	err = json.Unmarshal(bodyBuffer, &sp)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	// Test simplement qu'il y a au moins un sport de créé
	if len(sp) < 1 {
		t.Error("Number of sports expected: ", len(sp))
	}
}
