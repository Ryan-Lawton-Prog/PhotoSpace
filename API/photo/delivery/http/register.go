package http

import (
	"github.com/gin-gonic/gin"
	"ryanlawton.art/photospace-api/photo"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc photo.UseCase) {
	h := NewHandler(uc)

	photos := router.Group("/photos")
	{
		photos.POST("", h.Upload)
		photos.GET("", h.Fetch)
		// photos.DELETE("", h.Delete)
	}
}
