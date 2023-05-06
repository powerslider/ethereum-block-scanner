package jsonrpc

type clientConfig struct {
	allowUnknownFields bool
	defaultRequestID   int
}

func newRPCClientDefaultConfig() *clientConfig {
	return &clientConfig{
		allowUnknownFields: true,
		defaultRequestID:   0,
	}
}

func (o *clientConfig) applyOptions(opts ...ClientOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// ClientOption specifies a Client setting.
type ClientOption func(config *clientConfig)

// WithAllowUnknownFields specifies if unknown fields should be allowed in the response object.
func WithAllowUnknownFields(allowUnknownFields bool) ClientOption {
	return func(o *clientConfig) {
		o.allowUnknownFields = allowUnknownFields
	}
}

// WithDefaultRequestID sets the default ID for JSON-RPC request.
func WithDefaultRequestID(defaultRequestID int) ClientOption {
	return func(o *clientConfig) {
		o.defaultRequestID = defaultRequestID
	}
}
