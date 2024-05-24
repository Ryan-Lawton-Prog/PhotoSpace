package http

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"log"
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

type successResponse map[string]string
type fetchRequest map[string]interface{}

func createHTTPFetchRequest(t *testing.T, img *image.RGBA) (*http.Request, *httptest.ResponseRecorder) {
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

	// Create a new request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/photo", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, w
}

func TestUpload(t *testing.T) {
	// Set up test data
	testUser := &models.User{
		Username: "testuser",
		Password: "testpass",
	}

	// Create image
	buf := new(bytes.Buffer)
	img := image.NewRGBA(image.Rect(0, 1, 1, 0))
	err := png.Encode(buf, img)
	if err != nil {
		t.Error(err)
	}
	blob := new(models.PhotoBlob)
	*blob = models.PhotoBlob(buf.Bytes())

	// Set up router
	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.PhotoUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	// Mock the use case
	uc.On("UploadPhoto", testUser, blob, "someimg.png").Return("123", nil)

	req, res := createHTTPFetchRequest(t, img)
	r.ServeHTTP(res, req)

	// Test implementation
	_, err = json.Marshal(blob)
	assert.NoError(t, err)

	var result successResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalln(err)
	}

	responseTest := successResponse{
		"message":  "Photo uploaded successfully",
		"photo_id": "123",
	}

	assert.Equal(t, 201, res.Code)
	assert.Equal(t, responseTest, result)
}

func TestFetch(t *testing.T) {
	// Set up test data
	testUser := &models.User{
		Username: "testuser",
		Password: "testpass",
	}

	testRequestBody := fetchRequest{
		"photo_id": "123",
	}
	body, _ := json.Marshal(testRequestBody)

	// Create image
	buf := new(bytes.Buffer)
	img := image.NewRGBA(image.Rect(0, 1, 1, 0))
	err := png.Encode(buf, img)
	if err != nil {
		t.Error(err)
	}
	blob := new(models.PhotoBlob)
	*blob = models.PhotoBlob(buf.Bytes())

	// Set up router
	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.PhotoUseCaseMock)
	RegisterHTTPEndpoints(group, uc)

	// Mock the use case
	mockMetadata := &models.PhotoMetadata{
		ID:       "123",
		Filename: "",
	}
	uc.On("FetchPhoto", testUser, "123").Return(mockMetadata, *blob, nil)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/photo", bytes.NewReader(body))
	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, buf.Bytes(), res.Body.Bytes())
}
