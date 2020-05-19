// Gate service
package main

import (
	"context"
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

	log.Info("Serving on for Http 1/2 Service ...", endPoint)
	log.Fatal(httpsvr.Router(ctx).Run(endPoint))
}

// Run gRPC server
func gRpcServer(endPoint string) {
	log.Info("Serving on for gRPC Service ...", endPoint)

}

// Run websocket server
func websocketServer(endPoint string) {
	log.Info("Serving on for WebSocket Service ...", endPoint)
	log.Fatal(websocketsvr.Router().Run(endPoint))

}

// Run raw TCP server
func tcpServer(endPoint string) {
	log.Info("Serving on for raw TCP Service ...", endPoint)
	log.Fatal(tcpsvr.RunServer(endPoint))
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

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
