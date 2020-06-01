package httpsvr

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var prefix = "/air"

func Router(ctx context.Context) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		Ping(ctx, c)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})

	r.GET(prefix+"/city/:city", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(prefix+"/ip/:ip", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(prefix+"/geo/:lat/:lng", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(prefix+"/aqi", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})

	return r
}
