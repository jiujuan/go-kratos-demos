package main

import (
	"context"
	"os"

	v1 "github.com/go-kratos/examples/traces/api/message"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"

	// "go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

var (
	Name    = "message"
	Env     = "env"
	Version = "v1.0.0"
)

type server struct {
	v1.UnimplementedMessageServiceServer
	tracer trace.TracerProvider
}

// get trace provider
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// create the jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	// New trace provider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// always be sure to batch in production
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(Name), // service name,实例名称
				attribute.String("env", Env),        // environment
				attribute.String("ID", Version),     // version
			)),
	)
	return tp, nil
}

func (s *server) GetUserMessage(ctx context.Context, request *v1.GetUserMessageRequest) (*v1.GetUserMessageReply, error) {
	msgs := &v1.GetUserMessageReply{}
	for i := 0; i < int(request.Count); i++ {
		msgs.Messages = append(msgs.Messages, &v1.Message{Content: "Teletubbies say hello."})
	}
	return msgs, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "trace_id", tracing.TraceID())
	logger = log.With(logger, "span_id", tracing.SpanID())
	log := log.NewHelper(logger)

	url := "http://jaeger:14268/api/traces"
	if os.Getenv("jaeger_url") != "" {
		url = os.Getenv("jeager_url")
	}

	tp, err := tracerProvider(url) // tracer provider
	if err != nil {
		log.Error(err)
	}

	s := &server{}

	// grpc server
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				tracing.Server(tracing.WithTracerProvider(tp)), //设置trace，传入 trace provider
				logging.Server(logger),
			),
		),
	)

	v1.RegisterMessageServiceServer(grpcSrv, s)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
