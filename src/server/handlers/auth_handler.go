package handlers

import (
	"fmt"
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

func CheckUsername(ginCtx *gin.Context) {
	request := auth.UsernameData{}

	if err := ginCtx.ShouldBind(&request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid username"})
		return
	}

	if _, err := validator.ValidateStruct(request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid username"})
		utils.
			GetLogger().
			WithFields(log.Fields{"warning": err.Error()}).
			Warnf("Failed validation on authentication attempt for user %s", request.Username)
		return
	}

	isAvailable, err := auth.UsernameExists(request)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "error on username verification attempt"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on verification attempt for username %s", request.Username)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, isAvailable)
}

func CheckEmail(ginCtx *gin.Context) {
	request := auth.EmailData{}

	if err := ginCtx.ShouldBind(&request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid email"})
		return
	}

	if _, err := validator.ValidateStruct(request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid email"})
		utils.
			GetLogger().
			WithFields(log.Fields{"warning": err.Error()}).
			Warnf("Failed validation on authentication attempt for email %s", request.Email)
		return
	}

	isAvailable, err := auth.EmailExists(request)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "error on email verification attempt"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on verification attempt for email %s", request.Email)

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, isAvailable)
}

func Register(ginCtx *gin.Context) {
	request := auth.UserRegistrationData{}

	if err := ginCtx.ShouldBind(&request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(request); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": fmt.Sprintf("failed field validations - %s", err),
		})
		return
	}

	userData, err := auth.Register(request)
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
