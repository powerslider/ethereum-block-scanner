package jsonrpc

import (
	"fmt"

	"github.com/powerslider/ethereum-block-scanner/pkg/transport/client/httpx"

	"net/http"
	"strconv"
)

// RPCError represents a JSON-RPC error object if an RPC error occurred.
// See: http://www.jsonrpc.org/specification#error_object
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error function is provided to be used as error object.
func (e *RPCError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}

// HandleResponseError handles errouneous responses.
func HandleResponseError(
	err error, httpResp *http.Response, rpcReq *RPCRequest, redactedURL string, rpcResponse *RPCResponse) error {
	if err != nil {
		// if we have some http error, return it
		if httpResp.StatusCode >= 400 {
			return &httpx.Error{
				Code: httpResp.StatusCode,
				Err: fmt.Errorf(
					"rpc call %v() on %v status code: %v. could not decode body to rpc response: %w",
					rpcReq.Method, redactedURL, httpResp.StatusCode, err),
			}
		}

		return fmt.Errorf(
			"rpc call %v() on %v status code: %v. could not decode body to rpc response: %w",
			rpcReq.Method, redactedURL, httpResp.StatusCode, err)
	}

	// response body empty
	if rpcResponse == nil {
		// if we have some http error, return it
		if httpResp.StatusCode >= 400 {
			return &httpx.Error{
				Code: httpResp.StatusCode,
				Err: fmt.Errorf(
					"rpc call %v() on %v status code: %v. rpc response missing",
					rpcReq.Method, redactedURL, httpResp.StatusCode),
			}
		}

		return fmt.Errorf(
			"rpc call %v() on %v status code: %v. rpc response missing",
			rpcReq.Method, redactedURL, httpResp.StatusCode)
	}

	// if we have a response body, but also a http error situation, return both
	if httpResp.StatusCode >= 400 {
		if rpcResponse.Error != nil {
			return &httpx.Error{
				Code: httpResp.StatusCode,
				Err: fmt.Errorf(
					"rpc call %v() on %v status code: %v. rpc response error: %v",
					rpcReq.Method, redactedURL, httpResp.StatusCode, rpcResponse.Error),
			}
		}

		return &httpx.Error{
			Code: httpResp.StatusCode,
			Err: fmt.Errorf(
				"rpc call %v() on %v status code: %v. no rpc error available",
				rpcReq.Method, redactedURL, httpResp.StatusCode),
		}
	}

	return nil
}
