package httpx

import (
	"io"
	nethttp "net/http"
)

const (
	_defaultRequestTimeout      = 10
	_defaultTLSHandshakeTimeout = 5
	_defaultDialerTimeout       = 5
)

type clientConfig struct {
	requestTimeout      int
	tlsHandshakeTimeout int
	dialerTimeout       int
	redirectPolicy      RedirectPolicy
}

// RedirectPolicy specifies the policy for handling redirects.
type RedirectPolicy func(req *nethttp.Request, via []*nethttp.Request) error

// ClientOption specifies a HTTP client setting.
type ClientOption func(config *clientConfig)

func newHTTPClientDefaultConfig() *clientConfig {
	return &clientConfig{
		requestTimeout:      _defaultRequestTimeout,
		tlsHandshakeTimeout: _defaultTLSHandshakeTimeout,
		dialerTimeout:       _defaultDialerTimeout,
		redirectPolicy: func(req *nethttp.Request, via []*nethttp.Request) error {
			return nil
		},
	}
}

func (o *clientConfig) applyOptions(opts ...ClientOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithCustomTLSHandshakeTimeout specifies TLS handshake timeout.
func WithCustomTLSHandshakeTimeout(timeout int) ClientOption {
	return func(o *clientConfig) {
		o.tlsHandshakeTimeout = timeout
	}
}

// WithCustomRequestTimeout specifies request timeout.
func WithCustomRequestTimeout(timeout int) ClientOption {
	return func(o *clientConfig) {
		o.requestTimeout = timeout
	}
}

// WithCustomDialerTimeout specifies dialer timeout.
func WithCustomDialerTimeout(timeout int) ClientOption {
	return func(o *clientConfig) {
		o.dialerTimeout = timeout
	}
}

// WithCustomRedirectPolicy specifies a RedirectPolicy.
func WithCustomRedirectPolicy(policy RedirectPolicy) ClientOption {
	return func(o *clientConfig) {
		o.redirectPolicy = policy
	}
}

type requestConfig struct {
	queryParams map[string]string
	formParams  map[string]string
	body        io.Reader
	basicAuth   *basicAuth
}

type basicAuth struct {
	username string
	password string
}

// RequestOption specifies a HTTP request option.
type RequestOption func(config *requestConfig)

func newHTTPRequestDefaultConfig() *requestConfig {
	return &requestConfig{
		basicAuth:   &basicAuth{},
		queryParams: make(map[string]string),
		body:        nil,
	}
}

// ApplyOptions applies request options.
func (o *requestConfig) ApplyOptions(opts ...RequestOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithRequestBody specifies a request body.
func WithRequestBody(body io.Reader) RequestOption {
	return func(o *requestConfig) {
		o.body = body
	}
}

// WithQueryParam adds a query parameter.
func WithQueryParam(key, val string) RequestOption {
	return func(o *requestConfig) {
		o.queryParams[key] = val
	}
}

// WithFormParam adds a form parameter.
func WithFormParam(key, val string) RequestOption {
	return func(o *requestConfig) {
		o.formParams[key] = val
	}
}

// WithBasicAuth adds a basic authentication header.
func WithBasicAuth(username, password string) RequestOption {
	return func(o *requestConfig) {
		o.basicAuth.username = username
		o.basicAuth.password = password
	}
}
