package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	// JSON Web Token
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// JWTMiddleware handler the JWT middleware
func (a *AcquisitionService) JWTMiddleware(h http.Handler) http.Handler {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(a.keys.JWT), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	}).Handler(h)
}

// SecureHeaders adds secure headers to the API
/*func (a *AcquisitionService) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
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
		}

		// Add X-XSS-Protection header
		w.Header().Add("X-XSS-Protection", "1")

		// Add X-Content-Type-Options header
		w.Header().Add("X-Content-Type-Options", "nosniff")

		// Prevent page from being displayed in an iframe
		w.Header().Add("X-Frame-Options", "DENY")

		next.ServeHTTP(w, r)
	})
}*/

// Login handle jwt creation
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
		fmt.Println(string(body))
		var ad Admins
		err = json.Unmarshal(body, &ad)
		if err != nil {
			a.ErrorHandler(w, err)
			return
		}

		pass := []byte(ad.PassHash)

		fmt.Println("Email : ", ad.Email)
		fmt.Println("Pass : ", pass)
		db.Where("email = ?", ad.Email).First(&ad)

		if ad.Email != "" {
			if err = bcrypt.CompareHashAndPassword([]byte(ad.PassHash), pass); err != nil {
				a.ErrorHandler(w, err)
				return
			}

			// Create the token
			token := jwt.New(jwt.SigningMethodHS256)

			// Create a map to store our claims
			claims := token.Claims.(jwt.MapClaims)

			// Set token claims
			claims["admin"] = true
			claims["name"] = "admin"
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

			// Sign the token with our secret
			tokenString, _ := token.SignedString([]byte(a.keys.JWT))
			msg := map[string]string{"token": tokenString}

			Message(w, msg, http.StatusOK)
		} else {
			msg := map[string]string{"error": "Le mot de passe ou l'adresse email entrées est invalide."}
			Message(w, msg, http.StatusBadRequest)
		}
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}

// Authenticate provides Authentication middleware for handlers
func (a *AcquisitionService) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Get token from the Authorization header
		// format: Authorization: Bearer
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// If the token is empty...
		if token == "" {
			// If we get here, the required token is missing
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Now parse the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				return nil, msg
			}
			return []byte(a.keys.JWT), nil
		})
		if err != nil {
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}

		// Check token is valid
		if parsedToken != nil && parsedToken.Valid {
			// Everything worked! Set the user in the context.
			context.Set(r, "user", parsedToken)
			next.ServeHTTP(w, r)
		}

		// Token is invalid
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	})
}
