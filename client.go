package gofantasy

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"net/http"
	"time"
)

type IClient interface {
	// WithOptions allows providing additional client options such as WithHTTPDebugging. These are not commonly needed.
	WithOptions(opts ...ClientOption) IClient
	Yahoo() IYahooClient
	//ESPN() IEspnClient
}

type client struct {
	requestor *requestor
	cache     *lru.Cache[interface{}, interface{}]
}

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

func NewClient(opts ...ClientOption) IClient {

	r := &requestor{
		httpClient: defaultHTTPClient,
	}
	c := &client{
		requestor: r,
	}
	c.WithOptions(opts...)

	return c
}

func (c *client) WithOptions(opts ...ClientOption) IClient {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *client) Yahoo() IYahooClient {
	return &yahooClient{
		baseUrl:     YahooBaseURL,
		baseClient:  c,
		yahooOAuth2: &yahooOAuth2{},
	}
}

//func (c *client) ESPN() IEspnClient {
//	return &espnClient{
//		baseUrl: YahooBaseURL,
//	}
//}
