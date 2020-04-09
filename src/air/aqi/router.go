package aqi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/air/:city", func(c *gin.Context) {
		AirOfCity(c)
	})

	return r
}
