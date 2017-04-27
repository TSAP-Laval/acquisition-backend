//
// Fichier     : messages.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de gérer les types de messages à envoyer
//

package api

type simpleMessage struct {
	Body string `json:"body"`
}

type errorMessage struct {
	Error string `json:"error"`
}
