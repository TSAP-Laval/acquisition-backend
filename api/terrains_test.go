//
// TEST
//
// Fichier     : terrains_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de tester les interractions sur un terrain.
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

// TestGetTerrains permet de récupérer les terrains
// avec erreur de connexion à la base de données
func TestGetTerrainsErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetTerrains permet de récupérer les terrains
func TestGetTerrains(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var lo []api.Locations
	err = json.Unmarshal(bodyBuffer, &lo)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	// Test simplement qu'il y a au moins un terrain de créé
	if len(lo) < 1 {
		t.Error("Number of location expected: ", len(lo))
	}
}

// TestCreerTerrainErrBD test la création d'un terrain
// avec erreur de connexion à la base de données
func TestCreerTerrainErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LE terrain", 
			"City": "Quebec", 
			"Address": "1231 une rue"
		}`)
	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestCreerTerrain test la création d'un terrain
func TestCreerTerrain(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LE terrain", 
			"City": "Quebec", 
			"Address": "1231 une rue"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var l api.Locations
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Error(err)
	}

	// On garde en mémoire l'ID du terrain venant d'être créé
	// pour pouvoir le modifier et supprimer plus tard...
	rmID = fmt.Sprintf("%d", l.ID)

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if l.Name != "LE terrain" {
		t.Error("Name expected: ", l.Name)
	}

	if l.City != "Quebec" {
		t.Error("City expected: ", l.City)
	}

	if l.Address != "1231 une rue" {
		t.Error("Address expected: ", l.Address)
	}
}

// TestCreerTerrainErrEmpty test que la création d'un terrain avec des informations
// manquante retourne une erreur.
func TestCreerTerrainErrEmpty(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "", 
			"City": "", 
			"Addrese": ""
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerTerrainErrMauvaiseInfos test que la création d'un terrain
// avec un JSON bidon retourne une erreur
func TestCreerTerrainErrMauvaiseInfos(t *testing.T) {
	reader = strings.NewReader(
		`{
			"RIEN": ""
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerTerrainErrExiste test la création d'un terrain existante.
// Doit retourner une erreur
func TestCreerTerrainErrExiste(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LE terrain", 
			"City": "Quebec", 
			"Address": "1234 autre rue"
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Un terrain de même nom existe déjà.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerTerrainVide test la création d'un terrain sans envoyer d'information.
// Doit retourner une erreur
func TestCreerTerrainVide(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerTerrainMauvaiseInfo test la création d'un terrain avec de mauvaises
// informations. Doit retourner une erreur.
func TestCreerTerrainMauvaiseInfo(t *testing.T) {
	// Envoie de l'information au format JSON avec une erreur
	reader = strings.NewReader(
		`{
			"RIEN: ""
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/terrains", reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestGetTerrain test la récupération du terrain préalablement créé
// avec erreur de connexion à la base de données
func TestGetTerrainErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains/LE", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetTerrain test la récupération du terrain préalablement créé
func TestGetTerrain(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains/LE", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var l []api.Locations
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Error(err)
	}

	// On garde en mémoire l'ID du terrain venant d'être créé
	// pour pouvoir le modifier et supprimer plus tard...
	rmID = fmt.Sprintf("%d", l[0].ID)

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if l[0].Name != "LE terrain" {
		t.Error("Name expected: ", l[0].Name)
	}

	if l[0].City != "Quebec" {
		t.Error("City expected: ", l[0].City)
	}

	if l[0].Address != "1231 une rue" {
		t.Error("Address expected: ", l[0].Address)
	}
}

// TestModifierTerrainErrBD test la modification du terrain préalablement créé
// avec erreur de connexion à la base de données
func TestModifierTerrainErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LE terrain", 
			"City": "Montreal", 
			"Address": ""
		}`)
	// rmID est utilisé ici pour permettre la modification de le terrain créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/terrains/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestModifierTerrain test la modification du terrain préalablement créé
func TestModifierTerrain(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LE terrain", 
			"City": "Montreal", 
			"Address": ""
		}`)

	// rmID est utilisé ici pour permettre la modification de le terrain créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/terrains/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var l api.Locations
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if l.Name != "LE terrain" {
		t.Error("Name expected: ", l.Name)
	}

	if l.City != "Montreal" {
		t.Error("City expected: ", l.City)
	}

	if l.Address != "1231 une rue" {
		t.Error("Address expected: ", l.Address)
	}
}

// TestModifierTerrainAutresInfos test la modification du terrain préalablement créé
// avec toutes les données. Les données devraient rester les mêmes
func TestModifierTerrainAutresInfos(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "", 
			"City": "", 
			"Address": ""
		}`)

	// rmID est utilisé ici pour permettre la modification de le terrain créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/terrains/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var l api.Locations
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if l.Name != "LE terrain" {
		t.Error("Name expected: ", l.Name)
	}

	if l.City != "Montreal" {
		t.Error("City expected: ", l.City)
	}

	if l.Address != "1231 une rue" {
		t.Error("Address expected: ", l.Address)
	}
}

// TestModifierTerrainVide test que la modification d'un terrain avec aucune information
// retourne une erreur
func TestModifierTerrainVide(t *testing.T) {
	reader = strings.NewReader(``)

	// rmID est utilisé ici pour permettre la modification de le terrain créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/terrains/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifierTerrainMauvaiseInfo test que la modification d'un terrain avec de mauvaises
// informations retourne une erreur
func TestModifierTerrainMauvaiseInfo(t *testing.T) {
	reader = strings.NewReader(
		`{
			"RIEN: ""
		}`)

	// rmID est utilisé ici pour permettre la modification de le terrain créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/terrains/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestSupprimerTerrainErrBD test la suppression du terrain préalablement créé
// avec erreur de connexion à la base de données
func TestSupprimerTerrainErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/terrains/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestSupprimerTerrain test la suppression du terrain préalablement créé
func TestSupprimerTerrain(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/terrains/"+rmID, reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}

// TestGetTerrainErr test que le terrain a bel et bien été supprimé
func TestGetTerrainErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/terrains/LE", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var lo []api.Locations
	err = json.Unmarshal(bodyBuffer, &lo)
	if err != nil {
		t.Error(err)
	}

	if len(lo) != 0 {
		t.Errorf("Number of location expected: %d", len(lo))
	}
}

// TestSupprimerTerrainSupprime test la suppression d'un terrain déjà supprimé.
func TestSupprimerTerrainSupprime(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/terrains/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)
	if !strings.Contains(me.Err, "Aucun terrain ne correspond.") {
		t.Error("Error expected: ", me.Err)
	}
}
