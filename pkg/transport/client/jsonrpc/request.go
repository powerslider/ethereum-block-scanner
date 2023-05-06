package jsonrpc

import "reflect"

// RPCRequest represents a JSON-RPC request object.
//
// Method: string containing the method to be invoked
//
// Params: can be nil. if not must be an json array or object
//
// ID: may always be set to 0 (default can be changed) for single requests.
// Should be unique for every request in one batch request.
//
// JSONRPC: must always be set to "2.0" for JSON-RPC version 2.0
//
// See: http://www.jsonrpc.org/specification#request_object
//
// Most of the time you shouldn't create the RPCRequest object yourself.
// The following functions do that for you:
// Call(), CallFor(), NewRequest()
//
// If you want to create it yourself (e.g. in batch or CallRaw()), consider using Params().
// Params() is a helper function that uses the same parameter syntax as Call().
//
// e.g. to manually create an RPCRequest object:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: Params("Alex", 35, true),
//	}
//
// If you know what you are doing you can omit the Params() call to avoid some reflection,
// but potentially create incorrect rpc requests:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: 2, <-- invalid since a single primitive value must be wrapped in an array --> no magic without Params()
//	}
//
// correct:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: []int{2}, <-- invalid since a single primitive value must be wrapped in an array
//	}
type RPCRequest struct {
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
}

// NewRequest returns a new RPCRequest that can be created using the same convenient parameter syntax as Call()
//
// Default RPCRequest id is 0. If you want to use an id other than 0, use NewRequestWithID() or set the ID field
// of the returned RPCRequest manually.
//
// e.g. NewRequest("myMethod", "Alex", 35, true)
func NewRequest(method string, params ...any) *RPCRequest {
	request := &RPCRequest{
		Method:  method,
		Params:  Params(params...),
		JSONRPC: _jsonrpcVersion,
	}

	return request
}

// NewRequestWithID returns a new RPCRequest that can be created using the same convenient parameter syntax as Call()
//
// e.g. NewRequestWithID(123, "myMethod", "Alex", 35, true)
func NewRequestWithID(id int, method string, params ...any) *RPCRequest {
	request := &RPCRequest{
		ID:      id,
		Method:  method,
		Params:  Params(params...),
		JSONRPC: _jsonrpcVersion,
	}

	return request
}

// Params is a helper function that uses the same parameter syntax as Call().
// But you should consider to always use NewRequest() instead.
//
// e.g. to manually create an RPCRequest object:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: Params("Alex", 35, true),
//	}
//
// same with new request:
// request := NewRequest("myMethod", "Alex", 35, true)
//
// If you know what you are doing you can omit the Params() call but potentially create incorrect rpc requests:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: 2, <-- invalid since a single primitive value must be wrapped in an array --> no magic without Params()
//	}
//
// correct:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: []int{2}, <-- valid since a single primitive value must be wrapped in an array
//	}
func Params(params ...any) any {
	var finalParams any

	// if params was nil skip this and p stays nil
	if params != nil {
		switch len(params) {
		case 0: // no parameters were provided, do nothing so finalParam is nil and will be omitted
		case 1: // one param was provided, use it directly as is, or wrap primitive types in array
			if params[0] != nil {
				var typeOf reflect.Type

				// traverse until nil or not a pointer type
				for typeOf = reflect.TypeOf(params[0]); typeOf != nil && typeOf.Kind() == reflect.Ptr; typeOf = typeOf.Elem() {
				}

				if typeOf != nil {
					// now check if we can directly marshal the type or if it must be wrapped in an array
					switch typeOf.Kind() {
					// for these types we just do nothing, since value of p is already unwrapped from the array params
					case reflect.Struct:
						finalParams = params[0]
					case reflect.Array:
						finalParams = params[0]
					case reflect.Slice:
						finalParams = params[0]
					case reflect.Interface:
						finalParams = params[0]
					case reflect.Map:
						finalParams = params[0]
					default: // everything else must stay in an array (int, string, etc)
						finalParams = params
					}
				}
			} else {
				finalParams = params
			}
		default: // if more than one parameter was provided it should be treated as an array
			finalParams = params
		}
	}

	return finalParams
}
