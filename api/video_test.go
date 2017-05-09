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
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

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
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

}
