package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/nleeper/goment"
	log "github.com/sirupsen/logrus"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
	"strings"
	"time"
)

// TrackVisitations tracks every unique website visitation and adds the authenticated users IP to the database
// IP list
func TrackVisitations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()

		referer := ctx.Request.Referer()

		if referer != "https://all-the-best-recipes.vercel.app/" {
			return
		}

		dateNow, err := goment.New(time.Now())
		if err != nil {
			return
		}
		date := dateNow.UTC().Format("YYYY-MM-DD HH:mm:ss")

		collectVisitationToVisitationsTable(date, clientIP)

		if len(ctx.Request.Header["X-Authorization"]) != 0 && len(ctx.Request.Header["X-Authorization"][0]) != 0 {
			token := ctx.Request.Header["X-Authorization"][0]
			collectIPtoAuthenticatedUsersIPsList(clientIP, token)
		}

		isRecipeEndpoint := strings.Contains(ctx.Request.URL.String(), "/recipes")
		isGETRequest := ctx.Request.Method == "GET"
		recipeName, recipeNameParamExists := ctx.Params.Get("name")

		if isRecipeEndpoint && isGETRequest && recipeNameParamExists {
			addANewRecipeVisitation(recipeName)
		}
		ctx.Next()
	}
}

func addANewRecipeVisitation(recipeName string) {
	_, err := database.ExecuteNamedQuery(
		`UPDATE recipes
				SET visitations_count = visitations_count + 1
				WHERE recipe_name = :recipe_name;`,
		map[string]interface{}{"recipe_name": recipeName},
	)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on track recipe visitation attempt for recipe %s", recipeName)
		return
	}
}

func collectIPtoAuthenticatedUsersIPsList(ip string, token string) {
	claims, isValid, err := utils.ParseJWT(token)
	if err != nil || !isValid {
		return
	}

	userId := claims.Id

	_, err = database.ExecuteNamedQuery(
		`WITH same_ip_data AS (SELECT EXISTS (SELECT user_entity_id
                                     FROM user_entity_ip_addresses
                                     WHERE ip_addresses = :ip_address
                                       AND user_entity_id = :user_id) AS exists)
				INSERT
				INTO user_entity_ip_addresses (ip_addresses, user_entity_id)
				SELECT :ip_address, :user_id
				WHERE (SELECT exists FROM same_ip_data) = false;`,
		map[string]interface{}{"ip_address": ip, "user_id": userId},
	)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on track IP address attempt for user %s", claims.Username)
		return
	}
}

func collectVisitationToVisitationsTable(date, clientIP string) {
	var exists bool

	err := database.GetSingleRecordNamedQuery(
		&exists,
		`SELECT EXISTS (SELECT id
									FROM visitations
									WHERE ip_address = :client_ip
									  AND CAST(DATE_TRUNC('day', visited_at) AS DATE) = :date
									);`,
		map[string]interface{}{"date": date, "client_ip": clientIP},
	)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Errorf("Error on track visitation attempt for IP %s", clientIP)
		return
	}

	if !exists {
		_, err = database.ExecuteNamedQuery(
			`INSERT INTO visitations (ip_address, visited_at)
				VALUES (:ip_address, :date)`,
			map[string]interface{}{
				"date":       date,
				"ip_address": clientIP,
			})
		if err != nil {
			utils.
				GetLogger().
				WithFields(log.Fields{"error": err.Error()}).
				Errorf("Error on track visitation attempt for IP %s", clientIP)
			return
		}
	}
}
