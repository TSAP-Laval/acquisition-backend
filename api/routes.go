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

// AcquisitionService represents a single service instance
type AcquisitionService struct {
	logger *log.Logger
	config *AcquisitionConfiguration
	server *manners.GracefulServer
}

// New crée une nouvelle instance du service
func New(writer io.Writer, config *AcquisitionConfiguration) *AcquisitionService {

	return &AcquisitionService{
		logger: log.New(writer, "[acquisition-api] ", log.Flags()),
		config: config,
		server: manners.NewServer(),
	}
}

// Info écrit un message vers le logger du service
func (a *AcquisitionService) Info(message string) {
	a.logger.Println(message)
}

// Info écrit un message d'erreur vers le logger du service
func (a *AcquisitionService) Error(message string) {
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

	// Actions
	r.HandleFunc("/api/action/movementType", a.GetMovementTypeHandler).Methods("GET")
	r.HandleFunc("/api/action/actiontype", a.GetAllActionsTypes).Methods("GET")
	r.HandleFunc("/api/action/addactiontype", a.PostActionType).Methods("POST")
	//Coachs
	r.HandleFunc("/api/coachs/coachs", a.GetCoachsHandler).Methods("GET")
	r.HandleFunc("/api/coachs/addcoach", a.PostCoachHandler).Methods("POST")
	r.HandleFunc("/api/coachs/addCoachTeam/{id}", a.AssignerEquipeCoach).Methods("PUT")
	// Upload
	r.HandleFunc("/api/upload", a.UploadHandler)
	// Terrains
	r.HandleFunc("/api/terrains", a.GetTerrainsHandler).Methods("GET")
	r.HandleFunc("/api/terrains/{nom}", a.GetTerrainHandler).Methods("GET")
	r.HandleFunc("/api/terrains/{id}", a.TerrainsHandler).Methods("DELETE", "PUT")
	r.HandleFunc("/api/terrains", a.CreerTerrainHandler).Methods("POST")
	// Equipes
	r.HandleFunc("/api/equipes", a.GetEquipesHandler).Methods("GET")
	r.HandleFunc("/api/equipes/{nom}", a.GetEquipeHandler).Methods("GET")
	r.HandleFunc("/api/equipes/{id}", a.EquipesHandler).Methods("DELETE", "PUT")
	r.HandleFunc("/api/equipes", a.CreerEquipeHandler).Methods("POST")
	// Parties
	r.HandleFunc("/api/parties", a.PartiesHandler).Methods("GET", "POST")
	r.HandleFunc("/api/parties/{id}", a.SupprimerPartiesHandler).Methods("DELETE")
	// BD
	r.HandleFunc("/api/seed", a.RemplirBD).Methods("POST")
	r.HandleFunc("/api/bd", a.FaireBD).Methods("POST")
	// Autre
	r.HandleFunc("/api/actions", a.GetActions).Methods("GET")
	r.HandleFunc("/api/actions", a.PostAction).Methods("POST")
	r.HandleFunc("/api/joueur", a.HandleJoueur).Methods("POST")
	r.HandleFunc("/api/joueur/{id}", a.HandleJoueur).Methods("PUT")
	r.HandleFunc("/api/saison", a.GetSeasons).Methods("GET")
	r.HandleFunc("/api/saison", a.PostSaison).Methods("POST")
	r.HandleFunc("/api/sports", a.GetSports).Methods("GET")
	r.HandleFunc("/api/niveau", a.GetNiveau).Methods("GET")
	r.HandleFunc("/api/joueur", a.GetJoueurs).Methods("GET")
	return a.Middleware(r)
}

// Start démarre le service
func (a *AcquisitionService) Start() {
	go func() {
		a.server.Addr = a.config.Port
		a.server.Handler = a.getRouter()
		err := a.server.ListenAndServe()
		a.Info("Acquisition shutting down...")

		if err != nil {
			panic(err)
		}

	}()
	a.logger.Printf("TSAP-Acquisiton started on localhost%s... \n", a.config.Port)
}

// Stop arrête le service
func (a *AcquisitionService) Stop() {
	a.server.Close()
}
