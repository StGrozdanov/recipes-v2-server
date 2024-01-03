package handlers

import (
	validator "github.com/asaskevich/govalidator"
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

func GetByCategory(ginCtx *gin.Context) {
	query := ginCtx.Request.URL.Query().Get("name")

	if query == "" {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "category name should not be empty"})
		return
	}

	recipesData, err := recipes.SearchByCategory(query)
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
}

func GetRecipe(ginCtx *gin.Context) {
	recipeName, ok := ginCtx.Params.Get("name")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "recipe name was not found"})
		return
	}

	recipe, err := recipes.GetASingleRecipe(recipeName)
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
	ginCtx.JSON(http.StatusOK, recipe)
}

func GetRecipesByUser(ginCtx *gin.Context) {
	username, ok := ginCtx.Params.Get("username")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "owner username was not found"})
		return
	}

	recipesResults, err := recipes.GetRecipesFromUser(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting the recipes from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, recipesResults)
}

func GetUserFavouriteRecipes(ginCtx *gin.Context) {
	username, ok := ginCtx.Params.Get("username")

	if !ok {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"errors": "username was not found"})
		return
	}

	recipesResults, err := recipes.GetFavourites(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting the recipes from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, recipesResults)
}

func CheckIfRecipeIsInFavourites(ginCtx *gin.Context) {
	data := recipes.FavouritesRequest{}

	if err := ginCtx.ShouldBind(&data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	isInFavourites, err := recipes.IsInFavourites(data)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	ginCtx.JSON(http.StatusOK, isInFavourites)
}

func AddToFavourites(ginCtx *gin.Context) {
	data := recipes.FavouritesRequest{}

	if err := ginCtx.ShouldBind(&data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	err := recipes.AddToFavourites(data)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

func RemoveFromFavourites(ginCtx *gin.Context) {
	data := recipes.FavouritesRequest{}

	if err := ginCtx.ShouldBind(&data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(data); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	err := recipes.RemoveFromFavourites(data)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"success": true})
}
