package config

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Configuration structure
type Configuration struct {
	DBDriver   string `mapstructure:"db_driver" yaml:"db_driver" json:"-"`
	DBName     string `mapstructure:"db_name" yaml:"db_name" json:"-"`
	DBUsername string `mapstructure:"db_username" yaml:"db_username" json:"-"`
	DBPassword string `mapstructure:"db_password" yaml:"db_password" json:"-"`
	DBHost     string `mapstructure:"db_host" yaml:"dh_host" json:"-"`
	DBPort     string `mapstructure:"db_port" yaml:"dh_port" json:"-"`
	DBLog      bool   `mapstructure:"db_log" yaml:"dh_log" json:"-"`
	ListenPort int    `mapstructure:"listen_port" yaml:"listen_port" json:"-"`
	Timeout    int    `mapstructure:"timeout" yaml:"timeout" json:"-"`
	Version    string `json:"version"`
}

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
	viper.SetEnvPrefix("GSS")
	_ = viper.BindEnv("DB_DRIVER")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("DB_USERNAME")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_LOG")
	_ = viper.BindEnv("LISTEN_PORT")
	_ = viper.BindEnv("TIMEOUT")

	viper.ReadInConfig()
	viper.AutomaticEnv()
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
