//
// Fichier     : coaches_test.go
// Développeur : Mehdi Laribi
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

// TestPostNewCoach: test de la methode POST
func TestPostNewCoach(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Fname": "MEHEHHEHEHEDi", 
			"Lname": "Lariihhhiibi", 
			"Actif": "true", 
			"Email": "Mehdi@hotmale.com"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coaches/addcoach", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var l api.Coaches
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Error(err)
	}

	rmID = fmt.Sprintf("%d", l.ID)

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

	if l.Lname != "Lariihhhiibi" {
		t.Error("LName expected: ", l.Lname)
	}

	if l.Fname != "MEHEHHEHEHEDi" {
		t.Error("FName expected: ", l.Fname)
	}

	if l.Actif != "true" {
		t.Error("Actif expected: ", l.Actif)
	}

	if l.Email != "Mehdi@hotmale.com" {
		t.Error("Actif expected: ", l.Actif)
	}
}

//TestGetCoaches : test de la methode GET
func TestGetCoaches(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/coaches", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var co []api.Coaches
	err = json.Unmarshal(bodyBuffer, &co)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	// On s'assure qu'il y ait au moins un coach
	if len(co) <= 0 {
		t.Errorf("Number of Coaches expected: %d", len(co))
	}

}

//TestUpdateCoach : test de la methode PUT
func TestUpdateCoach(t *testing.T) {
	reader = strings.NewReader(`{"Fname": "m", "Lname": "l", "Actif": "false", "Email": "Mehdi@hotmale.com"}`)
	request, err := http.NewRequest("PUT", baseURL+"/api/coaches/updatecoach"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TestModifierCoach test la modification de l'équipe créée plus haut
func TestModifierCoach(t *testing.T) {
	reader = strings.NewReader(`{"Fname": "Mehdi", "Lname": "Laribi", "Actif": "false", "Email": "Mehdi@hotmale.com"}`)

	// rmID est utilisé ici pour permettre la modification créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/coaches/editcoach/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te api.Coaches
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if te.Fname != "Mehdi" {
		t.Error("New name expected: ", te.Fname)
	}

	if te.Lname != "Laribi" {
		t.Error("Last name expected: ", te.Lname)
	}
}

// TestCreerCoacheVide test que créer un coach sans
// information. Doit retourner une erreur.
func TestCreerCoacheVide(t *testing.T) {
	reader = strings.NewReader(``)

	request, err := http.NewRequest("POST", baseURL+"/coaches/addcoach", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}
