module github.com/jasonsoft/learning-opentelemetry

go 1.15

require (
	github.com/golang/protobuf v1.4.2
	github.com/jasonsoft/log/v2 v2.0.0-beta.3
	go.opentelemetry.io/otel v0.10.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.10.0
	go.opentelemetry.io/otel/sdk v0.10.0
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd
	google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215 // indirect
	google.golang.org/grpc v1.31.0
)
