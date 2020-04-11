package main

import (
	log "github.com/sirupsen/logrus"
	"spy/k8s"
)

func main() {

	log.Fatal(k8s.Router().Run("0.0.0.0:9041"))
}

