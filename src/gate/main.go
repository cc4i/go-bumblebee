package main

import (
	"gate/graphqlsvr"
	"gate/httpsvr"
	"gate/tcpsvr"
	"gate/websocketsvr"
	log "github.com/sirupsen/logrus"
	"os"
)

func httpServer(endPoint string) {
	log.Printf("Serving on %s for Http 1/2 Service ...", endPoint)
	log.Fatal(httpsvr.Router().Run(endPoint))
}

func gRpcServer(endPoint string) {
	log.Printf("Serving on %s for gRPC Service ...", endPoint)

}

func graphQLServer(endPoint string) {
	log.Printf("Serving on %s for GraphQL Service ...", endPoint)
	log.Fatal(graphqlsvr.Router().Run(endPoint))
}

func websocketServer(endPoint string) {
	log.Printf("Serving on %s for WebSocket Service ...", endPoint)
	log.Fatal(websocketsvr.Router().Run(endPoint))

}

func tcpServer(endPoint string) {
	log.Printf("Serving on %s for raw TCP Service ...", endPoint)
	log.Fatal(tcpsvr.RunServer(endPoint))
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	go httpServer("0.0.0.0:9010")

	go graphQLServer("0.0.0.0:9030")
	go websocketServer("0.0.0.0:9040")
	tcpServer("0.0.0.0:9050")

}
