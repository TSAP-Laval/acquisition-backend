package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
