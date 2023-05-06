package httpx

import (
	"io"
	"net"
	"net/http"
	"time"
)

// Error represents an error that occurred on HTTP level.
type Error struct {
	Code int
	Err  error
}

// Error function is provided to be used as error object.
func (e *Error) Error() string {
	return e.Err.Error()
}

// Client represents a HTTP client wrapper.
type Client struct {
	clientInstance *http.Client
}

// NewClient is a constructor function for Client.
func NewClient(options ...ClientOption) *Client {
	config := newHTTPClientDefaultConfig()

	config.applyOptions(options...)

	netTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(config.dialerTimeout) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: time.Duration(config.tlsHandshakeTimeout) * time.Second,
	}

	client := &http.Client{
		Timeout:       time.Duration(config.requestTimeout) * time.Second,
		Transport:     netTransport,
		CheckRedirect: config.redirectPolicy,
	}

	return &Client{
		clientInstance: client,
	}
}

// Do executes a passed HTTP request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.clientInstance.Do(req)
}

// Get executes a HTTP GET request.
func (c *Client) Get(url string, options ...RequestOption) (*http.Response, error) {
	req, err := GetRequest(url, options...)

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Post executes a HTTP POST request.
func (c *Client) Post(url string, body io.Reader, options ...RequestOption) (*http.Response, error) {
	req, err := PostRequest(url, body, options...)

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
