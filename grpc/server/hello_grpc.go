package main

import (
	"context"

	"github.com/jasonsoft/learning-opentelemetry/grpc/proto"

	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/log/v2"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Internal, "no metadata")
	}

	tracer := global.TracerProvider().Tracer("go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc")

	_, spanCtx := grpctrace.Extract(ctx, &md)

	// ctx = baggage.ContextWithMap(ctx, baggage.NewMap(baggage.MapUpdate{
	// 	MultiKV: entries,
	// }))

	if spanCtx.IsValid() {
		log.Debug("span is valid")
		ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
	} else {
		log.Debug("span is invalid")
	}

	// Start the server-side span, passing the remote
	// child span context explicitly.
	_, span := tracer.Start(
		ctx,
		"hello",
		//trace.WithAttributes(attrs...),
	)
	defer span.End()

	log.Debug("say hello is called")
	return &proto.HelloReply{Message: "Hello " + in.Name}, nil
}
