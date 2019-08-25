package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server      ServerConfiguration   `json:"-"`
	Database    DatabaseConfiguration `json:"-"`
	Version     string                `json:"version"`
	HealthCheck HealthCheck           `json:"-"`
}

var configuration Configuration

func InitConfig() Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return configuration
}

func GetConfig() Configuration {
	return configuration
}
