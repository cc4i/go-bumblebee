package aqi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

// API routes/path definition
func Router(ctx context.Context) *gin.Engine {
	r := gin.Default()

	// Get air quality data of given city
	// Example: ?c="city"
	r.GET("/air/city/:city", func(c *gin.Context) {
		AirOfCity(ctx, c)
	})

	// Get the nearest station close to the user location, based on the IP information.
	//
	// Destination: /feed/here/?token=:token
	//				api.ipstack.com/:ip?access_key=ad7c6834f8dba51e8943d96d3742fcc5
	r.GET("/air/ip/:ip", func(c *gin.Context) {
		AirOfIP(ctx, c)
	})

	// Get the nearest station close to the user location, based on the IP information.
	// Example: ?lat=&lng=
	//
	// Destination: /feed/geo::lat;:lng/?token=:token
	r.GET("/air/geo/:lat/:lng", func(c *gin.Context) {
		AirOfGeo(ctx, c)
	})

	// Get all the stations within a given lat/lng bounds
	// Example: ?lat1=39.379436&lng1=116.091230&lat2=40.235643&lng2=116.784382
	//
	// Destination: /map/bounds/?token=:token&latlng=:latlng
	r.POST("/station/bounds", func(c *gin.Context) {
		//TODO
		//lat1:lng1:lat2:lng2
		c.String(http.StatusHTTPVersionNotSupported, "Working in progress")
	})

	// Search stations by city name
	// Example: ?c="city"
	//
	// Destination: /search/?keyword=:keyword&token=:token
	r.POST("/station/city", func(c *gin.Context) {
		//TODO
		c.String(http.StatusHTTPVersionNotSupported, "Working in progress")
	})

	// Get AQIServer standard: Air Quality Index scale as defined by the US-EPA 2016 standard
	r.GET("/air/aqi", func(c *gin.Context) {
		c.String(http.StatusOK, `
			{
				"Standard": "Air Quality Index scale as defined by the US-EPA 2016 standard",
				"Definitions": [
					{
						"AQIServer": "0-50",
						"Level": "Good",
						"Implication": "Air quality is considered satisfactory, and air pollution poses little or no risk",
						"Caution": "None"
					},
					{
						"AQIServer": "51 -100",
						"Level": "Moderate",
						"Implication": "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution.",
						"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
					},
					{
						"AQIServer": "101-150",
						"Level": "Unhealthy for Sensitive Groups",
						"Implication": "Members of sensitive groups may experience health effects. The general public is not likely to be affected.",
						"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
					},
					{
						"AQIServer": "151-200",
						"Level": "Unhealthy",
						"Implication": "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects",
						"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid prolonged outdoor exertion; everyone else, especially children, should limit prolonged outdoor exertion"
					},
					{
						"AQIServer": "201-300",
						"Level": "Very Unhealthy",
						"Implication": "Health warnings of emergency conditions. The entire population is more likely to be affected.",
						"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid all outdoor exertion; everyone else, especially children, should limit outdoor exertion."
					},
					{
						"AQIServer": "300+",
						"Level": "Hazardous",
						"Implication": "Health alert: everyone may experience more serious health effects",
						"Caution": "Everyone should avoid all outdoor exertion"
					}
				]
			}	
		`)
	})

	// Basic health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// Metrics for Prometheus
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})
	// Version from ENV for test purpose
	r.GET("/version", func(c *gin.Context) {
		version := os.Getenv("AIR_VERSION")
		if version == "" {
			version = "v0.0.0"
		}
		c.String(http.StatusOK, version)
	})

	return r
}
