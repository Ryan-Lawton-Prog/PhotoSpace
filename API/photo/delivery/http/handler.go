package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ryanlawton.art/photospace-api/auth"
	"ryanlawton.art/photospace-api/models"
	"ryanlawton.art/photospace-api/photo"
)

type Photo struct {
	PhotoID string `json:"photo_id"`
	Title   string `json:"title"`
	Photo   []byte `json:"photo"`
}

type Handler struct {
	useCase photo.UseCase
}

func NewHandler(useCase photo.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Photo []byte `form:"photo"`
}

func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("photo")
	defer file.Close()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	fileb := make([]byte, header.Size)
	_, err = file.Read(fileb)

	if err != nil {
		log.Printf("Error reading file: %s\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.UploadPhoto(c, user, fileb, header.Filename); err != nil {
		log.Printf("Error uploading photo: %s\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Photo uploaded successfully"})
}

type fetchInput struct {
	PhotoID string `json:"photo_id"`
}

func (h *Handler) Fetch(c *gin.Context) {
	log.Printf("Hello")
	inp := new(fetchInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Printf("Input: %s\n", inp)

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	p, err := h.useCase.FetchPhoto(c.Request.Context(), user, inp.PhotoID)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Printf("Filename: %s\n", p.Filename)
	log.Printf("photo: %s", p.Photo)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", p.Filename))
	c.Data(http.StatusOK, "application/octet-stream", p.Photo)
}

// func NewHandler(useCase photo.UseCase) *Handler {
// 	return &Handler{
// 		useCase: useCase,
// 	}
// }

// func (h *Handler) Upload(c *gin.Context) {
// 	photo, err := c.FormFile("photo")

// 	if err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	f, _ := photo.Open()
// 	fmt.Printf("f: %v\n", f)

// 	//user := c.MustGet(auth.CtxUserKey).(*models.User)

// 	// for _, file := range files {
// 	// 	log.Println(file.Filename)
// 	// 	f, _ := file.Open()
// 	// 	fmt.Printf("f: %v\n", f)
// 	// 	// if err := h.useCase.UploadPhoto(c, user, f.Read()); err != nil {
// 	// 	// 	c.JSON(500, gin.H{"error": err.Error()})
// 	// 	// 	return
// 	// 	// }

// 	// 	// Upload the file to specific dst.
// 	// 	//c.SaveUploadedFile(file, dst)
// 	// }

// 	c.Status(http.StatusCreated)
// }
