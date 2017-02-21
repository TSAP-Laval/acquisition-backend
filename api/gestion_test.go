package api

import (
	"fmt"
	"io"
	"net/http/httptest"

	"github.com/TSAP-Laval/acquisition-backend/api"
)

var (
	server   *httptest.Server
	reader   io.Reader //Ignore this for now
	usersUrl string
)

func init() {
	server = httptest.NewServer(api.getRouter()) //Creating new server with the user handlers

	usersUrl = fmt.Sprintf("%s/users", server.URL) //Grab the address for the API endpoint
}
