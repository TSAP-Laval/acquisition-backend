//
// Fichier     : auth.go
// Développeur : Laurent Leclerc Poulin
//
// Permet d'ajouter une couche de sécurité avec
// l'authentification et la gestion de JSON Web Token (JWT)
//

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	// JSON Web Token
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/jinzhu/gorm"
	// Algorithme utilisé pour l'encryption du mot de passe
	"golang.org/x/crypto/bcrypt"
)

// Claims utilisé pour stocker des données dans les JSON Web Token (JWT)
type Claims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

// JWTMiddleware gère le middleware pour les JSON Web Token (JWT)
func (a *AcquisitionService) JWTMiddleware(h http.Handler) http.Handler {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(a.keys.JWT), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	}).Handler(h)
}

// SecureHeaders adds secure headers to the API
func (a *AcquisitionService) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*var err error
		if len(a.allowedHosts) > 0 {
			isGoodHost := false
			for _, allowedHost := range a.allowedHosts {
				if strings.EqualFold(allowedHost, r.Host) {
					isGoodHost = true
					break
				}
			}
			if !isGoodHost {
				Message(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		// If there was an error, do not continue request
		if err != nil {
			Message(w, "Failed to check allowed hosts", http.StatusInternalServerError)
			return
		}*/

		// Add X-XSS-Protection header
		w.Header().Add("X-XSS-Protection", "1")

		// Add X-Content-Type-Options header
		w.Header().Add("X-Content-Type-Options", "nosniff")

		// Prevent page from being displayed in an iframe
		w.Header().Add("X-Frame-Options", "DENY")

		next.ServeHTTP(w, r)
	})
}

// Login gère la création des JSON Web Token (JWT)
func (a *AcquisitionService) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		db, err := gorm.Open(a.config.DatabaseDriver, a.config.ConnectionString)
		defer db.Close()

		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		var ad Admins
		err = json.Unmarshal(body, &ad)
		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		var dbAd Admins
		db.Where("email = ?", ad.Email).First(&dbAd)

		if dbAd.Email != "" {
			// Vérification que le mot de passe envoyé est bel et bien celui contenu dans la base de données et
			// faisant référence à celui correspondant au email entré
			if err = bcrypt.CompareHashAndPassword([]byte(dbAd.PassHash), []byte(ad.PassHash)); err == nil && ad.PassHash != "" {
				if dbAd.TokenLogin != "" {
					token, _ := jwt.ParseWithClaims(dbAd.TokenLogin, &Claims{}, func(token *jwt.Token) (interface{}, error) {
						return []byte(a.keys.JWT), nil
					})

					if _, ok := token.Claims.(*Claims); ok && token.Valid {
						tokenString, _ := token.SignedString([]byte(a.keys.JWT))
						msg := map[string]string{"token": tokenString}
						Message(w, msg, http.StatusOK)
						return
					}
					// Crée un nouveau token, le sauvegarde dans la base de données
					// et l'envoie au client.
					setToken(db, dbAd, w, a)
					return
				}
				// Crée un nouveau token, le sauvegarde dans la base de données
				// et l'envoie au client.
				setToken(db, dbAd, w, a)
				return
			}
		}
		// Dans le cas où le mot de passe où l'adresse courriel est invalide,
		// on envoie un message d'erreur au client.
		msg := map[string]string{"error": "Le mot de passe ou l'adresse email entrés est invalide."}
		Message(w, msg, http.StatusBadRequest)
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}

// setToken crée un token avec les informations du client, le signe avec les secret
// contenu dans une variable d'environnement, le sauvegarde dans la base de données
// et l'envoie au client
func setToken(db *gorm.DB, ad Admins, w http.ResponseWriter, a *AcquisitionService) {
	// Création des informations sauvegardées dans la token
	claims := Claims{}
	claims.Admin = true
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	// Création du token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signature du token avec le secret contenu dans une variable d'environnement
	tokenString, _ := token.SignedString([]byte(a.keys.JWT))

	var adm Admins
	adm.Email = ad.Email
	adm.PassHash = ad.PassHash
	adm.TokenLogin = tokenString
	db.Model(&Admins{}).Where("email = ?", ad.Email).Updates(adm)

	msg := map[string]string{"token": tokenString}
	Message(w, msg, http.StatusOK)
}
