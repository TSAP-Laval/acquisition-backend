package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/TSAP-Laval/acquisition-backend/api"
)

func PostActionType(t *testing.T) {
	reader = strings.NewReader(`{"Name": "Passe", "Description": "Interceptée", "ControlType": "Neutre", "MovementType": "Postif"}`)
	request, err := http.NewRequest("POST", baseURL+"//addactiontype", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	// Buffer the body
	bodyBuffer, _ := ioutil.ReadAll(res.Body)
	t.Logf("Res: --> %s\n\n", bodyBuffer)

	var l api.Coaches
	err = json.Unmarshal(bodyBuffer, &l)
	if err != nil {
		t.Logf("ERR: --> %s\n\n", err)
	}

	rmID = fmt.Sprintf("%d", l.ID)
	t.Logf("ID: --> %s\n\n", rmID)

	if res.StatusCode != 201 {
		t.Errorf("Action posted: %d", res.StatusCode)
	}
}

func GetActionsTypes(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"//actiontype", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Error
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode) //Success
	}
}

func GetMovementType(t *testing.T) {
	reader = strings.NewReader("")
	request, err := http.NewRequest("GET", baseURL+"//movementType", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Error
	}

	if res.StatusCode != 200 {
		t.Errorf("Get coachs success: %d", res.StatusCode) //Success
	}
}
