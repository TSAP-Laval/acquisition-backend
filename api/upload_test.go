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

// Constantes sur les tailles possibles
const (
	GB = 1 << (10 * 3)
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

var gameID [6]string

// TestUploadVideoMP4MauvaiseDate Simule l'envoie d'une video au format mp4
// avec une mauvaise date
func TestUploadVideoMP4MauvaiseDate(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[0] = m.GameID

	if gameID[0] == "" || gameID[0] == "0" {
		LogErrors(Messages{t, "Game ID expected: %s", gameID[0], false, nil, nil})
	}

	path, err := filepath.Abs(testPath + "/small.mp4")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[0])
	if err != nil {
		t.Error(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var me MessageError
	err = json.Unmarshal(bodyBuffer, &me)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	if !strings.Contains(me.Err, "Une erreur est survenue lors de l'envoie de la vidéo!") {
		t.Error("Error expected : ", me.Err)
	}
}

// Simule l'envoie d'aucune video
func TestUploadNoVideo(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	var gameID = m.GameID

	if gameID == "" || gameID == "0" {
		t.Errorf("Game ID expected: %s", gameID)
	}

	reader = strings.NewReader("")
	request, err = http.NewRequest("POST", baseURL+"/api/upload/"+gameID, reader)

	if err != nil {
		t.Error(err)
	}

	me := BadRequestHandler(request, t)

	if !strings.Contains(me.Err, "Aucun fichier n'a été envoyé ! Veuillez réessayer !") {
		LogErrors(Messages{t, "Error expected: %s", me.Err, true, request, res})
	}
}

// Simule l'envoie d'une vidéo vide de contenu
func TestUploadEmptyVideo(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	var gameID = m.GameID

	if gameID == "" || gameID == "0" {
		t.Errorf("Game ID expected: %s", gameID)
	}

	reader = strings.NewReader("")
	request, err = http.NewRequest("POST", baseURL+"/api/upload/"+gameID, reader)
	res, err = SecureRequest(request)

	path, err := filepath.Abs(testPath + "/empty.mp4")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "Le fichier envoyé ne possède aucun contenu!") {
		LogErrors(Messages{t, "Error expected: %s", me.Err, true, request, res})
	}
}

// TestUploadVideoMP4MauvaisGameID Simule l'envoie d'une video au format mp4
// avec un mauvais identifiant de partie
func TestUploadVideoMP4MauvaisGameID(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[0] = m.GameID

	if gameID[0] == "" || gameID[0] == "0" {
		LogErrors(Messages{t, "Game ID expected: %s", gameID[0], false, nil, nil})
	}

	path, err := filepath.Abs(testPath + "/small.mp4")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, "aaa")
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// TestUploadVideoMP4ErrBD Simule l'envoie d'une video au format mp4
// avec erreur de connexion à la base de données
func TestUploadVideoMP4ErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/upload/1", reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// Simule l'envoie de toutes les vidéos
func TestUploadAllVideos(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[0] = m.GameID

	if gameID[0] == "" || gameID[0] == "0" {
		LogErrors(Messages{t, "Game ID expected: %s", gameID[0], false, nil, nil})
	}

	pathMp4, err := filepath.Abs(testPath + "/small.mp4")
	if err != nil {
		t.Error(err)
	}

	pathWebm, err := filepath.Abs(testPath + "/small.webm")
	if err != nil {
		t.Error(err)
	}

	pathMov, err := filepath.Abs(testPath + "/small.mov")
	if err != nil {
		t.Error(err)
	}

	pathMpg, err := filepath.Abs(testPath + "/small.mpg")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{pathMp4, pathWebm, pathMov, pathMpg, pathMp4, pathMp4, pathMp4}, t, gameID[0])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule l'envoie d'une video au format mp4
func TestUploadVideoMP4(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[0] = m.GameID

	if gameID[0] == "" || gameID[0] == "0" {
		LogErrors(Messages{t, "Game ID expected: %s", gameID[0], false, nil, nil})
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

// Simule l'envoie d'une video au format 3gp
func TestUploadVideo3GP(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
		return
	}

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[1] = m.GameID

	if gameID[1] == "" || gameID[1] == "0" {
		t.Errorf("Game ID expected: %s", gameID[1])
	}

	invalidFormatVideoRequest(request, "/small.3gp", gameID[1], t)
}

// Simule l'envoie d'une video au format flv
func TestUploadVideoFLV(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
		return
	}

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[2] = m.GameID

	if gameID[2] == "" || gameID[2] == "0" {
		t.Errorf("Game ID expected: %s", gameID[2])
	}

	invalidFormatVideoRequest(request, "/small.flv", gameID[2], t)
}

// Simule l'envoie d'une video au format webm
func TestUploadVideoWEBM(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

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
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule l'envoie d'une video au format ogv
func TestUploadVideoOGV(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
		return
	}

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[3] = m.GameID

	if gameID[3] == "" || gameID[3] == "0" {
		t.Errorf("Game ID expected: %s", gameID[3])
	}

	invalidFormatVideoRequest(request, "/small.ogv", gameID[3], t)
}

// Simule l'envoie d'une video au format mpg
func TestUploadVideoMPG(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[4] = m.GameID

	if gameID[4] == "" || gameID[4] == "0" {
		t.Errorf("Game ID expected: %s", gameID[4])
	}

	path, err := filepath.Abs(testPath + "/small.mpg")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[4])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule l'envoie d'une video au format mov
func TestUploadVideoMOV(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[5] = m.GameID

	if gameID[5] == "" || gameID[5] == "0" {
		t.Errorf("Game ID expected: %s", gameID[5])
	}

	path, err := filepath.Abs(testPath + "/small.mov")
	if err != nil {
		t.Error(err)
	}

	res, err = sendUploadRequest([]string{path}, t, gameID[5])
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule l'envoie d'une video au format wmv
func TestUploadVideoWMV(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
		return
	}

	var m MessageSuccess
	responseMapping(&m, res)
	gameID[2] = m.GameID

	if gameID[2] == "" || gameID[2] == "0" {
		t.Errorf("Game ID expected: %s", gameID[2])
	}

	invalidFormatVideoRequest(request, "/small.wmv", gameID[2], t)
}

// Simule l'envoie de 2 videos au format mp4 et webm
func TestUploadVideos(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

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
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule l'envoie d'une video au format webm pour suppression
func TestUploadVideoWEBMDel(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("POST", baseURL+"/api/parties", reader)
	res, err := SecureRequest(request)

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
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule l'envoie d'une requestuête d'options pour l'upload
func TestSendOptionsUpload(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("OPTIONS", baseURL+"/api/parties/0", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}
}

// Simule la suppression d'une fausse vidéo
func TestUploadDeleteFalseVideo(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/0", reader)
	res, err := SecureRequest(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	var m MessageError
	responseMapping(&m, res)

	if !strings.Contains(m.Err, "aucune partie ne correspond") {
		t.Errorf("Error expected: %s", m.Err)
	}
}

// TestUploadDeleteVideoMP4ErrBD Simule la suppression de la première partie (avec la vidéo)
// avec erreur de connexion à la base de données
func TestUploadDeleteVideoMP4ErrBD(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[0], reader)

	if err != nil {
		t.Error(err)
	}

	BDErrorHandler(request, t)
}

// Simule la suppression de la première partie (avec la vidéo)
func TestUploadDeleteVideoMP4(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[0], reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}

// Simule la suppression de la deuxième partie (avec la vidéo)
func TestUploadDeleteVideoWEBM(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[1], reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}

// Simule la suppression de la troisième partie (avec les vidéos)
func TestUploadDeleteVideos(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[2], reader)

	if err != nil {
		t.Error(err)
		return
	}

	DeleteHandler(request, t)
}

// Simule la suppression de la dernière partie (avec la vidéo supprimé `à bras`)
func TestUploadDeleteVideoWEBMDel(t *testing.T) {
	removeContents(videoPath)

	reader = strings.NewReader("")
	request, err := http.NewRequest("DELETE", baseURL+"/api/upload/"+gameID[3], reader)
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

	if !strings.Contains(m.Err, "Impossible de supprimer la video !") {
		t.Errorf("Error expected: %s", m.Err)
	}
}

// Vérifie que le dossier est bel et bien vide
func TestFolderEmpty(t *testing.T) {
	FolderEmpty(t)
}

func FolderEmpty(t *testing.T) {
	folder, err := filepath.Abs(videoPath)
	if err != nil {
		t.Error(err)
	}

	if empty, err := isDirEmpty(folder); !empty || err != nil {
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
func isDirEmpty(name string) (bool, error) {
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
	request, err := newfileUploadRequest(path, gameID)
	if err != nil {
		t.Error(err)
	}

	return SecureRequest(request)
}

func responseMapping(m interface{}, res *http.Response) {
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	json.Unmarshal(body, m)
}

// Créé une requestuête pour l'envoie de fichier
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

		var part io.Writer
		if gameID == "5" {
			part, err = writer.CreateFormFile("Mon Jan 02 2006 :04:05sdfsdf GMT-0400", fi.Name())
		} else {
			part, err = writer.CreateFormFile(fi.ModTime().Format("Mon Jan 02 2006 15:04:05 GMT-0400 (EDT)"), fi.Name())
		}

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

func invalidFormatVideoRequest(request *http.Request, path string, gameID string, t *testing.T) {
	path, err := filepath.Abs(testPath + path)
	if err != nil {
		t.Error(err)
	}

	res, err := sendUploadRequest([]string{path}, t, gameID)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 400 {
		LogErrors(Messages{t, "Response code expected: %d", res.StatusCode, true, request, res})
	}

	var me MessageError
	responseMapping(&me, res)

	if !strings.Contains(me.Err, "n'est pas une vidéo de format valide") {
		t.Errorf("Error expected: %s", me.Err)
	}
}
