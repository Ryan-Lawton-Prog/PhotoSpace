package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	token "ryanlawton.art/photospace/internal/pkg/models"
)

const (
	photoSpaceURL = "http://localhost:8000"
)

type SignInBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	JWT token.JWT `json:"token"`
}

func SignIn(body SignInBody) (bool, error) {
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(photoSpaceURL+"/auth/sign-in", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		token, _ := token.GetInstance()
		resBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return false, err
		}
		respJson := SignInResponse{}
		json.Unmarshal(resBody, &respJson)
		token.JWT = respJson.JWT

		return true, nil
	default:
		return false, errors.New("failed login")
	}
}
