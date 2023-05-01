package gofantasy

import (
	"context"

	"github.com/gofantasy/model/yahoo"
)

type IClient interface {
	// WithOptions allows providing additional client options such as WithHTTPDebugging. These are not commonly needed.
	WithOptions(opts ...ClientOption) *yahooClient
	Get(ctx context.Context, endpoint string, objType string) (*yahoo.FantasyContent, error)
	SetCache(c ICache)
}

type client struct {
	Client
}

type Client struct {
	Requestor *requestor
	Cache     ICache
}
