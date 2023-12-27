package main

import (
	log "github.com/sirupsen/logrus"
	"recipes-v2-server/config"
	"recipes-v2-server/database"
	"recipes-v2-server/internal/auth"
	"recipes-v2-server/server"
	"recipes-v2-server/utils"
)

func init() {
	app, err := config.Init()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on config initialization")
		return
	}
	if app.AppEnv == "LOC" {
		utils.PrettyPrint(app)
	}

	database.Init(
		app.DBHosts,
		app.DBUsername,
		app.DBPassword,
		app.DBPort,
		app.DBName,
	)

	utils.CreateS3Session(
		app.S3BucketName,
		app.S3BucketKey,
		app.S3BucketURL,
		app.S3BucketRegion,
		app.AWSAccessKey,
		app.AWSSecretKey,
		app.S3ACL,
	)

	utils.GetJWTKey(app.JWTSecret)

	auth.GetSaltRounds(app.Salt)
}

func main() {
	server.Run()
}
