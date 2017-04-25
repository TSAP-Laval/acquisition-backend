//
// TEST
//
// Fichier     : api_test.go
// Développeur : Laurent Leclerc Poulin
//
// Fichier `main` pour lancer les tests
//

package api_test

import (
	"io"
	"net/http/httptest"
	"os"
	"time"

	"github.com/TSAP-Laval/acquisition-backend/api"
	"github.com/kelseyhightower/envconfig"
)

var (
	server  *httptest.Server
	reader  io.Reader
	baseURL string
	rmID    string
)

// Permet de simuler le démarrage du serveur le temps des tests
func init() {
	var a api.AcquisitionConfiguration
	var k api.Keys

	err := envconfig.Process("TSAP", &a)
	err = envconfig.Process("KEYS", &k)

	if err != nil {
		panic(err)
	}

	service := api.New(os.Stdout, &a, &k)
	service.Start()

	// ** IMPORTANT **
	// Permet de s'assurer que le serveur a bel et bien démarré
	// avant de lancer les tests.
	time.Sleep(5 * time.Second)

	baseURL = "http://localhost" + a.Port
}
