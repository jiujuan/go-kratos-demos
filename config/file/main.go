package main

import (
	"log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

func main() {
	path := "./config.yaml"

	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	// 定义读取配置文件的结构
	var v struct {
		Service struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"service"`
	}

	if err := c.Scan(&v); err != nil {
		panic(err)
	}
	log.Printf("config: %+v", v)

	// 获取值
	name, err := c.Value("service.name").String()
	if err != nil {
		panic(err)
	}
	log.Printf("service: %s", name)

	// watch key
	if err := c.Watch("service.name", func(key string, value config.Value) {
		log.Printf("config changed: %s=%v\n", key, value)
	}); err != nil {
		panic(err)
	}

	<-make(chan struct{})
}
