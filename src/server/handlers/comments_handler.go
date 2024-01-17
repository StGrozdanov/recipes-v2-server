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
	dataFromContext, ok := ginCtx.Get("commentData")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "missing parameters"})
		return
	}

	parsedCommentData, ok := dataFromContext.(comments.CommentEditData)
	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid comment data type"})
		return
	}

	commentData, err := comments.Edit(parsedCommentData)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such comment"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on edit attempt for comment with id: %d", parsedCommentData.Id)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, commentData)
}

func DeleteComment(ginCtx *gin.Context) {
	dataFromContext, ok := ginCtx.Get("commentData")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "missing parameters"})
		return
	}

	commentData, ok := dataFromContext.(comments.CommentIdData)
	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid comment data type"})
		return
	}

	err := comments.Delete(commentData.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such comment"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on delete attempt for comment with id: %d", commentData.Id)

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

func GetCommentsCount(ginCtx *gin.Context) {
	commentsCount, err := comments.Count()
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting comments count")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"count": commentsCount})
}
