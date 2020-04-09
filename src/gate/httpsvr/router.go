package httpsvr

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		Ping(c)
	})

	r.GET("/air/:city", func(c *gin.Context) {
		AirOfCity(c)
	})
	return r
}
