package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"recipes-v2-server/server/handlers"
	"recipes-v2-server/server/middlewares"
	"recipes-v2-server/utils"
)

func setupRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(middlewares.Logger(utils.GetLogger()), gin.Recovery())
	router.Use(middlewares.CORS())

	router.GET("/healths", handlers.HealthCheck)
	router.GET("/metrics", handlers.Metrics)

	router.GET("/recipes", handlers.GetAllRecipes)
	router.GET("/recipes/category", handlers.GetByCategory)
	router.GET("/recipes/latest", handlers.GetLatestRecipes)
	router.GET("/recipes/most-popular", handlers.GetMostPopularRecipes)
	router.GET("/recipes/:name", handlers.GetRecipe)
	router.GET("/recipes/user/:username", handlers.GetRecipesByUser)
	router.GET("/recipes/favourites/:username", handlers.GetUserFavouriteRecipes)

	router.GET("/users/:username", handlers.GetUser)

	router.GET("/comments/latest", handlers.GetLatestComments)
	router.GET("/comments/:recipeName", handlers.GetRecipeComments)

	router.POST("/auth/login", handlers.Login)
	router.POST("/auth/check-username", handlers.CheckUsername)
	router.POST("/auth/check-email", handlers.CheckEmail)
	router.POST("/auth/register", handlers.Register)
	router.POST("/auth/generate-verification-code", handlers.GenerateVerificationCode)
	router.POST("/auth/verify-code", handlers.VerifyCode)
	router.POST("/auth/reset-password", handlers.ResetPassword)

	authGroup := router.Group("")
	authGroup.Use(middlewares.AuthMiddleware())
	{
		authGroup.GET("/notifications/:username", handlers.GetNotifications)
		authGroup.PUT("/notifications", handlers.MarkNotificationAsRead)
	}

	return
}

// Run defines the router endpoints and starts the server
func Run() {
	router := setupRouter()

	err := router.Run()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Unable to start web server")
	}

	utils.GetLogger().Debug("Web server started ...")
}
