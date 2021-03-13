package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "http-server",
			Tags: []attribute.KeyValue{
				attribute.String("version", "1.0"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return func() {
		flush()
	}
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello was called")

	uk := attribute.Key("username")
	ctx := req.Context()
	span := trace.SpanFromContext(ctx)
	username := baggage.Value(ctx, uk)
	span.AddEvent("handling this...", trace.WithAttributes(uk.String(username.AsString())))

	time.Sleep(1 * time.Second)

	_, _ = io.WriteString(w, "Hello, world!\n")

}

func main() {
	fn := initTracer()
	defer fn()

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "HelloOperation")
	http.Handle("/hello", otelHandler)

	err := http.ListenAndServe("localhost:7777", nil)
	if err != nil {
		panic(err)
	}
}
