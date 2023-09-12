package elasticService

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type ElasticService interface {
	CreateDocument(ctx context.Context, indexName string, id string, req interface{}) (*elastic.IndexResponse, error)
	SearchByKeyword(ctx context.Context, indexName string, keyword string) (*elastic.SearchResult, error)
}

type Elastic struct {
	client *elastic.Client
}

func NewElastic(client *elastic.Client) ElasticService {
	return &Elastic{
		client: client,
	}
}

func (e *Elastic) CreateDocument(ctx context.Context, indexName string, id string, req interface{}) (*elastic.IndexResponse, error) {
	return e.client.Index().Index(indexName).Id(id).BodyJson(req).Do(ctx)
}

func (e *Elastic) SearchByKeyword(ctx context.Context, indexName string, keyword string) (*elastic.SearchResult, error) {
	matchQuery := elastic.NewMatchQuery("title", keyword)
	return e.client.Search().Index(indexName).Query(matchQuery).Sort("update_time", false).From(0).Size(10).Pretty(true).Do(ctx)
}
