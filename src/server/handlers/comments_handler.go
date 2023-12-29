package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/comments"
	"recipes-v2-server/utils"
)

func GetLatestComments(ginCtx *gin.Context) {
	latestComments, err := comments.GetLatestComments()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting latest comments from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, latestComments)
}

func GetRecipeComments(ginCtx *gin.Context) {
	recipeName, ok := ginCtx.Params.Get("recipeName")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "recipe name was not found"})
		return
	}

	recipeComments, err := comments.GetCommentsForRecipe(recipeName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting recipe comments from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, recipeComments)
}
