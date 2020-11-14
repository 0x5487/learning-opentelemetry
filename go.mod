module github.com/jasonsoft/learning-opentelemetry

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/jasonsoft/log/v2 v2.0.0-beta.4
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.13.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	google.golang.org/grpc v1.33.2
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
)
