package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"

	"fmt"

	"github.com/powerslider/ethereum-block-scanner/pkg/transport/client/httpx"
)

const (
	_jsonrpcVersion = "2.0"
)

// RPCClient represents a JSON-RPC client.
type RPCClient struct {
	endpoint           string
	httpClient         *httpx.Client
	allowUnknownFields bool
	defaultRequestID   int
}

// NewDefaultClient returns a new RPCClient instance with default configuration.
func NewDefaultClient(endpoint string, opts ...ClientOption) *RPCClient {
	return NewClient(httpx.NewClient(), endpoint, opts...)
}

// NewClient returns a new RPCClient instance with custom configuration.
func NewClient(httpClient *httpx.Client, endpoint string, opts ...ClientOption) *RPCClient {
	config := newRPCClientDefaultConfig()

	config.applyOptions(opts...)

	rpcClient := &RPCClient{
		endpoint:   endpoint,
		httpClient: httpClient,
	}

	rpcClient.allowUnknownFields = config.allowUnknownFields
	rpcClient.defaultRequestID = config.defaultRequestID

	return rpcClient
}

// Call calls a JSON-RPC method with optional params.
func (c *RPCClient) Call(ctx context.Context, method string, params ...any) (*RPCResponse, error) {
	request := &RPCRequest{
		ID:      c.defaultRequestID,
		Method:  method,
		Params:  Params(params...),
		JSONRPC: _jsonrpcVersion,
	}

	return c.doCall(ctx, request)
}

// CallRaw calls a JSON-RPC method by passing a RPCRequest.
func (c *RPCClient) CallRaw(ctx context.Context, request *RPCRequest) (*RPCResponse, error) {
	return c.doCall(ctx, request)
}

// CallFor calls a JSON-RPC method and deserializes the response in a specified response object.
func (c *RPCClient) CallFor(ctx context.Context, out any, method string, params ...any) error {
	rpcResponse, err := c.Call(ctx, method, params...)
	if err != nil {
		return err
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	return rpcResponse.GetObject(out)
}

func (c *RPCClient) doCall(
	_ context.Context, rpcReq *RPCRequest, options ...httpx.RequestOption) (*RPCResponse, error) {
	body, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := httpx.PostRequest(c.endpoint, bytes.NewReader(body), options...)
	if err != nil {
		return nil, err
	}

	redactedURL := httpReq.URL.Redacted()

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("rpc call %v() on %v: %w", rpcReq.Method, redactedURL, err)
	}

	//nolint:errcheck
	defer httpResp.Body.Close()

	var rpcResponse *RPCResponse

	decoder := json.NewDecoder(httpResp.Body)
	if !c.allowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	decoder.UseNumber()

	err = decoder.Decode(&rpcResponse)

	return rpcResponse, HandleResponseError(err, httpResp, rpcReq, redactedURL, rpcResponse)
}
