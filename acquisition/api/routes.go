package api

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetRouter retourne les routes de l'API
func GetRouter() *gin.Engine {
	// Crée un routeur gin avec un middleware par défaut:
	// logger et recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "okok")
		println("okok")
	})

	router.POST("/upload", func(c *gin.Context) {
		println("okok")
	})

	return router
}
