package gofantasy

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type requestDecorator func(*http.Request) *http.Request

type requestor struct {
	baseURL                string
	AuthorizationDecorator requestDecorator
	requestDebugging       bool
	responseDebugging      bool
	httpClient             *http.Client
}

func (r *requestor) execute(
	ctx context.Context,
	path string, method string, toPostInput interface{}, into any,
	reqDecorator requestDecorator, d decoder,
) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", r.baseURL, path)

	// if we are handed a reader, then don't treat it as a json input
	var toPost io.Reader
	if toPostInput != nil {
		switch v := toPostInput.(type) {
		case io.Reader:
			toPost = v
		default:
			toPostBytes, err := json.Marshal(toPostInput)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal provided body to %s:%s", method, path)
			}
			toPost = bytes.NewReader(toPostBytes)
		}
	}

	req, err := http.NewRequest(method, url, toPost)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s:%s", method, path)
	}
	req = req.WithContext(ctx)

	if r.AuthorizationDecorator != nil {
		r.AuthorizationDecorator(req)
	}

	if reqDecorator != nil {
		reqDecorator(req)
	}

	if r.requestDebugging {
		bodyDebug := ""
		switch toPostInput.(type) {
		case io.Reader:
			bodyDebug = "reader provided, will not read"
		default:
			debugJson, debugErr := json.Marshal(toPostInput)
			bodyDebug = fmt.Sprintf("%v => %s", debugErr, string(debugJson))
		}
		logrus.WithFields(logrus.Fields{
			"url":     req.URL,
			"method":  method,
			"body":    bodyDebug,
			"headers": req.Header,
		}).Debug("API request")
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to executing request to %s:%s", method, path)
	}
	defer resp.Body.Close()

	var toDecode io.Reader = resp.Body

	if r.responseDebugging {
		body, debugErr := io.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"method":     req.Method,
			"url":        req.URL,
			"statusCode": resp.StatusCode,
			// "headers":    resp.Header,
			"body": fmt.Sprintf("%v => %s", debugErr, string(body)),
		}).Debug("API response")
		toDecode = bytes.NewReader(body)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		err = d.decode(resp.Body, &errRes)
		if err != nil || errRes.Error == nil {
			reqErr := RequestError{
				StatusCode: resp.StatusCode,
				Err:        err,
			}
			return nil, fmt.Errorf("error, %w", &reqErr)
		}
		errRes.Error.StatusCode = resp.StatusCode
		return nil, fmt.Errorf("error, status code: %d, message: %w", resp.StatusCode, errRes.Error)
	}

	if resp.StatusCode == 204 {
		// No Content
		// Do not bother with response body stuff; none expected
	} else {
		if into != nil {
			if err := d.decode(toDecode, into); err != nil {
				return resp, fmt.Errorf("failed parsing the response from %s:%s", method, path)
			}
		}
	}
	return resp, nil
}

type decoder interface {
	decode(reader io.Reader, into any) error
}

type xmlDecoder struct{}

func (*xmlDecoder) decode(reader io.Reader, into any) error {
	return xml.NewDecoder(reader).Decode(into)
}

type jsonDecoder struct{}

func (*jsonDecoder) jsonDecoder(reader io.Reader, into any) error {
	return json.NewDecoder(reader).Decode(into)
}

func (r *requestor) Get(ctx context.Context, path string, into any, decorator requestDecorator, d decoder) (*http.Response, error) {
	return r.execute(ctx, path, "GET", nil, into, decorator, d)
}

//func (r *requestor) Post(ctx context.Context, path string, toPost interface{}, into interface{}) (*http.Response, error) {
//	return r.execute(ctx, path, "POST", toPost, into, jsonDecorator, decoder)
//}
//
//func (r *requestor) Patch(ctx context.Context, path string, toPatch interface{}, into interface{}) (*http.Response, error) {
//	return r.execute(ctx, path, "PATCH", toPatch, into, jsonDecorator, xml)
//}
//
//func (r *requestor) Delete(ctx context.Context, path string, into interface{}) (*http.Response, error) {
//	return r.execute(ctx, path, "DELETE", nil, into, jsonDecorator)
//}

func jsonDecorator(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	return req
}

func xmlDecorator(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/xml")
	return req
}
