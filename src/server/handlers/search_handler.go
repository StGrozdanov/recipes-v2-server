package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/search"
	"recipes-v2-server/utils"
)

var possibleSearchKeys = map[string]func(searchQuery string) (results pq.StringArray, err error){
	"users":    search.UserSearch,
	"comments": search.CommentsSearch,
	"recipes":  search.RecipesSearch,
	"global":   search.UserSearch,
}

func Search(ctx *gin.Context) {
	var (
		searchQuery string
		searchKey   string
	)

	for key := range possibleSearchKeys {
		exists := ctx.Request.URL.Query().Has(key)
		if exists {
			searchQuery = ctx.Request.URL.Query().Get(key)
			searchKey = key
			break
		}
	}

	if searchQuery != "" {
		if searchKey == "global" {
			globalResults, err := search.Global(searchQuery)
			if err != nil {
				utils.
					GetLogger().
					WithFields(log.Fields{"error": err.Error()}).
					Error("Error on performing a user search")

				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
				return
			}
			ctx.JSON(http.StatusOK, globalResults)
			return
		} else {
			results, err := possibleSearchKeys[searchKey](searchQuery)
			if err != nil {
				utils.
					GetLogger().
					WithFields(log.Fields{"error": err.Error()}).
					Error("Error on performing a user search")

				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
				return
			}
			ctx.JSON(http.StatusOK, map[string]interface{}{"content": results})
			return
		}
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"content": []string{}})
}
