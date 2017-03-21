package api_test

import (
	"io"
	"net/http/httptest"
	"os"

	"github.com/TSAP-Laval/acquisition-backend/api"
	"github.com/kelseyhightower/envconfig"
)

var (
	server  *httptest.Server
	reader  io.Reader
	baseURL string
	rmID    string
)

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

	baseURL = "http://localhost" + a.Port
}
