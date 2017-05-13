//
// Fichier     : routes.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer le serveur et ses configurations
//

package api

import (
	"io"
	"log"
	"time"

	"encoding/json"

	"net/http"

	"github.com/braintree/manners"
	"github.com/didip/tollbooth"
	"github.com/gorilla/mux"
)

// AcquisitionConfiguration représente les paramètres
// exposés par l'application
type AcquisitionConfiguration struct {
	DatabaseDriver   string
	ConnectionString string
	Port             string
	Debug            bool
}

// Keys les clées utilisées pour les API de météo et géodécodage
type Keys struct {
	Geodecoder string
	Weather    string
	JWT        string
}

// AcquisitionService represents a single service instance
type AcquisitionService struct {
	logger *log.Logger
	config *AcquisitionConfiguration
	server *manners.GracefulServer
	keys   *Keys
}

// New crée une nouvelle instance du service
func New(writer io.Writer, config *AcquisitionConfiguration, keys *Keys) *AcquisitionService {

	return &AcquisitionService{
		logger: log.New(writer, "\x1b[32m[acquisition-api]\x1b[0m ", log.Flags()),
		config: config,
		server: manners.NewServer(),
		keys:   keys,
	}
}

// Info écrit un message vers le logger du service
func (a *AcquisitionService) Info(message string) {
	a.logger.SetPrefix("\x1b[32m[acquisition-api]\x1b[0m ")
	a.logger.Println(message)
}

// Info écrit un message d'erreur vers le logger du service
func (a *AcquisitionService) Error(message string) {
	a.logger.SetPrefix("\x1b[31m[acquisition-api]\x1b[0m ")
	a.logger.Printf("ERROR - %s\n", message)
}

// ErrorWrite écrit un message d'erreur en format JSON vers le writer
// passé en paramètre
func (a *AcquisitionService) ErrorWrite(message string, w io.Writer) error {
	bytes, err := json.Marshal(errorMessage{Error: message})

	if err != nil {
		return err
	}

	_, err = w.Write(bytes)

	return err
}

// Middleware applique les différents middleware
func (a *AcquisitionService) Middleware(h http.Handler) http.Handler {
	// Set CORS
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if a.config.Debug {
			// On ouvre l'accès de l'API si ce dernier est en debug
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		h.ServeHTTP(w, r)
	})
}

// RateLimiter est le middlewares limiteurs de débit utilisés par les endpoints
// 5 requêtes par seconde maximum
func (a *AcquisitionService) RateLimiter(h http.Handler) http.Handler {
	return tollbooth.LimitHandler(tollbooth.NewLimiter(20, time.Second), h)
}

// AddMiddleware adds middleware to a Handler
func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func (a *AcquisitionService) getRouter() http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Auth
	api.Handle("/auth",
		AddMiddleware(
			http.HandlerFunc(a.Login),
			a.RateLimiter,
		)).Methods("POST", "OPTIONS")
	// Actions
	api.HandleFunc("/action/movementType", a.GetMovementTypeHandler).Methods("GET")
	api.HandleFunc("/action/actiontype", a.GetAllActionsTypes).Methods("GET")
	api.HandleFunc("/action/addactiontype", a.PostActionType).Methods("POST")
	api.HandleFunc("/actionType/{id}", a.GetActionsTypeHandler).Methods("GET")
	api.HandleFunc("/action/addactiontype", a.PostActionType).Methods("POST")
	//Coachs
	api.HandleFunc("/coachs/coachs", a.GetCoachsHandler).Methods("GET")
	api.HandleFunc("/coachs/addcoach", a.PostCoachHandler).Methods("POST")
	api.HandleFunc("/coachs/addCoachTeam/{id}", a.AssignerEquipeCoach).Methods("PUT")

	// Upload
	api.Handle("/upload",
		AddMiddleware(
			a.SecureHeaders(http.HandlerFunc(a.UploadHandler)),
			a.JWTMiddleware,
			a.RateLimiter,
		))
	api.Handle("/upload/{game-id}",
		AddMiddleware(
			a.SecureHeaders(http.HandlerFunc(a.UploadHandler)),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("DELETE", "POST", "OPTIONS")

	// Videos
	api.Handle("/parties/{id}/videos/{part}",
		AddMiddleware(
			http.HandlerFunc(a.VideoHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET")

	// Terrains
	api.Handle("/terrains",
		AddMiddleware(
			http.HandlerFunc(a.GetTerrainsHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET")
	api.Handle("/terrains/{nom}",
		AddMiddleware(
			a.SecureHeaders(http.HandlerFunc(a.GetTerrainHandler)),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET")
	api.Handle("/terrains/{id}",
		AddMiddleware(
			http.HandlerFunc(a.TerrainsHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("DELETE", "PUT")
	api.Handle("/terrains",
		AddMiddleware(
			http.HandlerFunc(a.CreerTerrainHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("POST")

	// Equipes
	api.Handle("/equipes",
		AddMiddleware(
			http.HandlerFunc(a.GetEquipesHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET")
	api.Handle("/equipes/{nom}",
		AddMiddleware(
			http.HandlerFunc(a.GetEquipeHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET")
	api.Handle("/equipes/{id}",
		AddMiddleware(
			http.HandlerFunc(a.EquipesHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("DELETE", "PUT")
	api.Handle("/equipes",
		AddMiddleware(
			http.HandlerFunc(a.CreerEquipeHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("POST")

	// Parties
	api.Handle("/parties",
		AddMiddleware(
			http.HandlerFunc(a.PartiesHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET", "POST")
	api.Handle("/parties/{id}",
		AddMiddleware(
			http.HandlerFunc(a.PartieHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("GET", "PUT", "OPTIONS")
	api.Handle("/parties/{id}",
		AddMiddleware(
			http.HandlerFunc(a.SupprimerPartiesHandler),
			a.JWTMiddleware,
			a.RateLimiter,
		)).Methods("DELETE")
	// BD
	api.HandleFunc("/seed", a.RemplirBD).Methods("POST")
	api.HandleFunc("/bd", a.FaireBD).Methods("POST")
	// Actions
	api.HandleFunc("/actions", a.GetActions).Methods("GET")
	api.HandleFunc("/actions", a.PostAction).Methods("POST")
	// Joueurs
	api.HandleFunc("/joueurs", a.HandleJoueur).Methods("POST")
	api.HandleFunc("/joueurs/{id}", a.HandleJoueur).Methods("PUT", "OPTIONS", "DELETE")
	api.HandleFunc("/joueurs", a.GetJoueurs).Methods("GET")
	// Saisons
	api.HandleFunc("/saison", a.GetSeasons).Methods("GET")
	api.HandleFunc("/saison", a.PostSaison).Methods("POST")
	// Autres
	api.HandleFunc("/sports", a.GetSports).Methods("GET")
	api.HandleFunc("/niveau", a.GetNiveau).Methods("GET")
	return a.Middleware(api)
}

// Start démarre le service
func (a *AcquisitionService) Start() {
	go func() {
		a.server.Addr = a.config.Port
		a.server.Handler = a.getRouter()
		err := a.server.ListenAndServe()
		a.Error("Acquisition shutting down...")

		if err != nil {
			panic(err)
		}

	}()
	a.logger.Printf("TSAP-Acquisiton started on localhost%s... \n", a.config.Port)
}

// Stop arrête le service
func (a *AcquisitionService) Stop() {
	a.server.Close()
	a.Info("Acquisition shutting down...")
}
