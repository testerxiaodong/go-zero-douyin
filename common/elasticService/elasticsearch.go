package elasticService

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
)

type ElasticService interface {
	CreateDocument(ctx context.Context, indexName string, id string, req interface{}) (*index.Response, error)
}

type Elastic struct {
	client *elasticsearch.TypedClient
}

func NewElastic(client *elasticsearch.TypedClient) ElasticService {
	return &Elastic{
		client: client,
	}
}

func (e *Elastic) CreateDocument(ctx context.Context, indexName string, id string, req interface{}) (*index.Response, error) {
	return e.client.Index(indexName).Id(id).Request(req).Do(ctx)
}

func (e *Elastic) SearchByKeyword(ctx context.Context, indexName string, keyword string) {
}
