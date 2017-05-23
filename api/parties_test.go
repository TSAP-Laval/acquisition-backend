//
// TEST
//
// Fichier     : parties_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de tester les interractions sur une partie.
//

package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/TSAP-Laval/acquisition-backend/api"
)

// TestGetPartiesErrBD test la récupération de toutes les parties
// avec erreur de connexion à la base de données
func TestGetPartiesErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetParties test la récupération de toutes les parties
func TestGetParties(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var ga []api.Games
	err = json.Unmarshal(bodyBuffer, &ga)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	// On s'assure qu'il y ait au moins une partie à la base
	if len(ga) != 0 {
		t.Errorf("Number of teams expected: %d", len(ga))
	}
}

// TestCreerPartieErrBD test la création d'une partie.
// avec erreur de connexion à la base de données
func TestCreerPartieErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-21 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Rien",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestCreerPartie test la création d'une partie.
// Cette partie sera utilisée pour le reste des opérations
// (modification, suppression)
func TestCreerPartie(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-21 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Rien",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var ga api.Games
	err = json.Unmarshal(bodyBuffer, &ga)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
		var me MessageError
		responseMapping(&me, res)
		t.Errorf("Error: %s", me.Err)
	}

	if ga.Date != "2016-06-21 06:02" {
		t.Error("Date expected: ", ga.Date)
	}

	if ga.FieldCondition != "Correcte" {
		t.Error("FieldCondition expected: ", ga.FieldCondition)
	}

	if ga.LocationID != 1 {
		t.Errorf("LocationID expected: %d", ga.LocationID)
	}

	if ga.OpposingTeam != "Rien" {
		t.Error("OpposingTeam expected: ", ga.OpposingTeam)
	}

	if ga.SeasonID != 1 {
		t.Errorf("SeasonID expected: %d", ga.SeasonID)
	}

	if ga.Status != "Local" {
		t.Error("Status expected: ", ga.Status)
	}

	if ga.TeamID != 1 {
		t.Error("TeamID expected: ", ga.TeamID)
	}
}

