package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func createConfigFile(cfgFile string) {
	content := []byte(`database:
  driver: "sqlite"
  name: "dbname.db"
  log: "false"
listen_port: "3000"
timeout: "2000"
`)

	ioutil.WriteFile(cfgFile, content, 0644)
}

func createErrorConfigFile(cfgFile string) {
	content := []byte(`database:
driver: "sqlite1"
  name: "dbname.db"
  log: "false"
listen_port: "3000"
timeout: "2000"
`)

	ioutil.WriteFile(cfgFile, content, 0644)
}

func TestEnvConfig(t *testing.T) {
	viper.Reset()
	cfgFile := "./not-found.yaml"
	os.Setenv("GSS_DATABASE_NAME", "name1")
	os.Setenv("GSS_DATABASE_LOG", "false")
	os.Setenv("GSS_LISTEN_PORT", "5")
	os.Setenv("GSS_TIMEOUT", "6")
	os.Setenv("GSS_ENVIRONMENT", "Production")
	c := InitConfig("1.2.3", cfgFile)
	assert.Equal(t, "name1", c.Database.Name)
	assert.Equal(t, false, c.Database.Log)
	assert.Equal(t, 5, c.ListenPort)
	assert.Equal(t, 6, c.Timeout)
	assert.Equal(t, "Production", c.Environment)
	os.Unsetenv("GSS_DATABASE_NAME")
	os.Unsetenv("GSS_DATABASE_LOG")
	os.Unsetenv("GSS_LISTEN_PORT")
	os.Unsetenv("GSS_TIMEOUT")
	os.Unsetenv("GSS_ENVIRONMENT")
}

func TestIsProduction(t *testing.T) {
	cases := []struct {
		Environment  string
		IsProduction bool
	}{
		{"Production", true},
		{"", false},
		{"Dev", false},
		{"PRODUCTION", true},
		{"PrOdUcTiOn", true},
	}

	for _, c := range cases {
		viper.Reset()
		if c.Environment != "" {
			os.Setenv("GSS_ENVIRONMENT", c.Environment)
		}
		cfg := InitConfig("test", "not-found")
		if c.IsProduction {
			assert.True(t, cfg.IsProduction(), "IsProduction should be true")
		} else {
			assert.False(t, cfg.IsProduction(), "IsProduction should be false")
		}
		if c.Environment != "" {
			os.Unsetenv("GSS_ENVIRONMENT")
		}
	}

}

func TestDefaultConfig(t *testing.T) {
	viper.Reset()
	cfgFile := "./config.yaml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", "")
	assert.Equal(t, "dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	assert.Equal(t, false, c.Database.Log, "Database Log should be 'false'")
	assert.Equal(t, 3000, c.ListenPort, "Listen Port should be '3000'")
	assert.Equal(t, 2000, c.Timeout, "Timeout should be '2000'")
	assert.Equal(t, "1.2.3", c.Version, "Version should be '1.2.3'")
}

func TestConfigFile(t *testing.T) {
	viper.Reset()
	cfgFile := "./test-config-file.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", cfgFile)
	assert.Equal(t, "dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	assert.Equal(t, false, c.Database.Log, "Database Log should be 'false'")
	assert.Equal(t, 3000, c.ListenPort, "Listen Port should be '3000'")
	assert.Equal(t, 2000, c.Timeout, "Timeout should be '2000'")
	assert.Equal(t, "1.2.3", c.Version, "Version should be '1.2.3'")
}

func TestConfigFileErr(t *testing.T) {
	viper.Reset()
	cfgFile := "./test-config-file-error.yml"
	createErrorConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)
	c := InitConfig("1.2.3", cfgFile)
	assert.NotEqual(t, "dbname.db", c.Database.Name, "Database Name should not be 'dbname.db'")
	assert.Equal(t, false, c.Database.Log, "Database Log should be 'false'")
	assert.NotEqual(t, 3000, c.ListenPort, "Listen Port should not be '3000'")
	assert.NotEqual(t, 2000, c.Timeout, "Timeout should not be '2000'")
	assert.Equal(t, "1.2.3", c.Version, "Version should be '1.2.3'")
}

func TestGetConfig(t *testing.T) {
	viper.Reset()
	cfgFile := "./test-get-config.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)
	InitConfig("1.2.3", cfgFile)
	c := GetConfig()
	assert.Equal(t, "dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	assert.Equal(t, false, c.Database.Log, "Database Log should be 'false'")
	assert.Equal(t, 3000, c.ListenPort, "Listen Port should be '3000'")
	assert.Equal(t, 2000, c.Timeout, "Timeout should be '2000'")
	assert.Equal(t, "1.2.3", c.Version, "Version should be '1.2.3'")
}
