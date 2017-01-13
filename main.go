package main

import "github.com/TSAP-Laval/acquisition-backend/acquisition/api"

// GetRouter retourne les routes de l'API
func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := api.GetRouter()

	router.Run(":3000")
}
