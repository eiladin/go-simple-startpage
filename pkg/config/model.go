package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Filepath    string   `json:"-"`
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

func Load(version string, cfgFile string) *Config {
	c := &Config{
		Environment: "Development",
		Version:     version,
		ListenPort:  3000,
		Timeout:     2000,
		Database: Database{
			Driver: "sqlite",
			Name:   "simple-startpage.db",
			Log:    false,
		},
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
	_ = viper.BindEnv("FILEPATH")
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
	return c
}
