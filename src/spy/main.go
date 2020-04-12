package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"spy/k8s"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	log.Fatal(k8s.Router().Run("0.0.0.0:9041"))
}
