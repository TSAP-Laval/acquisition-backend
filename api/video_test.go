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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Simule l'envoie d'une video au format mp4
func TestUploadVideoMP4PourEnvoie(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[0] = m.GameID

	if gameID[0] == "" || gameID[0] == "0" {
		t.Errorf("Game ID expected: %s", gameID[0])
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
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// TestGetVideo permet d'evoyer une vidéo au client
func TestGetVideo(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"/api/parties/"+gameID[0]+"/videos/1", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

}

// Simule la suppression de la première partie (avec la vidéo)
func TestUploadDeleteVideoMP4Envoye(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[0], reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Response code expected: %d", res.StatusCode)
		var m MessageError
		responseMapping(&m, res)
		t.Errorf("Error: %s", m.Err)
	}
}

// Vérifie que le dossier est bel et bien vide
func TestFolderEmptyEnvoie(t *testing.T) {
	folder, err := filepath.Abs(videoPath)
	if err != nil {
		t.Error(err)
	}

	if empty, err := IsDirEmpty(folder); !empty || err != nil {
		t.Errorf("Folder not empty")
		t.Error(err)
		return
	} else {
		if err := os.Remove(folder); err != nil {
			t.Error(err)
		}
	}
}
