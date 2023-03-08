package main

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/examples/helloworld/helloworld"
	etcdregitry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	etcdclient "go.etcd.io/etcd/client/v3"
	srcgrpc "google.golang.org/grpc"
)

func main() {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	r := etcdregitry.New(client) // 传入 etcd client，也就是选择 etcd 为服务中心

	connGRPC, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"), // 服务发现
		grpc.WithDiscovery(r),                        // 传入etcd registry
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connGRPC.Close()

	connHTTP, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connHTTP.Close()

	for {
		callHTTP(connHTTP)
		callGRPC(connGRPC)
		time.Sleep(time.Second)
	}
}

func callHTTP(conn *http.Client) {
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "go-kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %+v\n", reply)
}

func callGRPC(conn *srcgrpc.ClientConn) {
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "go-kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}
