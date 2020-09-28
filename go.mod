module github.com/jasonsoft/learning-opentelemetry

go 1.15

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/jasonsoft/log/v2 v2.0.0-beta.4
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.12.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.12.0 // indirect
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728
	google.golang.org/grpc v1.32.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
