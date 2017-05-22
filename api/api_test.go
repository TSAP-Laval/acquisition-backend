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
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

// Permet de simuler le demarrage du serveur le temps des tests
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

func BadRequestHandler(request *http.Request, t *testing.T) (me MessageError) {
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return MessageError{}
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	return me
}

func DeleteHandler(request *http.Request, t *testing.T) {
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
		var m MessageError
		responseMapping(&m, res)
		t.Errorf("Error: %s", m.Err)
	}
}

// BDErrorHandler gère l'envoie d'une requête avec une erreur de connexion à la base de données
func BDErrorHandler(request *http.Request, t *testing.T) {
	acqConf.ConnectionString = "host=localhost user=aaaaa dbname=tsap_acquisition sslmode=disable password="
	defer goodConnectionString()

	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

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

func goodConnectionString() {
	acqConf.ConnectionString = "host=localhost user=postgres dbname=tsap_acquisition sslmode=disable password="
}
