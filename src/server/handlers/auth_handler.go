package handlers

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/auth"
	"recipes-v2-server/utils"
	"strings"
)

func Login(ginCtx *gin.Context) {
	request := auth.UserAuthData{}

	if err := ginCtx.ShouldBind(&request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid username or password"})
		utils.
			GetLogger().
			WithFields(log.Fields{"warning": err.Error()}).
			Warnf("Failed validation on authentication attempt for user %s", request.Username)
		return
	}

	userData, err := auth.Login(request)
	if err != nil {
		if err.Error() == "sql: no rows in result set" || strings.Contains(err.Error(), "crypto/bcrypt") {
			ginCtx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "invalid username or password"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on authentication attempt for user %s", request.Username)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, userData)
}
