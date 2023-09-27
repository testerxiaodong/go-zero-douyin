package gorse

import "github.com/zhenghaoz/gorse/client"

type RecommendSystem interface {
}

type RecommendClient struct {
	Gorse *client.GorseClient
}

func NewRecommendClient(gorse *client.GorseClient) *RecommendClient {
	return &RecommendClient{
		Gorse: gorse,
	}
}

func (rc *RecommendClient) GetItemRecommendWithCategory() {

}
