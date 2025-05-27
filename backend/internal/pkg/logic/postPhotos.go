package logic

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	token "ryanlawton.art/photospace/internal/pkg/models"
)

type PostPhotoBody struct {
	Photo []byte `json:"photo"`
}

func PostPhoto(photo []byte, fileName string) error {
	client := &http.Client{}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	part, err := w.CreateFormFile("photo", fileName)
	if err != nil {
		return err
	}

	part.Write(photo)
	w.Close()

	req, err := http.NewRequest("POST", photoSpaceURL+"/api/photo", buf)
	if err != nil {
		return err
	}

	token, _ := token.GetInstance()
	req.Header.Add("Authorization", "Bearer "+string(token.JWT))
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		_, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("failed to fetch ids")
	}
}
