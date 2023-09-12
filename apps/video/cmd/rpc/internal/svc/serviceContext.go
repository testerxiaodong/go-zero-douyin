package svc

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/syncx"
	"go-zero-douyin/apps/video/cmd/rpc/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/internal/dao"
	"go-zero-douyin/common/cache"
	"go-zero-douyin/common/elasticService"
)

type ServiceContext struct {
	Config        config.Config
	VideoDo       dao.VideoDo
	SectionDo     dao.SectionDo
	TagDo         dao.TagDo
	Redis         cache.RedisCache
	Rabbit        rabbitmq.Sender
	Elasticsearch elasticService.ElasticService
	SingleFlight  syncx.SingleFlight
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		VideoDo:       dao.NewVideoRepository(c.DataSource),
		SectionDo:     dao.NewSectionRepository(c.DataSource),
		TagDo:         dao.NewTagRepository(c.DataSource),
		Redis:         cache.NewRedisClient(c.RedisCache),
		Rabbit:        rabbitmq.MustNewSender(c.RabbitSenderConf),
		Elasticsearch: elasticService.NewElastic(NewElasticsearchClient()),
		SingleFlight:  syncx.NewSingleFlight(),
	}
}

func NewElasticsearchClient() *elasticsearch.TypedClient {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{Addresses: []string{"http://127.0.0.1:9200"}})
	if err != nil {
		panic(err)
	}
	return client
}
