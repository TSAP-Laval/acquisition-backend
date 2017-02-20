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

func TestGetTerrains(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestCreerTerrain(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "LE terrain", "Ville": "Quebec", "Adresse": "1231 une rue"}`)
	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)
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

func TestCreerTerrainErrEmpty(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "UN terrain", "Ville": "Quebec", "Adresse": ""}`)
	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestCreerTerrainErrExiste(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "LE terrain", "Ville": "Quebec", "Adresse": ""}`)
	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetTerrain(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestModifierTerrain(t *testing.T) {
	reader = strings.NewReader(`{"Nom": "LE terrain", "Ville": "Montreal", "Adresse": ""}`)
	request, err := http.NewRequest("PUT", baseURL+"/api/terrains/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetTerrainModifie(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// Buffer the body
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	t.Logf("Res: --> %s\n\n", bodyBuffer)

	// Ve receive an array of Lieu
	l := []api.Lieu{}
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Logf("ERR: --> %s\n\n", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	// We only look at the first element of the array
	if l[0].Ville != "Montreal" {
		t.Errorf("City expected: %s", l[0].Ville)
	}
	if l[0].Adresse != "1231 une rue" {
		t.Errorf("Address expected: %s", l[0].Adresse)
	}
}

func TestSupprimerTerrain(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/terrains/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestGetTerrainErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains/LE", reader)
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

func TestSupprimerTerrainSupprime(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/terrains/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}
