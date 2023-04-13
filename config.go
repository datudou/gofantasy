package gofantasy

import (
	"net/http"
)

const (
	yahooFantasyV2Uri              = "https://fantasysports.yahooapis.com/fantasy/v2"
	defaultEmptyMessagesLimit uint = 300
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken          string
	BaseURL            string
	HTTPClient         *http.Client
	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken:          authToken,
		BaseURL:            yahooFantasyV2Uri,
		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
