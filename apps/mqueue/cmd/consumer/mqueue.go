package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/service"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/listen"

	"github.com/zeromicro/go-zero/core/conf"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
)

var configFile = flag.String("f", "etc/mqueue.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	// logger、prometheus、trace、metricsUrl.
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range listen.Mqs(c) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
