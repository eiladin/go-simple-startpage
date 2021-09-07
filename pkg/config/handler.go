package config

type IHandler interface {
	Get() (*Config, error)
}

// Compile-time proof of interface implementation.
var _ IHandler = (*handler)(nil)

type handler struct {
	config *Config
}

func New(cfg *Config) IHandler {
	return &handler{
		config: cfg,
	}
}

func (c *handler) Get() (*Config, error) {
	return c.config, nil
}
