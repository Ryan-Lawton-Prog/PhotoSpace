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
	inp := new(fetchInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	p, err := h.useCase.FetchPhoto(c.Request.Context(), user, inp.PhotoID)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", p.Filename))
	c.Data(http.StatusOK, "application/octet-stream", p.Photo)
}

func (h *Handler) FetchAllIDs(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	ids, err := h.useCase.FetchPhotoAllIDs(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"photo_ids": ids})
}
