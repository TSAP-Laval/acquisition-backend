//
// Fichier     : coachs_test.go
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

// PostNewCoach: test de la methode POST
func PostNewCoach(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Fname": "MEHEHHEHEHEDi", 
			"Lname": "Lariihhhiibi", 
			"Actif": "true", 
			"Email": "Mehdi@hotmale.com"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coachs/addcoach", reader)
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
}

//GetCoaches : test de la methode GET
func GetCoaches(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/coachs", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode)
	}
}

//UpdateCoach : test de la methode PUT
func UpdateCoach(t *testing.T) {
	reader = strings.NewReader(`{"Fname": "m", "Lname": "l", "Actif": "false", "Email": "Mehdi@hotmale.com"}`)
	request, err := http.NewRequest("PUT", baseURL+"/api/coachs/updatecoach"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}
