package main

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
)

func main() {
	consulClient, err := api.NewClient(api.DefaultConfig()) // consul client
	if err != nil {
		panic(err)
	}
	r := consul.New(consulClient) // 把 consulClient 客户端连接添加到 go-kratos 中的registry

	// grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterClient(conn) // gRPC 方式调用 helloworld 中方法

	// http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///helloworld"), // 服务名
		http.WithDiscovery(r),                        // 这里用 consul 作为服务发现中心
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hConn.Close()
	hClient := helloworld.NewGreeterClient(hConn) // http 方式调用 helloworld 中方法

	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		callHTTP(hClient)
	}

}

// grpc 方式 request
func callGRPC(client helloworld.GreetClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

// http 方式 request
func callHTTP(client helloworld.GreeterHTTPClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}
