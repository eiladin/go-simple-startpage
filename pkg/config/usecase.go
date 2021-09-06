package config

type IConfig interface {
	Get() (*Config, error)
}

// Compile-time proof of interface implementation.
var _ IConfig = (*service)(nil)

type service struct {
	config *Config
}

func New(cfg *Config) IConfig {
	return &service{
		config: cfg,
	}
}

func (c *service) Get() (*Config, error) {
	return c.config, nil
}
