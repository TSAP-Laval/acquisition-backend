package api

import (
	"encoding/json"
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
		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
		defer db.Close()
		fmt.Println(err)

		db.AutoMigrate(&TypeAction{})
		db.AutoMigrate(&Sport{})
		db.AutoMigrate(&Niveau{})
		db.AutoMigrate(&Entraineur{})
		db.AutoMigrate(&Joueur{})
		db.AutoMigrate(&Equipe{})
		db.AutoMigrate(&Zone{})
		db.AutoMigrate(&Saison{})
		db.AutoMigrate(&Lieu{})
		db.AutoMigrate(&Video{})
		db.AutoMigrate(&Partie{})
		db.AutoMigrate(&Action{})
	})
	router.GET("/getJoueurs", func(c *gin.Context) {
		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
		defer db.Close()
		fmt.Println(err)

		user := []Joueur{}
		db.Find(&user)

		userJSON, _ := json.Marshal(user)
		fmt.Println(string(userJSON))

		c.Header().Set("Content-Type", "application/json")
		c.Write(userJSON)

	})
	router.GET("/GetActions", func(c *gin.Context) {
		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=tsapBack sslmode=disable password=alex1997")
		defer db.Close()
		fmt.Println(err)
		user := []TypeAction{}
		db.Find(&user)
		c.JSON(200, &user)

	})
	return router
}
func foo(w http.ResponseWriter, r *http.Request) {

}
