package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"net/http"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/instrumentation/httptrace"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "http-client",
			Tags: []kv.KeyValue{
				kv.String("version", "1.0"),
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
	err := tracer.WithSpan(ctx, "client http hello demo",
		func(ctx context.Context) error {
			req, _ := http.NewRequest("GET", "http://localhost:7777/hello", nil)

			//ctx, req = httptrace.W3C(ctx, req)
			httptrace.Inject(ctx, req)
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			_, err = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()
			span := trace.SpanFromContext(ctx)
			span.SetAttribute("request_id", "abc")
			span.SetStatus(codes.OK, "OK")
			return err
		})

	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
}
