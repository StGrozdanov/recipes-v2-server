package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"recipes-v2-server/utils"
)

// AuthMiddleware checks for valid x-authorization access token in the request headers
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Request.Header["X-Authorization"]) == 0 || len(ctx.Request.Header["X-Authorization"][0]) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := ctx.Request.Header["X-Authorization"][0]

		_, isValid, err := utils.ParseJWT(token)
		if err != nil || !isValid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"message": "Invalid token"})
			return
		}

		ctx.Next()
	}
}
