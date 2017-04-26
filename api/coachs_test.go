//
// Fichier     : coachs_test.go
// Développeur : ?
//
// Commentaire expliquant le code, les fonctions...
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

// TODO: Changer le nom du fichier et ses références pour coach au pluriel...
//		http://www.wordhippo.com/what-is/the-plural-of/coach.html

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

func GetCoachs(t *testing.T) {
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
