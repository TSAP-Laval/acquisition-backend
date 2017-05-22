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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/TSAP-Laval/acquisition-backend/api"
	"github.com/kelseyhightower/envconfig"
)

var (
	server  *httptest.Server
	reader  io.Reader
	baseURL string
	rmID    string
	service *api.AcquisitionService
	token   Token
	acqConf api.AcquisitionConfiguration
	keys    api.Keys
)

// Token est utilisé pour les communication sécurisées
type Token struct {
	Token string `json:"token"`
}

// Messages est utilisé pour logguer les messages d'erreurs
type Messages struct {
	Testing *testing.T
	Message string
	Object  interface{}
	Debug   bool
	Request *http.Request
	Reponse *http.Response
}

// Permet de simuler le démarrage du serveur le temps des tests
func init() {

	err := envconfig.Process("TSAP", &acqConf)
	err = envconfig.Process("KEYS", &keys)

	if err != nil {
		panic(err)
	}

	service = api.New(os.Stdout, &acqConf, &keys)
	service.Start()

	// ** IMPORTANT **
	// Permet de s'assurer que le serveur a bel et bien démarré
	// avant de lancer les tests. Sans ceci, les tests échouaient
	// une fois sur deux...
	time.Sleep(5 * time.Second)

	baseURL = "http://localhost" + acqConf.Port
}

// SecureRequest permet de faire des requêtes au serveur avec le token en en-tête
func SecureRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add("Authorization", "Bearer "+token.Token)
	return http.DefaultClient.Do(request)
}

// LogErrors permet de logguer les messages d'erreurs
func LogErrors(msg Messages) {
	if msg.Debug {
		if msg.Reponse != nil {
			bodyBuffer, _ := ioutil.ReadAll(msg.Reponse.Body)
			msg.Testing.Error("[DEBUG] Response : ", string(bodyBuffer))
		}

		if msg.Request != nil {
			msg.Testing.Error("[DEBUG] Header : ", msg.Request.Header)
		}
	}
	msg.Testing.Errorf(msg.Message, msg.Object)
}
