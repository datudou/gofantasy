package gofantasy

import (
	"bytes"
	"context"
	"net/http"
)

type requestBuilder interface {
	build(ctx context.Context, method, url string, request any) (*http.Request, error)
}

type httpRequestBuilder struct {
	marshaller marshaller
}

func newRequestBuilder() *httpRequestBuilder {
	return &httpRequestBuilder{
		marshaller: &xmlMarshaller{},
	}
}

func (h *httpRequestBuilder) build(ctx context.Context, method, url string, request any) (*http.Request, error) {
	if request == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	var reqBytes []byte
	reqBytes, err := h.marshaller.marshal(request)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewBuffer(reqBytes),
	)
}
