// Gate service
package main

import (
	"context"
	"gate/graphqlsvr"
	"gate/httpsvr"
	"gate/tcpsvr"
	"gate/websocketsvr"
	log "github.com/sirupsen/logrus"
	"io"
	"os"

	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

// Run HTTP server to handle Restful API request.
func httpServer(ctx context.Context, endPoint string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "httpServer")
	defer span.Finish()

	log.Printf("Serving on %s for Http 1/2 Service ...", endPoint)
	log.Fatal(httpsvr.Router(ctx).Run(endPoint))
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

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {

	tracer, closer := initJaeger("gate-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	//http
	span := tracer.StartSpan("http")
	span.SetTag("http-to", "http-9010")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	go httpServer(ctx, "0.0.0.0:9010")
	defer span.Finish()

	//gql
	span2 := tracer.StartSpan("gql")
	span2.SetTag("gql-to", "gql-9040")
	go graphQLServer("0.0.0.0:9030")
	defer span2.Finish()

	//ws
	span3 := tracer.StartSpan("ws")
	span3.SetTag("ws-to", "ws-9040")
	go websocketServer("0.0.0.0:9040")
	defer span3.Finish()

	//tcp
	span4 := tracer.StartSpan("tcp")
	span4.SetTag("tcp-to", "ws-9050")
	tcpServer("0.0.0.0:9050")
	defer span4.Finish()

}
