package main

import (
	"air/aqi"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)


func router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/air/:city", func(c *gin.Context) {
		aqi.AirOfCity(c)
	})

	return r
}

func main() {
	log.Fatal(router().Run("0.0.0.0:9011"))
}
