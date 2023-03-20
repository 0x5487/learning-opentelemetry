module github.com/0x5487/learning-opentelemetry

go 1.14

require (
	github.com/golang/protobuf v1.5.3
	github.com/nite-coder/blackbear v0.0.0-20230316123859-b7d04f486c2c
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.40.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.40.0
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/exporters/jaeger v1.14.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.14.0
	go.opentelemetry.io/otel/sdk v1.14.0
	go.opentelemetry.io/otel/trace v1.14.0
	golang.org/x/net v0.8.0
	google.golang.org/grpc v1.53.0
)
