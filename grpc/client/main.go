package main

import (
	"context"

	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	helloworldProto "github.com/jasonsoft/learning-opentelemetry/grpc/proto"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	address = "localhost:10051"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "grpc-client",
			Tags: []label.KeyValue{
				label.String("version", "1.0"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		panic(err)
	}

	return func() {
		flush()
	}
}

func main() {
	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	fn := initTracer()
	defer fn()

	tracer := global.Tracer("client-tracer")

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5,    // send pings every 5 seconds if there is no activity
			Timeout:             5,    // wait 5 second for ping ack before considering the connection dead
			PermitWithoutStream: true, // send pings even without active streams
		}),
		grpc.WithUnaryInterceptor(grpctrace.UnaryClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(tracer)),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := helloworldProto.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := "jason"
	r, err := c.SayHello(context.Background(), &helloworldProto.HelloRequest{Name: name})
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("Greeting: %s", r.Message)
}
