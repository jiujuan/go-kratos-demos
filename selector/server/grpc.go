package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: fmt.Sprintf("welcome %+v!", in.Name)}, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.NewHelper(logger).Fatal(err)
	}
	go runServer("1.0.0", logger, consulClient, 8000)
	go runServer("1.0.0", logger, consulClient, 8010)

	runServer("2.0.0", logger, consulClient, 8020)
}

func runServer(version string, logger log.Logger, client *api.Client, port int) {
	logger = log.With(logger, "version", version, "port:", port)
	log := log.NewHelper(logger)

	grpcSrv := grpc.NewServer(
		grpc.Address(fmt.Sprintf(":%d", port+1000)),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)

	s := &server{}
	helloworld.RegisterGreeterServer(grpcSrv, s)

	r := consul.New(client)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Version(version),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
