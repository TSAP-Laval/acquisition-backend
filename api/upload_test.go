//
// TEST
//
// Fichier     : upload_test.go
// Développeur : Laurent Leclerc Poulin
//
// Permet de tester l'upload de fichiers.
//

package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type MessageError struct {
	Err string `json:"error"`
}

type MessageSuccess struct {
	Success string `json:"success"`
	GameID  string `json:"game_id"`
}

type MessageGameID struct {
	GameID string `json:"game_id"`
}

const videoPath = "../videos"
const testPath = "../test"

var gameID [4]string

// Simule l'envoie d'aucune video
func TestUploadNoVideo(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	var gameID = m.GameID

	if gameID == "" || gameID == "0" {
		t.Errorf("Game ID expected: %s", gameID)
	}

	reader = strings.NewReader("")
	req, err = http.NewRequest("POST", baseURL+"/api/upload/"+gameID, reader)
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "Aucun fichier envoyé") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// Simule l'envoie d'une video au format mp4
func TestUploadVideoMP4(t *testing.T) {
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

// Simule l'envoie d'une video au format 3gp
func TestUploadVideo3GP(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[1] = m.GameID

	if gameID[1] == "" || gameID[1] == "0" {
		t.Errorf("Game ID expected: %s", gameID[1])
	}

	path, err := filepath.Abs(testPath + "/small.3gp")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[1])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "n'est pas une vidéo de format valide") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// Simule l'envoie d'une video au format flv
func TestUploadVideoFLV(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[2] = m.GameID

	if gameID[2] == "" || gameID[2] == "0" {
		t.Errorf("Game ID expected: %s", gameID[2])
	}

	path, err := filepath.Abs(testPath + "/small.flv")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[2])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "n'est pas une vidéo de format valide") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// Simule l'envoie d'une video au format ogv
func TestUploadVideoOGV(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[3] = m.GameID

	if gameID[3] == "" || gameID[3] == "0" {
		t.Errorf("Game ID expected: %s", gameID[3])
	}

	path, err := filepath.Abs(testPath + "/small.ogv")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[3])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "n'est pas une vidéo de format valide") {
		t.Errorf("Error expected: %s", me.Err)
	}
}

// Simule l'envoie d'une video au format webm
func TestUploadVideoWEBM(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[1] = m.GameID

	if gameID[1] == "" || gameID[1] == "0" {
		t.Errorf("Game ID expected: %s", gameID[1])
	}

	path, err := filepath.Abs(testPath + "/small.webm")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[1])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// Simule l'envoie de 2 videos au format mp4 et webm
func TestUploadVideos(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[2] = m.GameID

	if gameID[2] == "" || gameID[2] == "0" {
		t.Errorf("Game ID expected: %s", gameID[2])
	}

	pathMp4, err := filepath.Abs(testPath + "/small.mp4")
	if err != nil {
		t.Error(err)
	}

	pathWebm, err := filepath.Abs(testPath + "/small.webm")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{pathMp4, pathWebm}, t, gameID[2])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// Simule l'envoie d'une video au format webm pour suppression
func TestUploadVideoWEBMDel(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := http.DefaultClient.Do(req)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[3] = m.GameID

	if gameID[3] == "" || gameID[3] == "0" {
		t.Errorf("Game ID expected: %s", gameID[3])
	}

	path, err := filepath.Abs(testPath + "/small.webm")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[3])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// Simule l'envoie d'une requête d'options pour l'upload
func TestSendOptionsUpload(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("OPTIONS", baseURL+"/api/parties/0", reader)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}
}

// Simule la suppression d'une fausse vidéo
func TestUploadDeleteFalseVideo(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("DELETE", baseURL+"/api/upload/0", reader)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Errorf("Response code expected: %d", res.StatusCode)
	}

	var m MessageError
	responseMapping(&m, res)

	if !strings.Contains(m.Err, "aucune partie ne correspond") {
		t.Errorf("Error expected: %s", m.Err)
	}
}

// Simule la suppression de la première partie (avec la vidéo)
func TestUploadDeleteVideoMP4(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[0], reader)
	res, err := http.DefaultClient.Do(req)

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

// Simule la suppression de la deuxième partie (avec la vidéo)
func TestUploadDeleteVideoWEBM(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[1], reader)
	res, err := http.DefaultClient.Do(req)

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

// Simule la suppression de la troisième partie (avec les vidéos)
func TestUploadDeleteVideos(t *testing.T) {
	reader = strings.NewReader("")
	req, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[2], reader)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
		var m MessageError
		responseMapping(&m, res)
		t.Errorf("Error: %s", m.Err)
	}
}

// Simule la suppression de la dernière partie (avec la vidéo supprimé `à bras`)
func TestUploadDeleteVideoWEBMDel(t *testing.T) {
	removeContents(videoPath)

	reader = strings.NewReader("")
	req, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[3], reader)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Errorf("Response code expected: %d", res.StatusCode)
		return
	}

	var m MessageError
	responseMapping(&m, res)

	if !strings.Contains(m.Err, "Impossible de supprimer la video !") {
		t.Errorf("Error expected: %s", m.Err)
	}
}

// Vérifie que le dossier est bel et bien vide
func TestFolderEmpty(t *testing.T) {
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

// Vérifie que le dossier passé en paramètre est bel et bien vide
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// Remove all files in directory
func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Envoie la requête d'upload du fichier
func sendUploadRequest(path []string, t *testing.T, gameID string) (*http.Response, error) {
	req, err := newfileUploadRequest(path, gameID)
	if err != nil {
		t.Error(err)
	}

	return http.DefaultClient.Do(req)
}

func responseMapping(m interface{}, res *http.Response) {
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, m)
}

// Créé une requête pour l'envoie de fichier
func newfileUploadRequest(paths []string, gameID string) (*http.Request, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for _, path := range paths {

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}
		file.Close()

		part, err := writer.CreateFormFile("file", fi.Name())
		if err != nil {
			return nil, err
		}
		part.Write(fileContents)
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", baseURL+"/api/upload/"+gameID, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}
