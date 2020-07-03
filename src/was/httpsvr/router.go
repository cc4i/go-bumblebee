package httpsvr

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var airV1 = "/air/v1"
var airV2 = "/air/v2"

func Router(ctx context.Context) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		Ping(ctx, c)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})

	r.GET(airV1+"/city/:city", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV1+"/ip/:ip", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV1+"/geo/:lat/:lng", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV1+"/aqi", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV1+"/version", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})

	r.GET(airV2+"/city/:city", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV2+"/ip/:ip", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV2+"/geo/:lat/:lng", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV2+"/aqi", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})
	r.GET(airV2+"/version", func(c *gin.Context) {
		ProxyAir(ctx, c)
	})

	return r
}
