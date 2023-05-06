package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// RPCResponse represents a JSON-RPC response object.
//
// Result: holds the result of the rpc call if no error occurred, nil otherwise. can be nil even on success.
//
// Error: holds an RPCError object if an error occurred. must be nil on success.
//
// ID: may always be 0 for single requests. is unique for each request in a batch call (see CallBatch())
//
// JSONRPC: must always be set to "2.0" for JSON-RPC version 2.0
//
// See: http://www.jsonrpc.org/specification#response_object
type RPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	ID      int       `json:"id"`
}

// GetInt converts the rpc response to an int64 and returns it.
//
// If result was not an integer an error is returned.
func (r *RPCResponse) GetInt() (int64, error) {
	val, ok := r.Result.(json.Number)
	if !ok {
		return 0, fmt.Errorf("could not parse int64 from %s", r.Result)
	}

	i, err := val.Int64()
	if err != nil {
		return 0, err
	}

	return i, nil
}

// GetFloat converts the rpc response to float64 and returns it.
//
// If result was not an float64 an error is returned.
func (r *RPCResponse) GetFloat() (float64, error) {
	val, ok := r.Result.(json.Number)
	if !ok {
		return 0, fmt.Errorf("could not parse float64 from %s", r.Result)
	}

	f, err := val.Float64()
	if err != nil {
		return 0, err
	}

	return f, nil
}

// GetBool converts the rpc response to a bool and returns it.
//
// If result was not a bool an error is returned.
func (r *RPCResponse) GetBool() (bool, error) {
	val, ok := r.Result.(bool)
	if !ok {
		return false, fmt.Errorf("could not parse bool from %s", r.Result)
	}

	return val, nil
}

// GetString converts the rpc response to a string and returns it.
//
// If result was not a string an error is returned.
func (r *RPCResponse) GetString() (string, error) {
	val, ok := r.Result.(string)
	if !ok {
		return "", fmt.Errorf("could not parse string from %s", r.Result)
	}

	return val, nil
}

// GetObject converts the rpc response to an arbitrary type.
//
// The function works as you would expect it from json.Unmarshal()
func (r *RPCResponse) GetObject(toType any) error {
	js, err := json.Marshal(r.Result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, toType)
	if err != nil {
		return err
	}

	return nil
}
