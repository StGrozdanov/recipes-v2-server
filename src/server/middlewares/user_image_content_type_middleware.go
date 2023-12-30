package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ImageContentTypeMiddleware checks for valid content type header of images
func ImageContentTypeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, found := ctx.GetPostForm("username")
		if !found {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				map[string]interface{}{"error": "invalid parameters, expected username to be present in the form data"},
			)
			return
		}

		coverImageKey := fmt.Sprintf("%s-cover-image", username)
		avatarImageKey := fmt.Sprintf("%s-avatar-image", username)

		file, err := ctx.FormFile(coverImageKey)
		if err != nil {
			file, err = ctx.FormFile(avatarImageKey)
			if err != nil {
				ctx.AbortWithStatusJSON(
					http.StatusBadRequest,
					map[string]interface{}{"message": "the provided file key is not in the expected format."},
				)
				return
			}
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
