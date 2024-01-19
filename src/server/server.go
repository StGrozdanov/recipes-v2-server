package server

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"recipes-v2-server/database"
	"recipes-v2-server/server/handlers"
	"recipes-v2-server/server/middlewares"
	"recipes-v2-server/utils"
)

func setupRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	websocket := melody.New()

	router.Use(middlewares.Logger(utils.GetLogger()), gin.Recovery())
	router.Use(middlewares.CORS())
	router.Use(middlewares.TrackVisitations())
	router.Use(middlewares.FilterBlockedUsers())

	router.GET("/healths", handlers.HealthCheck)
	router.GET("/metrics", handlers.Metrics)

	router.GET("/recipes", handlers.GetAllRecipes)
	router.GET("/recipes/category", handlers.GetByCategory)
	router.GET("/recipes/latest", handlers.GetLatestRecipes)
	router.GET("/recipes/most-popular", handlers.GetMostPopularRecipes)
	router.GET("/recipes/:name", handlers.GetRecipe)
	router.GET("/recipes/user/:username", handlers.GetRecipesByUser)
	router.GET("/recipes/favourites/:username", handlers.GetUserFavouriteRecipes)
	router.POST("/recipes/is-favourite", handlers.CheckIfRecipeIsInFavourites)
	router.POST("/recipes/check-name", handlers.CheckRecipeName)

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

	router.GET("/realtime-notifications", func(ctx *gin.Context) {
		if err := websocket.HandleRequest(ctx.Writer, ctx.Request); err != nil {
			return
		}
		handlers.RealtimeNotifications(websocket)
	})

	authGroup := router.Group("")
	authGroup.Use(middlewares.AuthMiddleware())
	{
		authGroup.GET("/notifications/:username", handlers.GetNotifications)
		authGroup.PUT("/notifications", handlers.MarkNotificationAsRead)

		authGroup.POST("/recipes/add-to-favourites", handlers.AddToFavourites)
		authGroup.DELETE("/recipes/remove-from-favourites", handlers.RemoveFromFavourites)
		authGroup.POST("/recipes", handlers.CreateRecipe)
		authGroup.POST("/recipes/upload-image", handlers.UploadRecipeImage)

		authGroup.POST("/comments", handlers.CreateComment)

		imageUploadGroup := authGroup.Group("/upload/image/users")
		imageUploadGroup.Use(middlewares.ImageContentTypeMiddleware())
		{
			imageUploadGroup.POST("/cover-image", handlers.UploadCoverImage)
			imageUploadGroup.POST("/avatar-image", handlers.UploadAvatarImage)
		}
	}

	resourceOwnerGroup := router.Group("")
	resourceOwnerGroup.Use(middlewares.ResourceOwnerMiddleware())
	{
		resourceOwnerGroup.PATCH("/users/:username", handlers.EditUserData)
		resourceOwnerGroup.PUT("/recipes/:name", handlers.EditRecipe)
		resourceOwnerGroup.DELETE("/recipes/:name", handlers.DeleteRecipe)
		resourceOwnerGroup.PUT("/comments", handlers.EditComment)
		resourceOwnerGroup.DELETE("/comments", handlers.DeleteComment)
	}

	adminGroup := router.Group("/admin")
	adminGroup.Use(middlewares.AdminMiddleware())
	{
		adminGroup.GET("/recipes/count", handlers.GetRecipesCount)
		adminGroup.GET("/comments/count", handlers.GetCommentsCount)
		adminGroup.GET("/users/count", handlers.GetUsersCount)
		adminGroup.GET("/analytics/visitations", handlers.GetVisitationsForTheLastSixMonths)
		adminGroup.GET("/analytics/most-active-user", handlers.GetTheMostActiveUser)
		adminGroup.GET("/analytics/visitations/today", handlers.GetVisitationsForTheDay)
		adminGroup.GET("/search", handlers.Search)
		adminGroup.GET("/users", handlers.GetAllUsers)
		adminGroup.GET("/recipes", handlers.GetAllRecipesAdmin)
		adminGroup.GET("/comments", handlers.GetAllComments)
	}

	return
}

// Run defines the router endpoints and starts the server
func Run() {
	router := setupRouter()
	cronjob := cron.New()

	_, err := cronjob.AddFunc("0 2 * * 1", cleanUpResetPasswordRequests)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error adding clean up password requests job")
	}
	_, err = cronjob.AddFunc("0 0 * * *", cleanUpNotificationsMarkedAsRead)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error adding clean up notifications job")
	}

	err = router.Run()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Unable to start web server")
	}

	utils.GetLogger().Debug("Web server started ...")
}

func cleanUpResetPasswordRequests() {
	_, err := database.ExecuteQuery(
		`DELETE 
				FROM password_requests 
       			WHERE issued_at < CURRENT_TIMESTAMP - INTERVAL '1 day'`,
	)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error executing clean up password requests job")
	}
}

func cleanUpNotificationsMarkedAsRead() {
	_, err := database.ExecuteQuery(
		`DELETE FROM notifications WHERE is_marked_as_read IS TRUE;`,
	)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error executing clean up notifications requests job")
	}
}
