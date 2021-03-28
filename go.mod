module github.com/jasonsoft/learning-opentelemetry

go 1.15

require (
	github.com/golang/protobuf v1.5.1
	github.com/jasonsoft/log/v2 v2.0.0-beta.4
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.19.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.19.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.19.0
	go.opentelemetry.io/otel/sdk v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	golang.org/x/net v0.0.0-20210326220855-61e056675ecf
	google.golang.org/grpc v1.36.1
)
