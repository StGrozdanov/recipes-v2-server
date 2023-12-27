package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recipes-v2-server/internal/recipes"
	"recipes-v2-server/utils"
	"strconv"
)

func GetAllRecipes(ginCtx *gin.Context) {
	limitAsString := ginCtx.Request.URL.Query().Get("limit")
	cursorAsString := ginCtx.Request.URL.Query().Get("cursor")
	search := ginCtx.Request.URL.Query().Get("search")

	limit, limitErr := strconv.Atoi(limitAsString)
	cursor, cursorErr := strconv.Atoi(cursorAsString)
	if limitErr != nil && search == "" || cursorErr != nil && search == "" {
		ginCtx.JSON(http.StatusBadRequest, map[string]string{"error": "limit and cursor are required parameters and should be of type int. If they are not provided - search parameter should be provided instead."})
		return
	}

	if search != "" {
		recipesData, err := recipes.Search(search)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				ginCtx.JSON(http.StatusOK, map[string]interface{}{})
				return
			}

			utils.
				GetLogger().
				WithFields(log.Fields{"error": err.Error()}).
				Error("Error on search recipes from the database")

			ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
			return
		}
		ginCtx.JSON(http.StatusOK, recipesData)
		return
	}

	recipesData, err := recipes.GetAll(limit, cursor)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting latest recipes from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, recipesData)
}

func GetLatestRecipes(ginCtx *gin.Context) {
	latestRecipes, err := recipes.GetLatest()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting latest recipes from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, latestRecipes)
}

func GetMostPopularRecipes(ginCtx *gin.Context) {
	mostPopularRecipes, err := recipes.GetMostPopular()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting the most popular recipes from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, mostPopularRecipes)
}
