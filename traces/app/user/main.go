package main

import (
	"context"
	"os"
	"time"

	pbmsg "github.com/go-kratos/examples/traces/api/message"
	pb "github.com/go-kratos/examples/traces/api/user"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	Name    = "user"
	Version = "v1.0.0"
	Env     = "development"
)

type server struct {
	pb.UnimplementedUserServer
	tracer trace.TracerProvider
}

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// create the jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// always be sure to batch in production
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(Name),
				attribute.String("enviroment", Env),
				attribute.String("version", Version),
			)),
	)
	return tp, nil
}

// trace grpc client demo
func (s *server) GetMyMessages(ctx context.Context, in *pb.GetMyMessagesRequest) (*pb.GetMyMessagesReply, error) {
	// create conn grpc(client)
	// // only for demo, use single instance in production env
	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(middleware.Chain(
			tracing.Client( //trace client
				tracing.WithTracerProvider(s.tracer),
			),
			recovery.Recovery(),
		)),
		grpc.WithTimeout(time.Second*2),
	)
	if err != nil {
		return nil, err
	}
	msg := pbmsg.NewMessageServiceClient(conn)
	reply, err := msg.GetUserMessage(ctx, &pbmsg.GetUserMessageRequest{Id: 123, Count: in.Count})
	if err != nil {
		return nil, err
	}
	res := &pb.GetMyMessagesReply{}
	for _, v := range reply.Messages {
		res.Messages = append(res.Messages, &pb.Message{Content: v.Content})
	}
	return res, err
}

func main() {
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "trace_id", tracing.TraceID())
	logger = log.With(logger, "span_id", tracing.SpanID())
	log := log.NewHelper(logger)

	tp, err := tracerProvider("http://jaeger:14268/api/traces")
	if err != nil {
		log.Error(err)
	}

	httpSrv := http.NewServer(
		http.Address(":8080"),
		http.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				// Configuring tracing middleware
				tracing.Server(
					tracing.WithTracerProvider(tp), // 提供 trace provider
				),
				logging.Server(logger),
			),
		),
	)
	s := &server{}
	pb.RegisterUserHTTPServer(httpSrv, s)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
