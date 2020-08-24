package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/instrumentation/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "http-server",
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

func helloHandler(w http.ResponseWriter, req *http.Request) {
	tracer := global.Tracer("")

	// Extracts the conventional HTTP span attributes,
	// distributed context tags, and a span context for
	// tracing this request.
	attrs, entries, spanCtx := httptrace.Extract(req.Context(), req)
	ctx := req.Context()
	if spanCtx.IsValid() {
		ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
	}

	// Apply the correlation context tags to the request
	// context.
	req = req.WithContext(correlation.ContextWithMap(ctx, correlation.NewMap(correlation.MapUpdate{
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

	span.AddEvent(ctx, "handling this...", kv.Int("request-handled", 100))

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
