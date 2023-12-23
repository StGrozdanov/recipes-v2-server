package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"os"
	"time"
)

var (
	timeFormat = "02/Jan/2006:15:04:05 -0700"
	skip       map[string]struct{}
)

// Logger logging middleware that tracks the request and response data and execution times. Accepts fields
// that can be ignored.
func Logger(logger logrus.FieldLogger, notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	populateSkippedFields(notLogged...)

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		entry := logger.WithFields(logrus.Fields{
			"hostname":       hostname,
			"statusCode":     statusCode,
			"responseTimeMs": latency,
			"clientIP":       clientIP,
			"method":         c.Request.Method,
			"path":           path,
			"referer":        referer,
			"dataLength":     dataLength,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, dataLength, referer, latency)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}

func populateSkippedFields(notLogged ...string) {
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}
}
