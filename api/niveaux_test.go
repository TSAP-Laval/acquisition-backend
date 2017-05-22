//
// TEST
//
// Fichier     : niveaux_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de tester les interractions sur un niveau.
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

// TestGetNiveauxErrBD permet de récupérer les niveaux
// avec erreur de connexion à la base de données
func TestGetNiveauxErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/niveaux", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetNiveaux permet de récupérer tous les niveaux
func TestGetNiveaux(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/niveaux", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var ca []api.Categories
	err = json.Unmarshal(bodyBuffer, &ca)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	// Test simplement qu'il y a au moins un sport de créé
	if len(ca) < 1 {
		t.Error("Number of sports expected: ", len(ca))
	}
}
