package http

import (
	"encoding/json"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ryanlawton.art/photospace-api/auth"
	"ryanlawton.art/photospace-api/models"
	"ryanlawton.art/photospace-api/photo/usecase"
)

func TestUpload(t *testing.T) {
	testUser := &models.User{
		Username: "testuser",
		Password: "testpass",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.PhotoUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	inp := &createInput{
		Photo: []byte("testphoto"),
	}

	_, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc.On("UploadPhoto", testUser, inp.Photo).Return(nil)

	// encode image to it's type
	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("photo", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := image.NewRGBA(image.Rect(0, 1, 1, 0))

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/photos", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
