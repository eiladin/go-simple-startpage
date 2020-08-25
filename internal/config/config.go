package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config structure
type Config struct {
	Database    Database `json:"-"`
	ListenPort  int      `mapstructure:"listen_port" yaml:"listen_port" json:"-"`
	Timeout     int      `json:"-"`
	Version     string   `json:"version"`
	Environment string   `json:"-"`
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

var c Config

// IsProduction returns true if running in production mode
func (c Config) IsProduction() bool {
	return strings.ToUpper(c.Environment) == "PRODUCTION"
}

// InitConfig initializes application configuration
func InitConfig(version string, cfgFile string) Config {
	c = Config{
		Environment: "Development",
	}
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
	viper.BindEnv("ENVIRONMENT")

	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.Unmarshal(&c)
	c.Version = version
	return c
}

// GetConfig returns application configuration
func GetConfig() Config {
	return c
}
