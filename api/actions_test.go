//
// TEST
//
// Fichier     : actions_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'un type d'action.
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

// TestCreerActionsType test la création d'un type d'action
func TestCreerActionsType(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Passe", 
			"Description": "Interceptée", 
			"TypeAction": "Reception"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/types", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var at api.ActionsType
	err = json.Unmarshal(bodyBuffer, &at)
	if err != nil {
		t.Error(err)
	}

	rmID = fmt.Sprintf("%d", at.ID)

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestGetActionsTypes test la récupération du type d'action venant d'être créé
func TestGetActionsTypes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/actions/types/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}
