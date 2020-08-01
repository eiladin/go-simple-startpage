package config

import (
	"log"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/spf13/viper"
)

var configuration Configuration

// InitConfig initializes application configuration
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

// GetConfig returns application configuration
func GetConfig() Configuration {
	return configuration
}

// GetAppConfig handles /api/appconfig
func GetAppConfig(c *fiber.Ctx) {
	c.JSON(configuration)
}
