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

func IsDbWorking(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/bd", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Error
	}

	if res.StatusCode != 200 {
		t.Errorf("Is Success: %d", res.StatusCode)
	}
}

func IsSeeded(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/seed", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Error
	}

	if res.StatusCode != 200 {
		t.Errorf("Is seeded: %d", res.StatusCode)
	}
}

func PostNewCoach(t *testing.T) {
	reader = strings.NewReader(`{"Fname": "MEHEHHEHEHEDi", "Lname": "Lariihhhiibi", "Actif": "true", "Email": "Mehdi@hotmale.com"}`)
	request, err := http.NewRequest("POST", baseURL+"/api/coachs/addcoach", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	// Buffer the body
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	t.Logf("Res: --> %s\n\n", bodyBuffer)

	var l api.Coaches
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

func GetCoachs(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/coachs", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode)
	}
}

/*func UpdateCoach(t *testing.T) {
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
*/

/*func UpdateTeams(t *testing.T) {
	reader = strings.NewReader(`{"Teams" : {1,2,3,4}}`)
	request, err := http.NewRequest("PUT", baseURL+"/api/coachs/updatecoachteam/" + rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}
*/
