package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"ryanlawton.art/photospace/internal/app/models/user"
)

type GetIDsResponse struct {
	PhotoIDs []string `json:"photo_ids"`
}

func GetPhotoIDs() ([]string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", photoSpaceURL+"/api/photo/ids", nil)
	if err != nil {
		return nil, err
	}

	user, _ := user.GetInstance()
	req.Header.Add("Authorization", "Bearer "+string(user.JWT))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		resBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		respJson := GetIDsResponse{}
		json.Unmarshal(resBody, &respJson)

		return respJson.PhotoIDs, nil
	default:
		return nil, errors.New("failed to fetch ids")
	}
}

type GetPhotoBody struct {
	PhotoId string `json:"photo_id"`
}

func GetPhoto(id string) ([]byte, error) {
	client := &http.Client{}
	jsonBody, _ := json.Marshal(&GetPhotoBody{PhotoId: id})
	req, err := http.NewRequest("GET", photoSpaceURL+"/api/photo", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	user, _ := user.GetInstance()
	req.Header.Add("Authorization", "Bearer "+string(user.JWT))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		resBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		return resBody, nil
	default:
		return nil, errors.New("failed to fetch ids")
	}
}

func GetPhotoThumbnail(id string) ([]byte, error) {
	client := &http.Client{}
	jsonBody, _ := json.Marshal(&GetPhotoBody{PhotoId: id})
	req, err := http.NewRequest("GET", photoSpaceURL+"/api/photo/thumbnail", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	user, _ := user.GetInstance()
	req.Header.Add("Authorization", "Bearer "+string(user.JWT))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		resBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		return resBody, nil
	default:
		return nil, errors.New("failed to fetch ids")
	}
}
