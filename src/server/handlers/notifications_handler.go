package handlers

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/olahol/melody"
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

func RealtimeNotifications(websocket *melody.Melody) {
	websocket.HandleConnect(func(session *melody.Session) {
		err := session.Write([]byte("you are successfully connected to the recipes websocket."))
		if err != nil {
			return
		}
	})

	websocket.HandleMessage(func(session *melody.Session, message []byte) {
		request := notifications.NotificationRequest{}

		err := json.Unmarshal(message, &request)
		if err != nil {
			return
		}
		if _, err = validator.ValidateStruct(request); err != nil {
			errorMessage, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
			err = session.Write(errorMessage)
			if err != nil {
				utils.
					GetLogger().
					WithFields(log.Fields{"warning": err.Error(), "request": request}).
					Warn("Could not send message through the websocket")
			}
			return
		}

		receiversUsernames, err := notifications.Create(request)
		if err != nil {
			utils.
				GetLogger().
				WithFields(log.Fields{"error": err.Error(), "request": request}).
				Error("Error on creating notification")
			errorMessage, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
			err = session.Write(errorMessage)
			if err != nil {
				return
			}
			return
		}
		receiverIds, _ := json.Marshal(receiversUsernames)
		err = websocket.Broadcast(receiverIds)
		if err != nil {
			utils.
				GetLogger().
				WithFields(log.Fields{"warning": err.Error(), "request": request, "receivers": receiverIds}).
				Warn("Error on send receiver ids attempt")
			return
		}
	})
}
