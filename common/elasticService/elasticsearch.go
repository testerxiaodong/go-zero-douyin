package elasticService

import (
	"context"
	"github.com/olivere/elastic/v7"
	"go-zero-douyin/common/xconst"
)

type ElasticService interface {
	CreateDocument(ctx context.Context, indexName string, id string, req interface{}) (*elastic.IndexResponse, error)
	SearchByKeyword(ctx context.Context, indexName string, field string, keyword string, page int64, pageSize int64, sort string) (*elastic.SearchResult, error)
	DeleteDocument(ctx context.Context, indexName string, id string) (*elastic.DeleteResponse, error)
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

func (e *Elastic) SearchByKeyword(ctx context.Context, indexName string, field string, keyword string, page int64, pageSize int64, sort string) (*elastic.SearchResult, error) {
	matchQuery := elastic.NewMatchQuery(field, keyword)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = xconst.ElasticSearchVideoDefaultPageSize
	}
	return e.client.Search().Index(indexName).Query(matchQuery).Sort(sort, false).From(int((page - 1) * pageSize)).Size(int(pageSize)).Pretty(true).Do(ctx)
}

func (e *Elastic) DeleteDocument(ctx context.Context, indexName string, id string) (*elastic.DeleteResponse, error) {
	return e.client.Delete().Index(indexName).Id(id).Do(ctx)
}
