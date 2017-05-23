//
// TEST
//
// Fichier     : videos_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de tester l'envoie d'une vidéo au client
//

package api_test

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

// Simule l'envoie d'une video au format mp4
func TestUploadVideoMP4PourEnvoie(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageGameID
	responseMapping(&m, res)
	gameID[0] = m.GameID

	if gameID[0] == "" || gameID[0] == "0" {
		t.Errorf("Game ID expected: %s", gameID[0])
	}

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	path, err := filepath.Abs(testPath + "/small.mp4")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[0])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestGetVideoErrBD permet d'evoyer une vidéo au client
// avec erreur de connexion à la base de données
func TestGetVideoErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+gameID[0]+"/videos/1", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// TestGetVideo permet d'evoyer une vidéo au client
func TestGetVideo(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+gameID[0]+"/videos/1", reader)

	if err != nil {
		t.Error(err)
	}

	GetRequestHandler(request, t)
}

// TestGetVideoInexistante test la récupération d'une vidéo inexistante dans la base de données
func TestGetVideoInexistante(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+gameID[0]+"/videos/10", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
		return
	}

	var m MessageError
	responseMapping(&m, res)

	if !strings.Contains(m.Err, "Fichier inexistant") {
		t.Errorf("Error expected: %s", m.Err)
	}
}

// Simule la suppression de la première partie (avec la vidéo)
func TestUploadDeleteVideoMP4Envoye(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[0], reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}

// Vérifie que le dossier est bel et bien vide
func TestFolderEmptyEnvoie(t *testing.T) {
	FolderEmpty(t)
}

// TestFermetureServeur ferme le l'api
func TestFermetureServeur(t *testing.T) {
	service.Stop()
}
