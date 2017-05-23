//
// Fichier     : actions_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer toutes les interractions nécessaires à la création,
// la modification, la suppression et la récupération des informations
// d'une action.
//

package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/TSAP-Laval/acquisition-backend/api"
)

// TestStartServerStarted test le démarrage du serveur lorsqu'il est déjà démarré
func TestStartServerStarted(t *testing.T) {
	service.Start()
}

// TestBD test la création de la base de donnée avec erreur de connexion à la base de données
// avec erreur de connexion à la base de données
func TestBDErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/bd", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestBD test la création de la base de donnée
func TestBD(t *testing.T) {
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
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/seed", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestSeed test le remplissage de la base de donnée avec des informations bidons
func TestSeed(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/seed", reader)

	if err != nil {
		t.Error(err)
		return
	}

	PostRequestHandler(request, t)
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
	reader = strings.NewReader(
		`{
			"Email": "admin@admin.ca",
			"PassHash": "aaaaa"
		}`,
	)
	request, err := http.NewRequest("POST", baseURL+"/api/auth", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetTokenExpire test la récupération d'un token avec les informations d'authentification d'un administrateur
// mais ayant un token expiré
func TestGetTokenExpire(t *testing.T) {
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
	// Je n'utiliserai pas la fonction GetRequestHandler, car la requête doit être non sécurisée
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

	if !strings.Contains(string(bodyBuffer), "signature is invalid") && !strings.Contains(string(bodyBuffer), "Token is expired") {
		t.Error("Error expected :", string(bodyBuffer))
	}
}

// TestCreerPartiePourActions test la création d'une partie.
func TestCreerPartiePourActions(t *testing.T) {
	reader = strings.NewReader(
		`{
			"Date": "2016-06-25 06:02",
            "FieldCondition": "Correcte",
            "LocationID": 1,
            "OpposingTeam": "Team",
            "SeasonID": 1,
            "Status": "Local",
            "TeamID": 1
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	// On garde en mémoire l'ID de la partie venant d'être créée
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var ga api.Games
	err = json.Unmarshal(bodyBuffer, &ga)
	if err != nil {
		t.Error(err)
	}

	rmID = strconv.Itoa(int(ga.ID))

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
		var me MessageError
		responseMapping(&me, res)
		t.Errorf("Error: %s", me.Err)
	}
}

// TestCreerActionErrBD test la création d'une nouvelle action
// avec erreur de connexion à la base de données
func TestCreerActionErrBD(t *testing.T) {
	reader = strings.NewReader(
		`{
			"ActionTypeID" : 1,
			"ReceptionTypeID" : 1,
			"ZoneID" : 1,
			"GameID" : 1,
			"X1" : 1,
			"Y1" : 1,
			"X2" : 1,
			"Y2" : 1,
			"X3" : 1,
			"Y3" : 1,
			"Time" : 30,
			"HomeScore" : 0,
			"GuestScore" : 0,
			"PlayerID" : 0,
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestCreerActionErr test la création d'une nouvelle action
// avec erreur dans le JSON
func TestCreerActionErr(t *testing.T) {
	reader = strings.NewReader(
		`{
			"ActionTypeID 1,
			"ReceptionTypeID" : 1,
			"ZoneID" : 1,
			"GameID" : 1,
			"X1" : 1,
			"Y1" : 1,
			"X2" : 1,
			"Y2" : 1,
			"X3" : 1,
			"Y3" : 1,
			"Time" : 30,
			"HomeScore" : 0,
			"GuestScore" : 0,
			"PlayerID" : 0,
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions", reader)

	if err != nil {
		t.Error(err)
		return
	}

	BadRequestHandler(request, t)
}

/// TestCreerAction test la création d'une nouvelle action
func TestCreerAction(t *testing.T) {
	reader = strings.NewReader(
		`{
			"ActionTypeID" : 1,
			"ReceptionTypeID" : 1,
			"ZoneID" : 1,
			"GameID" : 1,
			"X1" : 1,
			"Y1" : 1,
			"X2" : 1,
			"Y2" : 1,
			"X3" : 1,
			"Y3" : 1,
			"Time" : 30,
			"HomeScore" : 0,
			"GuestScore" : 0,
			"PlayerID" : 0
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions", reader)

	if err != nil {
		t.Error(err)
		return
	}

	PostRequestHandler(request, t)
}

/// TestCreerActionExiste test la création d'une action existante
func TestCreerActionExiste(t *testing.T) {
	reader = strings.NewReader(
		`{
			"ActionTypeID" : 1,
			"ReceptionTypeID" : 1,
			"ZoneID" : 1,
			"GameID" : 1,
			"X1" : 1,
			"Y1" : 1,
			"X2" : 1,
			"Y2" : 1,
			"X3" : 1,
			"Y3" : 1,
			"Time" : 30,
			"HomeScore" : 0,
			"GuestScore" : 0,
			"PlayerID" : 0
		}`)

	request, err := http.NewRequest("POST", baseURL+"/api/actions", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Une action existe déjà à ce moment précis de la partie !") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestGetActionsErrBD test la récupération des actions
// avec erreur de connexion à la base de données
func TestGetActionsErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+rmID+"/actions", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetActions test la récupération des actions d'une partie
func TestGetActions(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+rmID+"/actions", reader)

	if err != nil {
		t.Error(err)
		return
	}

	GetRequestHandler(request, t)
}

// TestGetActionsErrID test la récupération des actions d'une partie avec un
// mauvais identifiant de partie
func TestGetActionsErrID(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/0/actions", reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Aucune partie ne correspond.") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestSupprimerActionsErrBD test la suppression d'une action
// avec erreur de connexion à la base de données
func TestSupprimerActionsErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/actions/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestSupprimerActions test la suppression d'une action
func TestSupprimerActions(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/actions/"+rmID, reader)

	if err != nil {
		t.Error(err)
	}

	DeleteHandler(request, t)
}

// TestSupprimerActionsErr test la suppression d'une action déjà supprimée
func TestSupprimerActionsErr(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/actions/"+rmID, reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	var me MessageError
	responseMapping(&me, res)

	if res.StatusCode != 404 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "Aucune action ne correspond. Elle doit déjà avoir été supprimée!") {
		t.Error("Error expected : ", me.Err)
	}
}

// TestSupprimerPartiePourAction test la suppression de la partie préalablement créée
func TestSupprimerPartiePourAction(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/parties/"+rmID, reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}
