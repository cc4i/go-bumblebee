package k8s

import "github.com/gin-gonic/gin"

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
	return r
}
