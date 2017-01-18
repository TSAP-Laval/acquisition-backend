package api

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/gin-gonic/gin.v1"
)

// GetRouter retourne les routes de l'API
func GetRouter() *gin.Engine {
	// Crée un routeur gin avec un middleware par défaut:
	// logger et recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "bbb")
	})
	router.GET("/actions", func(c *gin.Context) {
		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapTest sslmode=disable password=alex1997")
		defer db.Close()
		fmt.Println(err)
		types := []Actions_type{}
		db.Find(&types)
		fmt.Println(types)

	})
	router.GET("/actionsTest", func(c *gin.Context) {
		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapTest sslmode=disable password=alex1997")
		defer db.Close()
		fmt.Println(err)
		user := Actions_type{ID: "test", Name: "passe", Description: "une passe", Id_movement_type: 1}
		if db.NewRecord(user) {
			fmt.Println("Test")
			db.Create(&user)
			db.NewRecord(user) // => return `false` after `user` created
		} else {
			fmt.Println("Test22")
		}

	})

	return router
}
