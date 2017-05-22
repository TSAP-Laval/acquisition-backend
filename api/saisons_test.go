//
// TEST
//
// Fichier     : saisons_test.go
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

// TestGetSaisonsErrBD permet de récupérer les saisons
// avec erreur de connexion à la base de données
func TestGetSaisonsErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/saisons", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestGetSaisons permet de récupérer toutes les saisons
func TestGetSaisons(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/saisons", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var sa []api.Seasons
	err = json.Unmarshal(bodyBuffer, &sa)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	// Test simplement qu'il y a au moins une saison de créé
	if len(sa) < 1 {
		t.Error("Number of saison expected: ", len(sa))
	}
}

// TestCreerSaisonErrBD test la création d'une saison
// avec erreur de connexion à la base de données
func TestCreerSaisonErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader(
		`{
			"Years": "2015-2016"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/saisons", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestCreerSaisonErr test la création d'une saison
// avec erreur de dans le JSON
func TestCreerSaisonErr(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader(
		`{
			"Years "2015-2016"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/saisons", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestCreerSaison test la création d'une saison
func TestCreerSaison(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Years": "2015-2016" 
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/saisons", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var s api.Seasons
	err = json.Unmarshal(bodyBuffer, &s)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if s.Years != "2015-2016" {
		t.Error("Year expected: ", s.Years)
	}
}

// TestCreerSaisonExiste test la création d'une saison
// qui existe déjà dans la base de données
func TestCreerSaisonExiste(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Years": "2015-2016"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/saisons", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "La saison entrée existe déjà !") {
		t.Error("Error expected : ", me.Err)
	}
}
