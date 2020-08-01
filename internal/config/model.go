package config

// Configuration structure
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

// ServerConfiguration structure
type ServerConfiguration struct {
	Port string
}
