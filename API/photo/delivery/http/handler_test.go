package http

import (
	"bytes"
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
	// Set up test data
	testUser := &models.User{
		Username: "testuser",
		Password: "testpass",
	}

	buf := new(bytes.Buffer)
	img := image.NewRGBA(image.Rect(0, 1, 1, 0))
	err := png.Encode(buf, img)
	if err != nil {
		t.Error(err)
	}

	// Set up router
	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.PhotoUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

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

	inp := &createInput{
		Photo: buf.Bytes(),
	}

	// Mock the use case
	uc.On("UploadPhoto", testUser, inp.Photo, "someimg.png").Return(nil)

	// Create a new request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/photos", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	// Test implementation
	_, err = json.Marshal(inp)
	assert.NoError(t, err)

	assert.Equal(t, 201, w.Code)
}
