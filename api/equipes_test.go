//
// TEST
//
// Fichier     : equipe_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la seppression et la récupération des informations
// d'une équipes.
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

// TestBD test la création de la base de donnée
func TestBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/bd", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// TestSeed test le remplissage de la base de donnée avec des informations bidons
func TestSeed(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/seed", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// TestGetEquipes test la récupération de toutes les équipes
func TestGetEquipes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te []api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	// On s'assure qu'il y ait au moins une équipe à la base
	if len(te) <= 0 {
		t.Errorf("Number of teams expected: %d", len(te))
	}
}

// TestCreerEquipe test la création d'une équipe.
// Cette équipe sera utilisée pour le reste des opérations
// (modification, suppression)
func TestCreerEquipe(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Lequipe", 
			"City": "Quebec", 
			"CategoryID": 1, 
			"SportID": 1
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	// On garde en mémoire l'ID de la partie venant d'être créée
	// pour pouvoir la modifier et supprimer plus tard...
	rmID = fmt.Sprintf("%d", te.ID)

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if te.Name != "Lequipe" {
		t.Error("Name expected: ", te.Name)
	}

	if te.City != "Quebec" {
		t.Error("City expected: ", te.City)
	}

	if te.CategoryID != 1 {
		t.Errorf("CategoryID expected: %d", te.CategoryID)
	}

	if te.SportID != 1 {
		t.Errorf("SportID expected: %d", te.SportID)
	}
}

// TestCreerEquipeErrEmpty test que créer une équipe avec des informations
// manquante retourne une erreur
func TestCreerEquipeErrEmpty(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "UNE equipe", 
			"": "", 
			"CategoryID": 1, 
			"SportID": 1
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerEquipeMauvaiseInfo test que créer une équipe avec de
// mauvaises informations. Doit retourner une erreur.
func TestCreerEquipeMauvaiseInfo(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name: "UNE equipe", 
		`)

	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// TestCreerEquipeVide test que créer une équipe sans
// information. Doit retourner une erreur.
func TestCreerEquipeVide(t *testing.T) {
	reader = strings.NewReader(``)

	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerEquipeErrExiste test que de créer une équipe qui
// existe déjà retourne une erreur
func TestCreerEquipeErrExiste(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "Lequipe", 
			"City": "Quebec", 
			"CategoryID": 1, 
			"SportID": 1
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/equipes", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 401 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "Une équipe de même nom existe déjà") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestGetEquipe test la récupération de l'équipe créée
func TestGetEquipe(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te []api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if te[0].Name != "Lequipe" {
		t.Error("Name expected: ", te[0].Name)
	}

	if te[0].City != "Quebec" {
		t.Error("City expected: ", te[0].City)
	}

	if te[0].CategoryID != 1 {
		t.Errorf("CategoryID expected: %d", te[0].CategoryID)
	}

	if te[0].SportID != 1 {
		t.Errorf("SportID expected: %d", te[0].SportID)
	}
}

// TestModifierEquipe test la modification de l'équipe créée plus haut
func TestModifierEquipe(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LE equipe", 
			"City": "Montreal", 
			"CategoryID": 1, 
			"SportID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if te.Name != "LE equipe" {
		t.Error("Name expected: ", te.Name)
	}

	if te.City != "Montreal" {
		t.Error("City expected: ", te.City)
	}
}

// TestModifierEquipeName test que seul le nom est modifié
func TestModifierEquipeName(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "LES equipe",
			"City": "",
			"CategoryID": 1,
			"SportID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if te.Name != "LES equipe" {
		t.Error("Name expected: ", te.Name)
	}

	if te.City != "Montreal" {
		t.Error("City expected: ", te.City)
	}
}

// TestModifierEquipeCity test que seul la ville est modifiée
func TestModifierEquipeCity(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Name": "",
			"City": "Toronto",
			"CategoryID": 1,
			"SportID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	if te.Name != "LES equipe" {
		t.Error("Name expected: ", te.Name)
	}

	if te.City != "Toronto" {
		t.Error("City expected: ", te.City)
	}
}

// TestModifierEquipeSansModif test qu'il n'y ait pas de modification
func TestModifierEquipeVide(t *testing.T) {
	reader = strings.NewReader(``)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)
	if !strings.Contains(me.Err, "Veuillez choisir au moins un champs à modifier.") {
		t.Error("Error expected: ", me.Err)
	}
}

// TestModifierEquipeMauvaiseInfo test que la modification avec une
// erreur dans le JSON retourne une erreur
func TestModifierEquipeMauvaiseInfo(t *testing.T) {
	reader = strings.NewReader(`{
		"ERREUR : ""
	}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// TestSupprimerEquipe test la suppression de l'équipe préalablement créée
func TestSupprimerEquipe(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// TestGetEquipeErr test la récupération de l'équipe supprimée
func TestGetEquipeErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/equipes/LE", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var te []api.Teams
	err = json.Unmarshal(bodyBuffer, &te)
	if err != nil {
		t.Error(err)
	}

	if len(te) > 0 {
		t.Errorf("Number of team expected: %d", len(te))
	}
}

// TestSupprimerEquipeSupprime test le retour si l'on tente de supprimer
// une équipe qui l'est déjà été
func TestSupprimerEquipeSupprime(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/equipes/"+rmID, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)
	if !strings.Contains(me.Err, "Aucune equipe ne correspond. Elle doit déjà avoir été supprimée!") {
		t.Error("Error expected: ", me.Err)
	}
}
