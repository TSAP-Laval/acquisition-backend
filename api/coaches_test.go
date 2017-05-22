//
// Fichier     : coachs_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la suppression et la récupération des informations
// d'un entraineur.
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

// TestCreerCoachErrBD test la création d'un nouvel entraineur
// avec erreur de connexion à la base de données
func TestCreerCoachErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader(
		`{
			"Fname": "entraineur", 
			"Lname": "tony", 
			"Email": "entraineur@entraineur.ca",
			"PassHash": "$2a$10$txBDGNabCC0j.n8wFURChO9KazKeQFOyPtUliyH.V5b7DbTkwsJxe"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coaches", reader)
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

// TestCreerCoach test la création d'un nouvel entraineur
// avec des erreurs dans le JSON
func TestCreerCoachErr(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader(
		`{
			"Fname " entraineur", 
			"Lname": "tony", 
			"Email": "entraineur@entraineur.ca",
			"PassHash": "$2a$10$txBDGNabCC0j.n8wFURChO9KazKeQFOyPtUliyH.V5b7DbTkwsJxe"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coaches", reader)
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

	if !strings.Contains(me.Err, "Certaines informations entrées sont invalides!") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestCreerCoach test la création d'un nouvel entraineur
func TestCreerCoach(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Fname": "entraineur", 
			"Lname": "tony", 
			"Email": "entraineur@entraineur.ca",
			"PassHash": "$2a$10$txBDGNabCC0j.n8wFURChO9KazKeQFOyPtUliyH.V5b7DbTkwsJxe"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coaches", reader)
	res, err := SecureRequest(request)

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
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TestCreerCoachInfosManquantes test la création d'un nouvel entraineur
// avec certaines informations manquantes
func TestCreerCoachInfosManquantes(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Fname": "entraineur", 
			"Lname": "tony", 
			"Email": "",
			"PassHash": "$2a$10$txBDGNabCC0j.n8wFURChO9KazKeQFOyPtUliyH.V5b7DbTkwsJxe"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coaches", reader)
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

	if !strings.Contains(me.Err, "Certaines informations sont manquantes!") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestCreerCoachCree test la création d'un entraineur déjà créé
func TestCreerCoachCree(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Fname": "entraineur", 
			"Lname": "tony", 
			"Email": "entraineur@entraineur.ca",
			"PassHash": "$2a$10$txBDGNabCC0j.n8wFURChO9KazKeQFOyPtUliyH.V5b7DbTkwsJxe"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/coaches", reader)
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

	if !strings.Contains(me.Err, "Un entraineur avec la même adresse courriel existe déjà") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestGetCoachsErrBD test la récupération des entraineurs
// avec erreur de connexion à la base de données
func TestGetCoachsErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/coaches", reader)
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

// TestGetCoachs test la récupération des entraineurs
func TestGetCoachs(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/coaches", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode)
	}
}

// TestUpdateCoachErrBD test l'ajout d'un entraineur dans une équipe
// avec erreur de connexion dans la base de données
func TestUpdateCoachErrBD(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("PUT", baseURL+"/api/coaches/"+rmID+"/equipes/1", reader)
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

// TestUpdateCoach test l'ajout d'un entraineur dans une équipe
func TestUpdateCoach(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("PUT", baseURL+"/api/coaches/"+rmID+"/equipes/1", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		bodyBuffer, _ := ioutil.ReadAll(res.Body)
		var me MessageError
		err = json.Unmarshal(bodyBuffer, &me)
		if err != nil {
			t.Error(err)
			return
		}
		t.Error("Error expected : ", me.Err)
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

// TestUpdateCoachDeja test l'ajout d'un entraineur dans une équipe dont il fait déjà partie
func TestUpdateCoachDeja(t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
	reader = strings.NewReader("")
	request, err := http.NewRequest("PUT", baseURL+"/api/coaches/"+rmID+"/equipes/1", reader)
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

	if !strings.Contains(me.Err, "L'entraineur fait déjà parti de l'équipe") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestUpdateCoachExistePas test l'ajout d'un entraineur dans une équipe avec un entraineur inexistant
func TestUpdateCoachExistePas(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("PUT", baseURL+"/api/coaches/100/equipes/1", reader)
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

	if !strings.Contains(me.Err, "Aucun entraineur ne correspond") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestUpdateCoachTeamExistePas test l'ajout d'un entraineur dans une équipe avec une équipe inexistant
func TestUpdateCoachTeamExistePas(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("PUT", baseURL+"/api/coaches/1/equipes/100", reader)
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

	if !strings.Contains(me.Err, "Aucune équipe ne correspond") {
		t.Error("Error expected : ", me.Err)
	}
}
