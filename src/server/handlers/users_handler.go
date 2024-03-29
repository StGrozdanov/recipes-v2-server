package handlers

import (
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/users"
	"recipes-v2-server/utils"
	"strconv"
)

func GetUser(ginCtx *gin.Context) {
	username, ok := ginCtx.Params.Get("username")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "username was not found"})
		return
	}

	user, err := users.GetUser(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting a user from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, user)
}

func UploadCoverImage(ginCtx *gin.Context) {
	username, found := ginCtx.GetPostForm("username")
	if !found {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"error": "invalid parameters, expected username to be present in the form data"},
		)
		return
	}

	imageKey := fmt.Sprintf("%s-cover-image", username)

	coverImage, err := ginCtx.FormFile(imageKey)
	if err != nil {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"error": fmt.Sprintf("the expected key - %s was not found in the form data", imageKey)},
		)
		return
	}

	imageURL, err := users.UploadCoverImage(coverImage, imageKey, username)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on attempting to upload cover image")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusCreated, map[string]interface{}{"coverImageURL": imageURL})
}

func UploadAvatarImage(ginCtx *gin.Context) {
	username, found := ginCtx.GetPostForm("username")
	if !found {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"error": "invalid parameters, expected username to be present in the form data"},
		)
		return
	}

	imageKey := fmt.Sprintf("%s-avatar-image", username)

	avatarImage, err := ginCtx.FormFile(imageKey)
	if err != nil {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"error": fmt.Sprintf("the expected key - %s was not found in the form data", imageKey)},
		)
		return
	}

	imageURL, err := users.UploadAvatarImage(avatarImage, imageKey, username)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on attempting to upload avatar image")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusCreated, map[string]interface{}{"avatarImageURL": imageURL})
}

func EditUserData(ginCtx *gin.Context) {
	username, ok := ginCtx.Params.Get("username")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "username was not found"})
		return
	}

	data := users.User{}

	if err := ginCtx.ShouldBind(&data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	userData, err := users.EditData(username, data)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such user"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on edit attempt for user %s", username)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, userData)
}

func GetUsersCount(ginCtx *gin.Context) {
	usersCount, err := users.Count()
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting users count")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"count": usersCount})
}

func GetAllUsers(ginCtx *gin.Context) {
	usersData, err := users.GetAllUsers()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting users from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, usersData)
}

func DeleteUser(ctx *gin.Context) {
	userId, ok := ctx.Params.Get("id")

	if !ok {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "user id was not found"})
		return
	}

	userIdAsNumber, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": err.Error()})
		return
	}

	err = users.Delete(userIdAsNumber)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on delete attempt for user %s", userId)

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func ChangeRole(ctx *gin.Context) {
	var data users.UserChangeRoleData

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	err := users.ChangeRole(data)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such user found"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on change role attempt for user %d", data.UserId)

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func BlockUser(ctx *gin.Context) {
	var data users.BlockUserData

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	err := users.Block(data)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such user found"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on block attempt for user %d", data.UserId)

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func UnblockUser(ctx *gin.Context) {
	userId, ok := ctx.Params.Get("id")

	if !ok {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "user id was not found"})
		return
	}

	userIdAsNumber, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": err.Error()})
		return
	}

	err = users.Unblock(userIdAsNumber)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "no such user found"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on unblock attempt for user %s", userId)

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}
