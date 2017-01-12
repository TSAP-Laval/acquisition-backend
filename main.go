package main

import (
	"log"
	"net/http"
	"tsap/back-end/acquisition-backend/acquisition/api"
)

func main() {
	router := api.GetRouter()

	log.Fatal(http.ListenAndServe(":3000", router))
}
