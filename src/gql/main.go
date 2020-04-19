package main

import (
	log "github.com/sirupsen/logrus"
	"gql/graph"
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

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
		log.WithFields(log.Fields{
			"Jaeger": jaeger.JaegerClientVersion,
		}).Error("cannot init Jaeger: ", err)
		panic(err)
	}
	return tracer, closer
}

func main() {

	tracer, closer := initJaeger("gql-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	log.Fatal(graph.Router().Run("0.0.0.0:9030"))
}
