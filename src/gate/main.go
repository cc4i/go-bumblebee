package main

import (
	"gate/httpsvr"
	"gate/tcpsvr"
	"gate/websocketsvr"
	log "github.com/sirupsen/logrus"
	"os"
)



func httpServer(endPoint string) {
	log.Fatal(httpsvr.Router().Run(endPoint))
}

func gRpcServer(endPoint string) {

}

func graphQLServer(endPoint string) {}

func websocketServer(endPoint string) {

	log.Fatal(websocketsvr.Router().Run(endPoint))

}

func tcpServer(endPoint string) {
	log.Fatal(tcpsvr.RunServer(endPoint))
}


func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	go httpServer("0.0.0.0:9010")
	go websocketServer("0.0.0.0:9040")
	tcpServer("0.0.0.0:9050")
}
