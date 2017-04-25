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

	"encoding/json"

	"net/http"

	"github.com/braintree/manners"
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

func (a *AcquisitionService) getRouter() http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Actions
	api.HandleFunc("/action/movementType", a.GetMovementTypeHandler).Methods("GET")
	api.HandleFunc("/action/actiontype", a.GetAllActionsTypes).Methods("GET")
	api.HandleFunc("/action/addactiontype", a.PostActionType).Methods("POST")
	//Coachs
	r.HandleFunc("/api/coachs", a.GetCoachsHandler).Methods("GET")
	r.HandleFunc("/api/coachs/{coachID}", a.GetCoachsHandler).Methods("GET")
	r.HandleFunc("/api/coachs/addcoach", a.PostCoachHandler).Methods("POST")
	r.HandleFunc("/api/coachs/editcoach/{id}", a.AssignerEquipeCoach).Methods("PUT")
	// Upload
	api.HandleFunc("/upload", a.UploadHandler)
	api.HandleFunc("/upload/{game-id}", a.UploadHandler).Methods("DELETE", "POST", "OPTIONS")
	// Terrains
	api.HandleFunc("/terrains", a.GetTerrainsHandler).Methods("GET")
	api.HandleFunc("/terrains/{nom}", a.GetTerrainHandler).Methods("GET")
	api.HandleFunc("/terrains/{id}", a.TerrainsHandler).Methods("DELETE", "PUT")
	api.HandleFunc("/terrains", a.CreerTerrainHandler).Methods("POST")
	// Equipes
	api.HandleFunc("/equipes", a.GetEquipesHandler).Methods("GET")
	api.HandleFunc("/equipes/{nom}", a.GetEquipeHandler).Methods("GET")
	api.HandleFunc("/equipes/{id}", a.EquipesHandler).Methods("DELETE", "PUT")
	api.HandleFunc("/equipes", a.CreerEquipeHandler).Methods("POST")
	// Parties
	api.HandleFunc("/parties", a.PartiesHandler).Methods("GET", "POST")
	api.HandleFunc("/parties/{id}", a.PartieHandler).Methods("GET", "PUT", "OPTIONS")
	api.HandleFunc("/parties/{id}", a.SupprimerPartiesHandler).Methods("DELETE")
	// BD
	api.HandleFunc("/seed", a.RemplirBD).Methods("POST")
	api.HandleFunc("/bd", a.FaireBD).Methods("POST")
	// Actions
	api.HandleFunc("/actions", a.GetActions).Methods("GET")
	api.HandleFunc("/actions", a.PostAction).Methods("POST")
	// Joueurs
	api.HandleFunc("/joueur", a.HandleJoueur).Methods("POST")
	api.HandleFunc("/joueur/{id}", a.HandleJoueur).Methods("PUT", "OPTIONS")
	api.HandleFunc("/joueur", a.GetJoueurs).Methods("GET")
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
