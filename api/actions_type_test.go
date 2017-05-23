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

// TestGetAllActionsTypesErrToken test la récupération de tous les type d'action avec un token expiré
func TestGetAllActionsTypesErrToken(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types", reader)
	request.Header.Add("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE0OTQ4MTE4MjMsImV4cCI6MTQ5NDgxMTgyNCwiYXVkIjoiIiwic3ViIjoiIiwiYWRtaW4iOiJ0cnVlIn0.n-eE_AJErTBpbuR78Cb4wCEEeleBEmG_6j0N_DDlVYs")
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 401 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(string(bodyBuffer), "signature is invalid") { // TODO : Token is expired
		t.Error("Error expected :", string(bodyBuffer))
	}
}

// TestGetAllActionsTypesErrBD test la récupération de tous les type d'action
func TestGetAllActionsTypesErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetAllActionsTypes test la récupération de tous les type d'action
func TestGetAllActionsTypes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var at []api.ActionsType
	err = json.Unmarshal(bodyBuffer, &at)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if len(at) < 1 {
		LogErrors(Messages{t, "Number of action types expected: %d", len(at), true, request, res})
	}
}

// TestCreerActionsType test la création d'un type d'action
func TestCreerActionsTypeErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Passe", 
			"Description": "Interceptée", 
			"TypeAction": "Reception"
		}`)
	request, err := http.NewRequest("POST", baseURL+"/api/actions/types", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestCreerActionsType test la création d'un type d'action
func TestCreerActionsType(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Passe", 
			"Description": "Interceptée", 
			"TypeAction": "Reception"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions/types", reader)
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

	// Ici, rmID est utilisé pour permettre la récupération de
	// l'élément créé
	rmID = fmt.Sprintf("%d", at.ID)

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestCreerActionsTypeExisteDeja test la création d'un type d'action qui existe déjà
func TestCreerActionsTypeExisteDeja(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Passe", 
			"Description": "Interceptée", 
			"TypeAction": "Reception"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions/types", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Un type d'action avec le même nom existe déjà") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerActionsTypeExisteDeja test la création d'un type d'action qui existe déjà
func TestCreerActionsTypeErr(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name "Passe", 
			"Description": "Interceptée", 
			"TypeAction": "Reception"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions/types", reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestGetActionsTypes test la récupération du type d'action venant d'être créé
func TestGetActionsTypes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types/"+rmID, reader)

	if err != nil {
		t.Error(err)
		return
	}

	GetRequestHandler(request, t)
}

// TestGetActionsTypesExistePas test la récupération du type d'action n'existant point
func TestGetActionsTypesExistePas(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types/100", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "Aucun type d'action ne correspond à celui entré") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestGetAllReceptionTypeErrBD test la création d'un type d'action
// avec erreur de connexion à la base de données
func TestGetAllReceptionTypeErrBD(t *testing.T) {
	reader = strings.NewReader(``)
	request, err := http.NewRequest("GET", baseURL+"/api/receptions", reader)

	if err != nil {
		t.Error(err)
		return
	}

	BDErrorHandler(request, t)
}

// TestGetAllReceptionType test la récupération du tous les types de réceptions
func TestGetAllReceptionType(t *testing.T) {
	reader = strings.NewReader(``)
	request, err := http.NewRequest("GET", baseURL+"/api/receptions", reader)
	res, err := SecureRequest(request)

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var rt []api.ReceptionType
	err = json.Unmarshal(bodyBuffer, &rt)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if len(rt) < 1 {
		LogErrors(Messages{t, "Number of reception types expected: %d", len(rt), true, request, res})
	}
}
