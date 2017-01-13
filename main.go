package api

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetRouter retourne les routes de l'API
func GetRouter() *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "aaaaa")
	})

	return router
}
