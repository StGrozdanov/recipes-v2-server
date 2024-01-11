package handlers

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/notifications"
	"recipes-v2-server/utils"
)

func GetNotifications(ginCtx *gin.Context) {
	username, ok := ginCtx.Params.Get("username")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "username was not found"})
		return
	}

	notificationsResults, err := notifications.GetNotificationsForUser(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting notifications")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, notificationsResults)
}

func MarkNotificationAsRead(ginCtx *gin.Context) {
	request := notifications.NotificationMarkAsReadData{}

	if err := ginCtx.ShouldBind(&request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "id should be a number"})
		return
	}

	err := notifications.MarkAsRead(request.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on marking notification as read")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func CreateNotifications(ginCtx *gin.Context) {
	request := notifications.NotificationRequest{}

	if err := ginCtx.ShouldBind(&request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": err.Error()})
		return
	}

	err := notifications.Create(request)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error(), "request": request}).
			Error("Error on creating notification")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}
