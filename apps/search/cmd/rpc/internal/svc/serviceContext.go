package svc

import (
	"github.com/olivere/elastic/v7"
	"go-zero-douyin/apps/search/cmd/rpc/internal/config"
	"go-zero-douyin/common/elasticService"
)

type ServiceContext struct {
	Config        config.Config
	ElasticSearch elasticService.ElasticService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ElasticSearch: elasticService.NewElastic(NewElasticsearchClient(c.ElasticsearchConf)),
	}
}

func NewElasticsearchClient(config []string) *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(config...), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return client
}
