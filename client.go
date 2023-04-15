package gofantasy

import (
	"net/http"
	"time"
)

type IClient interface {
	// WithOptions allows providing additional client options such as WithHTTPDebugging. These are not commonly needed.
	WithOptions(opts ...ClientOption) IClient

	Yahoo() IYahooClient
	ESPN() IEspnClient
}
type ICache interface {
}

type client struct {
	requestor *requestor
	cache     ICache
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
		baseUrl:   YahooBaseURL,
		requestor: c.requestor,
	}
}

func (c *client) ESPN() IEspnClient {
	return &espnClient{
		baseUrl: YahooBaseURL,
	}
}
