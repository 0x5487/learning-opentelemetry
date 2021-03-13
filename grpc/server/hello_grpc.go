package main

import (
	"context"
	"time"

	"github.com/jasonsoft/learning-opentelemetry/grpc/proto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/jasonsoft/log/v2"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Debug("== begin sayHello ==")
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, grpc.Errorf(codes.Internal, "no metadata")
	// }

	// tracer := otel.TracerProvider().Tracer("go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc")

	// _, spanCtx := grpctrace.Extract(ctx, &md)

	// // ctx = baggage.ContextWithMap(ctx, baggage.NewMap(baggage.MapUpdate{
	// // 	MultiKV: entries,
	// // }))

	// if spanCtx.IsValid() {
	// 	log.Debug("span is valid")
	// 	ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
	// } else {
	// 	log.Debug("span is invalid")
	// }

	time.Sleep(50 * time.Millisecond)

	span := trace.SpanFromContext(ctx)
	defer span.End()

	label2 := attribute.KeyValue{
		Key:   attribute.Key("key_aa"),
		Value: attribute.StringValue("value_aa"),
	}
	evt := trace.WithAttributes(label2)

	span.AddEvent("myEvent", evt)

	return &proto.HelloReply{Message: "Hello " + in.Name}, nil
}
