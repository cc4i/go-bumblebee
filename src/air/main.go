package main

import (
	"air/apis"
	"context"
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"io"
)

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
	tracer, closer := initJaeger("air-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span := tracer.StartSpan("http")
	span.SetTag("air", "http-9011")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	log.Fatal(apis.Router(ctx).Run("0.0.0.0:9011"))
	defer span.Finish()

}
