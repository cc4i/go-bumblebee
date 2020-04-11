package k8s

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		web := WebContext{}
		web.Handler(c.FullPath(), c)
	})
	//Realtime service: query air info with realtime update/push.
	r.GET("/spy", func(c *gin.Context) {
		web := WebContext{}
		web.Handler(c.FullPath(), c)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})

	return r
}
