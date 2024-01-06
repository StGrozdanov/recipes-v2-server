package handlers

import (
	validator "github.com/asaskevich/govalidator"
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

func EditComment(ginCtx *gin.Context) {
	data := comments.CommentData{}

	if err := ginCtx.ShouldBind(&data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	commentData, err := comments.Edit(data)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such comment"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on edit attempt for comment with id: %d", data.Id)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, commentData)
}

func DeleteComment(ginCtx *gin.Context) {
	id, ok := ginCtx.Params.Get("id")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "comment id was not found"})
		return
	}

	err := comments.Delete(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such comment"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on delete attempt for comment with id: %s", id)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func CreateComment(ginCtx *gin.Context) {
	data := comments.CommentData{}

	if err := ginCtx.ShouldBind(&data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	commentData, err := comments.Create(data)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on create attempt for comment with content: %s", data.Content)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusCreated, commentData)
}
