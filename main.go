package main

import (
	"log"
	"net/http"

	"github.com/TSAP-Laval/acquisition-backend/acquisition/api"
)

func main() {
	router := api.GetRouter()

	log.Fatal(http.ListenAndServe(":3000", router))
}
