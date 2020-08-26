package main

import (
	"net"
	"time"

	helloworldProto "github.com/jasonsoft/learning-opentelemetry/grpc/proto"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	port = "127.0.0.1:10051"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "grpc-server",
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

// server is used to implement helloworld.GreeterServer.
func main() {
	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	fn := initTracer()
	defer fn()

	tracer := global.Tracer("server-tracer")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    (time.Duration(5) * time.Second), // Ping the client if it is idle for 5 seconds to ensure the connection is still active
				Timeout: (time.Duration(5) * time.Second), // Wait 5 second for the ping ack before assuming the connection is dead
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             (time.Duration(2) * time.Second), // If a client pings more than once every 2 seconds, terminate the connection
				PermitWithoutStream: true,                             // Allow pings even when there are no active streams
			},
		),
		grpc.UnaryInterceptor(grpctrace.UnaryServerInterceptor(tracer)),
		grpc.StreamInterceptor(grpctrace.StreamServerInterceptor(tracer)),
	)

	server := NewServer()
	helloworldProto.RegisterGreeterServer(s, server)

	log.Debug("server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
