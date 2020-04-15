// Gate service
package main

import (
	"gate/graphqlsvr"
	"gate/httpsvr"
	"gate/tcpsvr"
	"gate/websocketsvr"
	log "github.com/sirupsen/logrus"
	"os"
)

// Run HTTP server to handle Restful API request.
func httpServer(endPoint string) {
	log.Printf("Serving on %s for Http 1/2 Service ...", endPoint)
	log.Fatal(httpsvr.Router().Run(endPoint))
}

// Run gRPC server
func gRpcServer(endPoint string) {
	log.Printf("Serving on %s for gRPC Service ...", endPoint)

}

// Run graphQL server
func graphQLServer(endPoint string) {
	log.Printf("Serving on %s for GraphQL Service ...", endPoint)
	log.Fatal(graphqlsvr.Router().Run(endPoint))
}

// Run websocket server
func websocketServer(endPoint string) {
	log.Printf("Serving on %s for WebSocket Service ...", endPoint)
	log.Fatal(websocketsvr.Router().Run(endPoint))

}

// Run raw TCP server
func tcpServer(endPoint string) {
	log.Printf("Serving on %s for raw TCP Service ...", endPoint)
	log.Fatal(tcpsvr.RunServer(endPoint))
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	//http
	go httpServer("0.0.0.0:9010")
	//gql
	go graphQLServer("0.0.0.0:9030")
	//ws
	go websocketServer("0.0.0.0:9040")
	//tcp
	tcpServer("0.0.0.0:9050")

}
