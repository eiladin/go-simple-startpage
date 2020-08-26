package models

import "strings"

type Config struct {
	Database    Database `json:"-"`
	ListenPort  int      `mapstructure:"listen_port" yaml:"listen_port" json:"-"`
	Timeout     int      `json:"-"`
	Version     string   `json:"version"`
	Environment string   `json:"-"`
}

type Database struct {
	Driver   string
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	Log      bool
}

func (c Config) IsProduction() bool {
	return strings.ToUpper(c.Environment) == "PRODUCTION"
}
