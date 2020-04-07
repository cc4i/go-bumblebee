package main

import (
	"air/aqi"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/air/:city", func(c *gin.Context) {
		aqi.AirOfCity(c)
	})

	r.Run("0.0.0.0:9011")
}
