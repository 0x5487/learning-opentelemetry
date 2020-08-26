package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"net/http"

	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace"
	"go.opentelemetry.io/otel/api/global"

	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "http-client",
			Tags: []label.KeyValue{
				label.String("version", "1.0"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		flush()
	}
}

func main() {
	fn := initTracer()
	defer fn()

	client := http.DefaultClient
	// ctx := correlation.NewContext(context.Background(),
	// 	kv.String("username", "donuts"),
	// )

	ctx := context.Background()
	tracer := global.Tracer("http-client")

	ctx, span := tracer.Start(ctx, "client http hello demo")
	defer span.End()

	req, _ := http.NewRequest("GET", "http://localhost:7777/hello", nil)
	httptrace.Inject(ctx, req)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	_, err = ioutil.ReadAll(res.Body)
	_ = res.Body.Close()

	span.SetAttribute("request_id", "abc")
	span.SetStatus(codes.OK, "OK")

	time.Sleep(time.Second * 2)
}
