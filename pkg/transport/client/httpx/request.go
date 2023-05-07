package httpx

import (
	"context"
	"io"
	"net/http"
	nethttp "net/http"
	"strings"
)

// Method represents an enum type for an HTTP method.
type Method int

const (
	// GetMethod represents a GET request method.
	GetMethod Method = iota
	// PostMethod represents a POST request method.
	PostMethod
)

// String returns the string value for an HTTP method.
func (m Method) String() string {
	return [...]string{"GET", "POST"}[m]
}

// GetRequest returns a GET HTTP request object.
func GetRequest(ctx context.Context, url string, options ...RequestOption) (*nethttp.Request, error) {
	return NewRequest(ctx, GetMethod, url, options...)
}

// PostRequest returns a POST HTTP request object.
func PostRequest(ctx context.Context, url string, body io.Reader, options ...RequestOption) (*nethttp.Request, error) {
	options = append(options, WithRequestBody(body))

	return NewRequest(ctx, PostMethod, url, options...)
}

// NewRequest builds and returns a customizable HTTP request object.
func NewRequest(ctx context.Context, method Method, url string, options ...RequestOption) (*nethttp.Request, error) {
	config := newHTTPRequestDefaultConfig()

	config.ApplyOptions(options...)

	var body io.Reader

	if config.body != nil {
		body = config.body
	}

	req, err := http.NewRequestWithContext(ctx, method.String(), url, body)

	if len(config.queryParams) > 0 {
		reqURL := req.URL
		params := reqURL.Query()

		for k, v := range config.queryParams {
			params.Add(k, v)
		}

		reqURL.RawQuery = params.Encode()
	}

	if len(config.formParams) > 0 {
		reqURL := req.URL
		form := reqURL.Query()

		for k, v := range config.formParams {
			form.Add(k, v)
		}

		req.PostForm = form

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		config.body = strings.NewReader(form.Encode())
	}

	auth := config.basicAuth
	if auth != nil {
		req.SetBasicAuth(auth.username, auth.password)
	}

	return req, err
}
