package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/api/baggage"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
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
			ServiceName: "http-server",
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

func helloHandler(w http.ResponseWriter, req *http.Request) {
	tracer := global.Tracer("")

	// Extracts the conventional HTTP span attributes,
	// distributed context tags, and a span context for
	// tracing this request.
	attrs, entries, spanCtx := otelhttptrace.Extract(req.Context(), req)
	ctx := req.Context()
	if spanCtx.IsValid() {
		ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
	}

	// Apply the correlation context tags to the request
	// context.
	req = req.WithContext(baggage.ContextWithMap(ctx, baggage.NewMap(baggage.MapUpdate{
		MultiKV: entries,
	})))

	// Start the server-side span, passing the remote
	// child span context explicitly.
	_, span := tracer.Start(
		req.Context(),
		"hello",
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	time.Sleep(1 * time.Second)

	span.AddEvent(ctx, "handling this...", label.Int("request-handled", 100))

	_, _ = io.WriteString(w, "Hello, world!\n")
	fmt.Println("hello is called")
}

func main() {
	fn := initTracer()
	defer fn()

	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe("localhost:7777", nil)
	if err != nil {
		panic(err)
	}
}
