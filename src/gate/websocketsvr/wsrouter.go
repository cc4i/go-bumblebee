package websocketsvr

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	r := gin.Default()

	// Home page for embedding  websocket page
	r.GET("/", func(c *gin.Context) {
		Home(c)
	})

	// Echo service: echo your message
	r.GET("/echo", func(c *gin.Context) {
		Echo(c)
	})

	//Realtime service: query air info with realtime update/push.
	r.GET("/spy", func(c *gin.Context) {
		Spy(c)
	})
	return r
}
