package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"recipes-v2-server/utils"
	"sync"
	"time"
)

type db struct {
	host     string
	user     string
	password string
	port     string
	database string
	DB       *sqlx.DB
}

var instance *db

// Init initialises database connection
func Init(hosts string, user string, password string, port string, database string) {
	var syncOnce sync.Once
	if instance == nil {
		syncOnce.Do(
			func() {
				instance = &db{
					host:     utils.GetRandomHost(hosts),
					user:     user,
					password: password,
					port:     port,
					database: database,
				}
				connect()
			},
		)
	}
}

func connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", instance.host, instance.port, instance.user, instance.password, instance.database)
	sqlxConnection, err := sqlx.Open("postgres", psqlInfo)
	if err == nil {
		instance.DB = sqlxConnection
		instance.DB.SetMaxOpenConns(15)
		instance.DB.SetMaxIdleConns(15)
		instance.DB.SetConnMaxIdleTime(5 * time.Minute)
		instance.DB.SetConnMaxLifetime(5 * time.Minute)
	} else {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on connection attempt to the Database")
	}
}
