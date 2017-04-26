//
// TEST
//
// Fichier     : actions.go
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

// TODO: Linter le code...
// TODO: Mettre des commentaire au dessus des fonctions

func PostActionType(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Passe", 
			"Description": "Interceptée", 
			"ControlType": "Neutre", 
			"MovementType": "Postif"
		}`)

	request, err := http.NewRequest("POST", baseURL+"//addactiontype", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var c api.Coaches
	err = json.Unmarshal(bodyBuffer, &c)
	if err != nil {
		t.Error(err)
	}

	rmID = fmt.Sprintf("%d", c.ID)

	if res.StatusCode != 201 {
		t.Errorf("Action posted: %d", res.StatusCode)
	}
}

func GetActionsTypes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"//actiontype", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode)
	}
}

func GetMovementType(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"//movementType", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Error
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode)
	}
}
