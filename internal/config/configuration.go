package config

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var configuration Configuration

// InitConfig initializes application configuration
func InitConfig(version string, cfgFile string) Configuration {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
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

	configuration.Version = version
	return configuration
}

// GetConfig returns application configuration
func GetConfig() Configuration {
	return configuration
}

// GetAppConfig handles /api/appconfig
func GetAppConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, configuration)
}
