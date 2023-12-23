package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ImageContentTypeMiddleware checks for valid content type header of images
func ImageContentTypeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("image")
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				map[string]interface{}{"message": "the provided file should be with the name 'image'."},
			)
			return
		}

		contentTypeIsOfTypeImage := strings.Contains(file.Header.Get("Content-Type"), "image/")

		if !contentTypeIsOfTypeImage {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				map[string]interface{}{"message": "provided file can only be of type image"},
			)
			return
		}

		ctx.Next()
	}
}
