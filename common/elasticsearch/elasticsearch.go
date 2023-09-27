package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
	"go-zero-douyin/common/xconst"
)

type ElasticService interface {
	CreateDocument(ctx context.Context, indexName string, id string, req interface{}) (*elastic.IndexResponse, error)
	SearchByKeyword(ctx context.Context, indexName string, field string, keyword string, page int64, pageSize int64, sort string, highLight int64) (*elastic.SearchResult, error)
	DeleteDocument(ctx context.Context, indexName string, id string) (*elastic.DeleteResponse, error)
	Suggestion(ctx context.Context, indexName string, suggestionName string, input string) (*elastic.SearchResult, error)
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

func (e *Elastic) SearchByKeyword(ctx context.Context, indexName string, field string, keyword string, page int64, pageSize int64, sort string, highLight int64) (*elastic.SearchResult, error) {
	matchQuery := elastic.NewMatchQuery(field, keyword)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = xconst.ElasticSearchVideoDefaultPageSize
	}
	if highLight == 1 {
		hl := elastic.NewHighlight()
		hl = hl.Fields(elastic.NewHighlighterField(field).RequireFieldMatch(false))
		hl = hl.PreTags("<em>").PostTags("</em>")
		return e.client.Search().Index(indexName).Highlight(hl).Query(matchQuery).Sort(sort, false).From(int((page - 1) * pageSize)).Size(int(pageSize)).Pretty(true).Do(ctx)
	}
	return e.client.Search().Index(indexName).Query(matchQuery).Sort(sort, false).From(int((page - 1) * pageSize)).Size(int(pageSize)).Pretty(true).Do(ctx)
}

func (e *Elastic) DeleteDocument(ctx context.Context, indexName string, id string) (*elastic.DeleteResponse, error) {
	return e.client.Delete().Index(indexName).Id(id).Do(ctx)
}

func (e *Elastic) Suggestion(ctx context.Context, indexName string, suggestionName string, input string) (*elastic.SearchResult, error) {
	suggester := elastic.NewCompletionSuggester(suggestionName).
		Field("suggestion").
		Prefix(input).
		Size(10) // 限制返回的建议结果数量

	searchService := e.client.Search().
		Index(indexName).
		Suggester(suggester)

	return searchService.Do(ctx)
}
