package main

import (
	"air/aqi"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Fatal(aqi.Router().Run("0.0.0.0:9011"))
}
