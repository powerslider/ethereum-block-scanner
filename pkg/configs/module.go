package configs

// InitializeConfig wires all dependencies for the config module.
func InitializeConfig() *Config {
	conf, err := NewConfig()
	if err != nil {
		panic("error loading application config")
	}

	return conf
}
