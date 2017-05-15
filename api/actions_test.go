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

// TestBD test la création de la base de donnée avec erreur de connexion à la base de données
// avec erreur de connexion à la base de données
func TestBDErr(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/bd", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestBD test la création de la base de donnée
func TestBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/bd", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestSeedErr test le remplissage de la base de donnée avec des informations bidons
// et erreur de connexion à la base de données
func TestSeedErr(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/seed", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestSeed test le remplissage de la base de donnée avec des informations bidons
func TestSeed(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/seed", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestGetTokenErr test la récupération d'un token avec des informations au mauvais format
// et avec erreur de connexion à la base de données
func TestGetTokenErr(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Email "admin@admin.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBuffer, &token)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}

	if token.Token != "" {
		t.Error("Token expected : ", token.Token)
	}
}

// TestGetTokenInvalid test la récupération d'un token avec des informations invalides
func TestGetTokenInvalid(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.ca",
			"PassHash": "1234534523"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBuffer, &token)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}

	if token.Token != "" {
		t.Error("Token :", token.Token)
		t.Error("Token is not empty !")
	}
}

// TestGetTokenInvalidEmail test la récupération d'un token avec des informations invalides
func TestGetTokenInvalidEmail(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.com",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBuffer, &token)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}

	if token.Token != "" {
		t.Error("Token :", token.Token)
		t.Error("Token is not empty !")
	}
}

// TestGetTokenErrBD test la récupération d'un token avec les informations d'authentification d'un administrateur
// avec erreur de connexion à la base de données
func TestGetTokenErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestGetTokenExpire test la récupération d'un token avec les informations d'authentification d'un administrateur
// mais ayant un token expiré
func TestGetTokenExpire(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader(
		`{
			"Email": "mauvais@mauvais.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBuffer, &token)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		var me MessageError
		err = json.Unmarshal(bodyBuffer, &me)
		if err != nil {
			t.Error(err)
			return
		}
		t.Error("Error : ", me.Err)
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}

	if token.Token == "" {
		t.Error("Token is empty")
	}
}

// TestGetToken test la récupération d'un token avec les informations d'authentification d'un administrateur
func TestGetToken(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBuffer, &token)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}

	if token.Token == "" {
		t.Error("Token is empty")
	}
}

// TestGetTokenSecondTime test la récupération d'un token avec les informations d'authentification d'un administrateur
func TestGetTokenSecondTime(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(bodyBuffer, &token)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}

	if token.Token == "" {
		t.Error("Token is empty")
	}
}

// TestGetTokenOptions test l'envoie d'un requête d'options pour l'authetification
func TestGetTokenOptions(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("OPTIONS", baseURL+"/api/auth", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, nil, res})
	}
}

// TestGetAllActionsTypesErrTokenSignature test la récupération de tous les type d'action avec un token mal signé
func TestGetAllActionsTypesErrTokenSignature(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types", reader)
	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNDk0NzgyNjg5fQ.5yuANzVeFq7HPMAmNQIk_QqWZxh2ZfgiWqvDoMtUTGs")
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 401 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(string(bodyBuffer), "signature is invalid") {
		t.Error("Error expected :", string(bodyBuffer))
	}
}

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

	if !strings.Contains(string(bodyBuffer), "Token is expired") {
		t.Error("Error expected :", string(bodyBuffer))
	}
}

// TestGetAllActionsTypesErrBD test la récupération de tous les type d'action
func TestGetAllActionsTypesErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestGetAllActionsTypes test la récupération de tous les type d'action
func TestGetAllActionsTypes(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
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
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
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

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "pq: role \"aaaaa\" does not exist") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestCreerActionsType test la création d'un type d'action
func TestCreerActionsType(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
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

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

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
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestGetActionsTypes test la récupération du type d'action venant d'être créé
func TestGetActionsTypes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/actions/types/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
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
