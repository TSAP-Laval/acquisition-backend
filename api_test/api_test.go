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

	err := envconfig.Process("TSAP", &a)

	if err != nil {
		panic(err)
	}

	service := api.New(os.Stdout, &a)
	service.Start()

	baseURL = "http://localhost:3000"
}
