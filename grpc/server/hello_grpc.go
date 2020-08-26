package main

import (
	"context"

	"github.com/jasonsoft/learning-opentelemetry/grpc/proto"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/correlation"
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

	tracer := global.Tracer("")

	entries, spanCtx := grpctrace.Extract(ctx, &md)

	ctx = correlation.ContextWithMap(ctx, correlation.NewMap(correlation.MapUpdate{
		MultiKV: entries,
	}))

	if spanCtx.IsValid() {
		ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
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
