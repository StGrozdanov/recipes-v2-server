package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/analytics"
	"recipes-v2-server/utils"
)

func GetVisitationsForTheLastSixMonths(ctx *gin.Context) {
	visitationsData, err := analytics.VisitationsForTheLastSixMonths()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusOK, map[string]interface{}{"error": "no data available"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting analytics for the last 6 months")

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, visitationsData)
}

func GetTheMostActiveUser(ctx *gin.Context) {
	user, err := analytics.MostActiveUser()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusOK, map[string]interface{}{"error": "no data available"})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting the most active user from the database")

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func GetVisitationsForTheDay(ctx *gin.Context) {
	visitationCount, err := analytics.VisitationsForToday()
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting visitations for today")

		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"count": visitationCount})
}
