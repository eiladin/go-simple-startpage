package config

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Configuration structure
type Configuration struct {
	Database   Database `json:"-"`
	ListenPort int      `mapstructure:"listen_port" yaml:"listen_port" json:"-"`
	Timeout    int      `yaml:"timeout" json:"-"`
	Version    string   `json:"version"`
}

// Database structure
type Database struct {
	Driver   string
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	Log      bool
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
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.BindEnv("DATABASE.DRIVER")
	viper.BindEnv("DATABASE.NAME")
	viper.BindEnv("DATABASE.USERNAME")
	viper.BindEnv("DATABASE.PASSWORD")
	viper.BindEnv("DATABASE.HOST")
	viper.BindEnv("DATABASE.PORT")
	viper.BindEnv("DATABASE.LOG")
	viper.BindEnv("LISTEN_PORT")
	viper.BindEnv("TIMEOUT")

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
