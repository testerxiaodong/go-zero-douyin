package http

import "github.com/go-resty/resty/v2"

type RestClient interface {
	Get(queryParams map[string]string, url string, response interface{}) error
	Post(body interface{}, url string, response interface{}) error
}

type RestyClient struct {
	c *resty.Client
}

func NewRestyClient() RestClient {
	c := resty.New()
	return &RestyClient{
		c: c,
	}
}

func (rc *RestyClient) Get(queryParams map[string]string, url string, response interface{}) error {
	_, err := rc.c.R().SetQueryParams(queryParams).
		SetHeader("Accept", "application/json").
		SetResult(response).Get(url)
	return err
}

func (rc *RestyClient) Post(body interface{}, url string, response interface{}) error {
	_, err := rc.c.R().
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		SetResult(response).
		Post(url)
	return err
}
