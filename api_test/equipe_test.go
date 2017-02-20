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

func TestSeed(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/seeders", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetEquipes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestCreerEquipe(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "Lequipe", "Ville": "Quebec", "NiveauID": 1, "SportID": 1}`)
	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	// Buffer the body
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	t.Logf("Res: --> %s\n\n", bodyBuffer)

	var l api.Lieu
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Logf("ERR: --> %s\n\n", err)
	}

	rmID = fmt.Sprintf("%d", l.ID)
	t.Logf("ID: --> %s\n\n", rmID)

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestCreerEquipeErrEmpty(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "UNE equipe", "":"", "NiveauID": 1, "SportID": 1}`)
	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestCreerEquipeErrExiste(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "Lequipe", "Ville": "Quebec", "NiveauID": 1, "SportID": 1}`)
	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 401 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetEquipe(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestModifierEquipe(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "LE equipe", "Ville": "Montreal", "NiveauID": 1, "SportID": 1}`)
	request, err := http.NewRequest("PUT", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetEquipeModifie(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// Buffer the body
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	t.Logf("Res: --> %s\n\n", bodyBuffer)

	// Ve receive an array of Equipe
	e := []api.Equipe{}
	err = json.Unmarshal(bodyBuffer, &e)
	if err != nil {
		t.Logf("ERR: --> %s\n\n", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	// We only look at the first element of the array
	if e[0].Ville != "Montreal" {
		t.Errorf("City expected: %s", e[0].Ville)
	}
}

func TestSupprimerEquipe(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetEquipeErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	if len(bodyBuffer) > 20 {
		t.Errorf("Body response length expected: %d", len(bodyBuffer))
	}
}

func TestSupprimerEquipeSupprime(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}
