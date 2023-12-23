package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"io"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetOutput(io.Discard)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
		},
	})
}

// GetLogger gets the configurated logrus logger
func GetLogger() *log.Entry {
	return log.WithFields(map[string]interface{}{
		"app_name": "recipes-v2-server",
	})
}
