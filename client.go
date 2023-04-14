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

type client struct {
	yahooClient *yahooClient
	espnClient  *espnClient
	requestor   *requestor
}

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 30,
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

	c.yahooClient = &yahooClient{
		baseUrl:   YahooBaseURL,
		requestor: r,
	}
	return c
}

func (c *client) WithOptions(opts ...ClientOption) IClient {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *client) Yahoo() IYahooClient {
	return c.yahooClient
}

func (c *client) ESPN() IEspnClient {
	return c.espnClient
}
