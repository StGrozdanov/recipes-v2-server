package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
)

// FilterBlockedUsers filters out blacklisted users
func FilterBlockedUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()

		var blockedFor string

		err := database.GetSingleRecordNamedQuery(
			&blockedFor,
			`SELECT CASE
				   WHEN EXISTS (SELECT ip_address FROM blacklist WHERE ip_address = :ip_address)
					   THEN (SELECT reason FROM blacklist WHERE ip_address = :ip_address)
				   ELSE ''
				   END AS reason;`,
			map[string]interface{}{"ip_address": clientIP},
		)
		if err != nil {
			utils.
				GetLogger().
				WithFields(log.Fields{"error": err.Error()}).
				Errorf("Error on checking if user is blocked for client ip - %s", clientIP)
		}

		if blockedFor != "" {
			ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, Errors{
				Info: Info{
					Message: blockedFor,
					Cause:   "Blocked",
				},
			})
			return
		}

		ctx.Next()
	}
}
