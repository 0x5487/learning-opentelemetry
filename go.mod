module github.com/0x5487/learning-opentelemetry

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/nite-coder/blackbear v0.0.0-20210710135651-97a27fc0a4df
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.21.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/net v0.7.0
	google.golang.org/grpc v1.39.0
)
