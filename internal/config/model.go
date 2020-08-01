package config

// Configuration stores application configuration
type Configuration struct {
	Server      ServerConfiguration   `json:"-"`
	Database    DatabaseConfiguration `json:"-"`
	Version     string                `json:"version"`
	HealthCheck HealthCheck           `json:"-"`
}

// DatabaseConfiguration structure
type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	Log      bool
}

// HealthCheck structure
type HealthCheck struct {
	Timeout int
}
