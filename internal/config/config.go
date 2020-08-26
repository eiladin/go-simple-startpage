package config

import (
	"strings"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/spf13/viper"
)

var c models.Config

// InitConfig initializes application configuration
func New(version string, cfgFile string) models.Config {
	c = models.Config{
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
	_ = viper.BindEnv("DATABASE.DRIVER")
	_ = viper.BindEnv("DATABASE.NAME")
	_ = viper.BindEnv("DATABASE.USERNAME")
	_ = viper.BindEnv("DATABASE.PASSWORD")
	_ = viper.BindEnv("DATABASE.HOST")
	_ = viper.BindEnv("DATABASE.PORT")
	_ = viper.BindEnv("DATABASE.LOG")
	_ = viper.BindEnv("LISTEN_PORT")
	_ = viper.BindEnv("TIMEOUT")
	_ = viper.BindEnv("ENVIRONMENT")

	_ = viper.ReadInConfig()
	viper.AutomaticEnv()
	_ = viper.Unmarshal(&c)
	c.Version = version
	return c
}

// GetConfig returns application configuration
func GetConfig() models.Config {
	return c
}
