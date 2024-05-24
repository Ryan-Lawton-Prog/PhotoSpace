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

func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	blob := make(models.PhotoBlob, header.Size)
	_, err = file.Read(blob)

	if err != nil {
		log.Printf("Error reading file: %s\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	photoId, err := h.useCase.UploadPhoto(c, user, &blob, header.Filename)
	if err != nil {
		log.Printf("Error uploading photo: %s\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Photo uploaded successfully", "photo_id": photoId})
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

	pm, blob, err := h.useCase.FetchPhoto(c.Request.Context(), user, inp.PhotoID)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched file name: %s", pm.Filename)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pm.Filename))
	c.Data(http.StatusOK, "application/octet-stream", blob)
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
