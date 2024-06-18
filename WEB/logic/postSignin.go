package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"ryanlawton.art/photospace-ui/models/user"
)

const (
	photoSpaceURL = "http://localhost:8000"
)

type SignInBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	JWT user.JWT `json:"token"`
}

func SignIn(body SignInBody) (bool, error) {
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(photoSpaceURL+"/auth/sign-in", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		user, _ := user.GetInstance()
		resBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return false, err
		}
		respJson := SignInResponse{}
		json.Unmarshal(resBody, &respJson)
		user.JWT = respJson.JWT

		return true, nil
	default:
		return false, errors.New("failed login")
	}
}
