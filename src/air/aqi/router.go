package aqi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/air/:city", func(c *gin.Context) {
		AirOfCity(c)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})

	return r
}
