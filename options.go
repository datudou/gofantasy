package gofantasy

import (
	"net/http"
)

type ClientOption func(*client)

// WithHTTPClient allows providing a custom underlying http client.  It is good practice to _not_ use the default http client
// that Go provides as it has no timeouts.  If you do not provide your own default client, a reasonable one will be created for you.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) {
		c.requestor.httpClient = httpClient
	}
}

// WithHTTPDebugging will trigger logrus debug messages to be emitted with the raw request and response information.
// This should only be used for debugging purposes as it can decode entire messages into memory.
func WithHTTPDebugging(request bool, response bool) ClientOption {
	return func(c *client) {
		c.requestor.requestDebugging = request
		c.requestor.responseDebugging = response
	}
}

// WithCache provide lru cache for caching responses.  This is useful for reducing the number of requests made to the server side.
// The cache use golang-lru library, so the size is the number of items to store in the cache.
func WithCache(cache ICache) ClientOption {
	return func(c *client) {
		c.cache = cache
	}
}
