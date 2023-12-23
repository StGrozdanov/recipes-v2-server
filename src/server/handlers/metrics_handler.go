package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics endpoint to retrieve prometheus metrics
func Metrics(ginCtx *gin.Context) {
	promhttp.Handler().ServeHTTP(ginCtx.Writer, ginCtx.Request)
}
