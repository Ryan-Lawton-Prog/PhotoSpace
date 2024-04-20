package http

import (
	"github.com/gin-gonic/gin"
	"ryanlawton.art/photospace-api/photo"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc photo.UseCase) {
	h := NewHandler(uc)

	photos := router.Group("/photo")
	{
		photos.POST("", h.Upload)
		photos.GET("/ids", h.FetchAllIDs)
		// photos.DELETE("", h.Delete)
		photos.GET("", h.Fetch)
	}
}