// TestCreerPartieErrEmpty test que créer une partie avec des informations
// manquante retourne une erreur
func TestCreerPartieErrEmpty(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-22 06:02",
            "": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Rien",
            "": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs!") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestCreerPartieMauvaiseInfo test que créer une partie avec de
// mauvaises informations. Doit retourner une erreur.
func TestCreerPartieMauvaiseInfo(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Dt: "Uen Deat", 
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestCreerPartieVide test que créer une partie sans
// information. Doit retourner une erreur.
func TestCreerPartieVide(t *testing.T) {
	reader = strings.NewReader(``)

	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	// On garde en mémoire l'ID de la partie venant d'être créée
	// pour pouvoir la modifier et supprimer plus tard...
	var me MessageGameID
	responseMapping(&me, res)
	rmID = me.GameID

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestCreerPartieErr test que créer une partie avec
// une erreur dans le JSON. Doit retourner une erreur.
func TestCreerPartieErr(t *testing.T) {
	reader = strings.NewReader(`{
		"ERREUR" "
	}`)

	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestCreerPartieErrExiste test que de créer une partie qui
// existe déjà retourne une erreur
func TestCreerPartieErrExiste(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-21 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Rien",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Une partie de même date avec les mêmes equipes existe déjà!") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// TestGetPartieErrBD test la récupération de la partie créée
// avec erreur de connexion à la base de données
func TestGetPartieErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetPartie test la récupération de la partie créée
func TestGetPartie(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	GetRequestHandler(request, t)
}

// TestModifierPartieErrBD test la modification de la partie créée plus haut
// avec erreur de connexion à la base de données
func TestModifierPartieErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-22 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Ok",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)
	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestModifierPartie test la modification de la partie créée plus haut
// avec une erreur dans la clé du géodécodeur
func TestModifierPartieErrKeyGeodecoder(t *testing.T) {
	keys.Geodecoder = ""
	reader = strings.NewReader(
		`{
			"Date": "2016-06-22 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Ok",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifierPartieErrKeyWeather test la modification de la partie créée plus haut
// avec un erreur dans la clé de l'API de récupération de la température
func TestModifierPartieErrKeyWeather(t *testing.T) {
	keys.Geodecoder = "gCVEOBYTObcAaHsG5MXE3Uy0PF1kgkg0"
	keys.Weather = ""
	reader = strings.NewReader(
		`{
			"Date": "2016-06-22 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Ok",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifierPartieErrDate test la modification de la partie créée plus haut
// avec une erreur dans le format de la date
func TestModifierPartieErrDate(t *testing.T) {
	keys.Weather = "1e471424e05991d2f9ed9e39b9749ae0"
	reader = strings.NewReader(
		`{
			"Date": "2016-06-22 06:02 EDT",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Ok",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifierPartie test la modification de la partie créée plus haut
func TestModifierPartie(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-22 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Ok",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var ga api.Games
	err = json.Unmarshal(bodyBuffer, &ga)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
		var me MessageError
		json.Unmarshal(bodyBuffer, &me)
		if me.Err != "" {
			t.Error(me.Err)
		}
	}

	if ga.Date != "2016-06-22 06:02" {
		t.Error("Date expected: ", ga.Date)
	}

	if ga.FieldCondition != "Correcte" {
		t.Error("FieldCondition expected: ", ga.FieldCondition)
	}

	if ga.LocationID != 1 {
		t.Errorf("LocationID expected: %d", ga.LocationID)
	}

	if ga.OpposingTeam != "Ok" {
		t.Error("OpposingTeam expected: ", ga.OpposingTeam)
	}

	if ga.SeasonID != 1 {
		t.Errorf("SeasonID expected: %d", ga.SeasonID)
	}

	if ga.Status != "Local" {
		t.Error("Status expected: ", ga.Status)
	}

	if ga.TeamID != 1 {
		t.Error("TeamID expected: ", ga.TeamID)
	}

	if ga.Degree == "" {
		t.Error("Degree expected: ", ga.Degree)
	}

	if ga.Temperature == "" {
		t.Error("Temperature expected: ", ga.Temperature)
	}
}

// TestModifierPartieVide test qu'il n'y ait pas de modification
func TestModifierPartieVide(t *testing.T) {
	reader = strings.NewReader(``)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Veuillez remplir tous les champs.") {
		t.Error("Error expected: ", me.Err)
	}
}

// TestModifierPartieMauvaiseInfo test que la modification avec une
// erreur dans les les informations entrées et retourne une erreur
func TestModifierPartieMauvaiseInfo(t *testing.T) {
	reader = strings.NewReader(`{
		"ERREUR : ""
	}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifierPartieErr test que la modification avec une
// erreur dans le JSON et retourne une erreur
func TestModifierPartieErr(t *testing.T) {
	reader = strings.NewReader(`{
		"ERREUR ""
	}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BadRequestHandler(request, t)
}

// TestModifierPartieCreer test que la modification crée
// une nouvelle partie
func TestModifierPartieCreer(t *testing.T) {
	reader = strings.NewReader(`{
	}`)

	// rmID est utilisé ici pour permettre la modification de la partie créée plus haut
	request, err := http.NewRequest("PUT", baseURL+"/api/parties/0", reader)

	if err != nil {
		t.Error(err)
	}

	PostRequestHandler(request, t)
}

// TestSupprimerPartieErrBD test la suppression de la partie préalablement créée
// avec erreur de connexion à la base de données
func TestSupprimerPartieErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestSupprimerPartie test la suppression de la partie préalablement créée
func TestSupprimerPartie(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}

// TestGetPartieErr test la récupération de la partie supprimée
func TestGetPartieErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	var me MessageError
	responseMapping(&me, res)
	if !strings.Contains(me.Err, "Aucune partie ne correspond") {
		t.Error("Error expected: ", me.Err)
	}
}

// TestSupprimerPartieSupprime test le retour si l'on tente de supprimer
// une partie qui l'est déjà été
func TestSupprimerPartieSupprime(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Aucune partie ne correspond. Elle doit déjà avoir été supprimée!") {
		t.Error("Error expected: ", me.Err)
	}
}

// TestGetPartiesMulti test la récupération de toutes les parties.
// Doit y en avoir plusieurs
func TestGetPartiesMulti(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	var ga []api.Games
	err = json.Unmarshal(bodyBuffer, &ga)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	// On s'assure qu'il y ait au moins une partie à la base
	if len(ga) < 2 {
		t.Errorf("Number of teams expected: %d", len(ga))
	}
}

// Simule l'envoie d'une requête d'options pour une partie
func TestSendOptionsParies(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("OPTIONS", baseURL+"/api/upload/0", reader)

	if err != nil {
		t.Error(err)
	}

	GetRequestHandler(request, t)
}
