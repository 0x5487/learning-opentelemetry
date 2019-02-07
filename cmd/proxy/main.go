package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"net/url"

	"github.com/jasonsoft/request"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func main() {

	tracer, closer := Init("bifrost") // tracer app name
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	helloTo := "google"

	span := tracer.StartSpan("begin")
	span.SetTag("hello-to", helloTo)
	span.SetTag("request-id", "abcd")
	span.SetBaggageItem("mybaggage", "abcd")

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)

	// simulate that send a message to mq
	ext.SpanKindProducer.Set(span)
	textMap := map[string]string{}
	span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(textMap),
	)
	span.Finish()
	sendToMQ(textMap)
}

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
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

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://localhost:8081/format?" + v.Encode()

	// set tags
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	tempHeader := map[string][]string{}
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(tempHeader),
	)

	for k, v := range tempHeader {
		fmt.Printf("tempHeader: %s,%v\n", k, v)
	}

	traceInfo := tempHeader["Uber-Trace-Id"][0]

	resp, err := request.
		GET(url).
		Set("Uber-Trace-Id", traceInfo).
		End()

	if err != nil {
		panic(err.Error())
	}

	helloStr := resp.String()

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}

func sendToMQ(textMapCarrier opentracing.TextMapCarrier) {
	time.Sleep(1 * time.Second)
	tracer, closer := Init("mq") // tracer app name
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	spanCtx, _ := tracer.Extract(opentracing.TextMap, textMapCarrier)
	span := tracer.StartSpan("mq-start", ext.RPCServerOption(spanCtx))
	ext.SpanKindConsumer.Set(span)
	defer span.Finish()

	mqStr := "Hello MQ"
	time.Sleep(1 * time.Second)
	println(mqStr)
	span.LogKV("event", "send-to-mq")
}
