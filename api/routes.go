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

	r.HandleFunc("/api/GetMovementType", a.GetMovementTypeHandler)
	r.HandleFunc("/api/GetActionType", a.GetAllActionsTypes)
	r.HandleFunc("/api/PostActionType", a.PostActionType)
	r.HandleFunc("/api/video", a.VideoHandler)
	r.HandleFunc("/api/edition/GetJoueurs", a.GetJoueurs)
	r.HandleFunc("/api/seeders", a.Remplir)
	r.HandleFunc("/api/seeders/FaireBD", a.FaireBD)
	r.HandleFunc("/api/edition/GetActions", a.GetActions)
	r.HandleFunc("/api/edition/PostJoueur", a.PostJoueur)

	return a.Middleware(r)
}

// Start démarre le service
func (a *AcquisitionService) Start() {
	go func() {
		a.server.Addr = a.config.Port
		a.server.Handler = a.getRouter()
		a.server.ListenAndServe()
		a.Info("Acquisition shutting down...")
	}()
	a.logger.Printf("TSAP-Acquisiton started on localhost%s... \n", a.config.Port)
}

// Stop arrête le service
func (a *AcquisitionService) Stop() {
	a.server.Close()
}
