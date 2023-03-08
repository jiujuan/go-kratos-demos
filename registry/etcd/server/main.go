package main

import (
	"context"
	"fmt"
	"log"

	etcdregitry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	pb "github.com/go-kratos/examples/helloworld/helloworld"
	etcdclient "go.etcd.io/etcd/client/v3"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("welcome %+v!", in.Name)}, nil
}

func main() {
	// 创建 etcd client 连接
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化 http server
	httpSrv := http.NewServer(
		http.Address(":8080"),
		http.Middleware(
			recovery.Recovery(),
		),
	)

	// 初始化 grpc server
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	// 在服务器上注册服务
	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	// 创建一个 registry 对象，就是对 ectd client 操作的一个包装
	r := etcdregitry.New(client)

	app := kratos.New(
		kratos.Name("helloworld"), // 服务名称
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(r), // 填入etcd连接(etcd作为服务中心)
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
