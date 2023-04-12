package main

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/v2/selector/filter"
	"github.com/go-kratos/v2/selector/wrr"
	"github.com/go-kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gitub.com/go-kratos/examples/helloworld/helloworld"
)

func main() {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	r := consul.New(consulCli)

	// grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r), // consul作为服务发现中心
		// 负载均衡 和 filter，weighted round robin算法
		grpc.WithBalancerName(wrr.Name),
		grpc.WithFilter(
			filter.Version("1.0.0"), //静态version=1.0.0的Filter
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterClient(conn)

	for {
		time.Sleep(time.Second)

		CallGRPC(gClient)
	}
}

func CallGRPC(client helloworld.NewGreeterClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "go-kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v \n", reply)
}
