package main

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"

	"github.com/go-kratos/kratos/v2/config"
)

func main() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./config/log",
		CacheDir:            "./config/cache",
		LogLevel:            "debug",
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Panic(err)
	}

	c := config.New(
		config.WithSource(
			knacos.NewConfigSource(
				client,
				knacos.WithGroup("defaulttest_group"),
				knacos.WithDataID("defaulttest.yaml"),
			),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	name, err := c.Value("service.name").String()
	if err != nil {
		panic(err)
	}
	log.Println("GET service.name: ", name)

}
