package httpsvr

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(ctx context.Context) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		Ping(ctx, c)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})

	r.GET("/air/city/:city", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET("/air/ip/:ip", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET("/air/geo/:lat/:lng", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET("/air/aqi", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})

	return r
}
