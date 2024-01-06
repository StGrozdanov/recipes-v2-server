package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
)

// ResourceOwnerMiddleware checks if the request comes from the resource owner or ADMINISTRATOR and only then
// the request is authorized
func ResourceOwnerMiddleware() gin.HandlerFunc {
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

		claims, isValid, err := utils.ParseJWT(token)
		if err != nil || !isValid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"message": "Invalid token"})
			return
		}

		username, isUsersRelated := ctx.Params.Get("username")
		recipeName, isRecipesRelated := ctx.Params.Get("name")

		if !isUsersRelated && !isRecipesRelated {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"message": "Missing request identifier"})
			return
		}

		var permissionsAreValid bool

		if isUsersRelated {
			permissionsAreValid = validateUserEditRequest(claims, username)
		} else {
			permissionsAreValid = validateRecipeEditRequest(claims, recipeName)
		}

		if !permissionsAreValid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"message": "Missing permissions to access this recipe"})
			return
		}

		ctx.Next()
	}
}

func validateRecipeEditRequest(claims *utils.TokenClaims, recipeName string) (validPermissions bool) {
	var ownerId int

	err := database.GetSingleRecordNamedQuery(
		&ownerId,
		`SELECT owner_id FROM recipes WHERE recipe_name = :recipe_name`,
		map[string]interface{}{"recipe_name": recipeName},
	)
	if err != nil {
		return
	}

	if claims.Role != "ADMINISTRATOR" && claims.Role != "MODERATOR" && claims.Id != ownerId {
		return
	}

	return true
}

func validateUserEditRequest(claims *utils.TokenClaims, username string) (validPermissions bool) {
	var ownerId int

	err := database.GetSingleRecordNamedQuery(
		&ownerId,
		`SELECT id FROM users WHERE username = :username`,
		map[string]interface{}{"username": username},
	)
	if err != nil {
		return
	}

	if claims.Id != ownerId && claims.Role != "ADMINISTRATOR" {
		return
	}

	return true
}
