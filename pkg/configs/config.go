package configs

import (
	"github.com/joeshaw/envdecode"
	"github.com/pkg/errors"
)

// Config represents all HTTP server configuration options.
type Config struct {
	Host         string `env:"SERVER_HOST"`
	Port         int    `env:"SERVER_PORT"`
	EthereumHost string `env:"ETHEREUM_HOST"`
}

// NewConfig constructs a new instance of Config via decoding
// the mapped env vars with envdecode library.
func NewConfig() (*Config, error) {
	var config Config

	err := envdecode.Decode(&config)

	if err != nil {
		if err != envdecode.ErrNoTargetFieldsAreSet {
			return nil, errors.WithStack(err)
		}
	}

	return &config, nil
}
