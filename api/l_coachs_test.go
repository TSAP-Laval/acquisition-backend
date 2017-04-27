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

	//Mehdi Fix
	_ "github.com/lib/pq"
)

// TestPostNewCoach: test de la methode POST
func TestPostNewCoach(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Fname": "MEHEHHEHEHEDi", 
			"Lname": "Lariihhhiibi", 
			"Actif": "true", 
			"Email": "Mehdi@hotmale.com",
			"TeamsIDs": "1,2",
			"SeasonID": 5
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
