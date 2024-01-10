package middlewares

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"net/http"
	"recipes-v2-server/database"
	"recipes-v2-server/internal/comments"
	"recipes-v2-server/utils"
	"strings"
)

type Info struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

type Errors struct {
	Info `json:"error"`
}

// ResourceOwnerMiddleware checks if the request comes from the resource owner or ADMINISTRATOR and only then
// the request is authorized
func ResourceOwnerMiddleware() gin.HandlerFunc {
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

		username, isUsersRelated := ctx.Params.Get("username")
		recipeName, isRecipesRelated := ctx.Params.Get("name")
		isCommentsRelated := strings.Contains(ctx.Request.URL.String(), "/comments")

		if !isUsersRelated && !isRecipesRelated && !isCommentsRelated {
			ctx.AbortWithStatusJSON(http.StatusForbidden, Errors{Info: Info{Message: "Missing identifier", Cause: "Failed to identify resource"}})
			return
		}

		var permissionsAreValid bool

		if isUsersRelated {
			permissionsAreValid = validateUserEditRequest(claims, username)
		} else if isRecipesRelated {
			permissionsAreValid = validateRecipeEditRequest(claims, recipeName)
		} else {
			permissionsAreValid = validateCommentEditRequest(claims, ctx)
		}

		if !permissionsAreValid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, Errors{Info: Info{Message: "You don't have permissions to access this resource", Cause: "Missing permissions"}})
			return
		}

		ctx.Next()
	}
}

func validateCommentEditRequest(claims *utils.TokenClaims, ctx *gin.Context) (validPermissions bool) {
	var ownerId int

	if ctx.Request.Method == "PUT" {
		editData := comments.CommentEditData{}
		if err := ctx.ShouldBind(&editData); err != nil {
			return
		}

		if _, err := validator.ValidateStruct(editData); err != nil {
			return
		}

		err := database.GetSingleRecordNamedQuery(
			&ownerId,
			`SELECT owner_id FROM comments WHERE id = :id`,
			map[string]interface{}{"id": editData.Id},
		)
		if err != nil {
			return
		}
		ctx.Set("commentData", editData)
	} else if ctx.Request.Method == "DELETE" {
		deleteData := comments.CommentIdData{}

		if err := ctx.ShouldBind(&deleteData); err != nil {
			return
		}

		if _, err := validator.ValidateStruct(deleteData); err != nil {
			return
		}

		selectErr := database.GetSingleRecordNamedQuery(
			&ownerId,
			`SELECT owner_id FROM comments WHERE id = :id`,
			map[string]interface{}{"id": deleteData.Id},
		)
		if selectErr != nil {
			return
		}
		ctx.Set("commentData", deleteData)
	}

	if claims.Role != "ADMINISTRATOR" && claims.Role != "MODERATOR" && claims.Id != ownerId {
		return
	}

	return true
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
