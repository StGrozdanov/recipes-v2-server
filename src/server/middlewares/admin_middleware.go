package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"recipes-v2-server/utils"
)

// AdminMiddleware checks if the request comes from user with ADMINISTRATOR role and only then the request is authorized
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Request.Header["X-Authorization"]) == 0 || len(ctx.Request.Header["X-Authorization"][0]) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := ctx.Request.Header["X-Authorization"][0]

		claims, isValid, err := utils.ParseJWT(token)
		if err != nil || !isValid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, Errors{
				Info: Info{
					Message: "Invalid Token",
					Cause:   "Auth Token",
				},
			})
			return
		}

		if claims.Role != "ADMINISTRATOR" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, Errors{Info: Info{Message: "You don't have permissions to access this resource", Cause: "Missing permissions"}})
			return
		}

		ctx.Next()
	}
}
