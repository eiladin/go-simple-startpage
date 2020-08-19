package config

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func createConfigFile(cfgFile string) {
	content := []byte(`database:
  driver: "sqlite"
  name: "dbname.db"
  username: "user"
  password: "pass"
  host: "host"
  port: "1234"
  log: "false"
listen_port: "3000"
timeout: "2000"
`)

	ioutil.WriteFile(cfgFile, content, 0644)
}

func TestEnvConfig(t *testing.T) {
	viper.Reset()
	cfgFile := "./not-found.yaml"
	os.Setenv("GSS_DATABASE_DRIVER", "driver")
	os.Setenv("GSS_DATABASE_NAME", "name")
	os.Setenv("GSS_DATABASE_USERNAME", "username")
	os.Setenv("GSS_DATABASE_PASSWORD", "password")
	os.Setenv("GSS_DATABASE_HOST", "host")
	os.Setenv("GSS_DATABASE_PORT", "1")
	os.Setenv("GSS_DATABASE_LOG", "false")
	os.Setenv("GSS_LISTEN_PORT", "2")
	os.Setenv("GSS_TIMEOUT", "3")
	c := InitConfig("1.2.3", cfgFile)
	assert.Equal(t, "driver", c.Database.Driver)
	assert.Equal(t, "name", c.Database.Name)
	assert.Equal(t, "username", c.Database.Username)
	assert.Equal(t, "password", c.Database.Password)
	assert.Equal(t, "host", c.Database.Host)
	assert.Equal(t, "1", c.Database.Port)
	assert.Equal(t, false, c.Database.Log)
	assert.Equal(t, 2, c.ListenPort)
	assert.Equal(t, 3, c.Timeout)
	os.Unsetenv("GSS_DATABASE_DRIVER")
	os.Unsetenv("GSS_DATABASE_NAME")
	os.Unsetenv("GSS_DATABASE_USERNAME")
	os.Unsetenv("GSS_DATABASE_PASSWORD")
	os.Unsetenv("GSS_DATABASE_HOST")
	os.Unsetenv("GSS_DATABASE_PORT")
	os.Unsetenv("GSS_DATABASE_LOG")
	os.Unsetenv("GSS_LISTEN_PORT")
	os.Unsetenv("GSS_TIMEOUT")
}

func TestDefaultConfig(t *testing.T) {
	viper.Reset()
	cfgFile := "./config.yaml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", "")
	assert.Equal(t, "sqlite", c.Database.Driver, "Database Driver should be 'sqlite'")
	assert.Equal(t, "dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	assert.Equal(t, "user", c.Database.Username, "Database Username should be 'user'")
	assert.Equal(t, "pass", c.Database.Password, "Database Password should be 'pass'")
	assert.Equal(t, "host", c.Database.Host, "Database Host should be 'host'")
	assert.Equal(t, "1234", c.Database.Port, "Database Port should be '1234'")
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
	assert.Equal(t, "sqlite", c.Database.Driver, "Database Driver should be 'sqlite'")
	assert.Equal(t, "dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	assert.Equal(t, "user", c.Database.Username, "Database Username should be 'user'")
	assert.Equal(t, "pass", c.Database.Password, "Database Password should be 'pass'")
	assert.Equal(t, "host", c.Database.Host, "Database Host should be 'host'")
	assert.Equal(t, "1234", c.Database.Port, "Database Port should be '1234'")
	assert.Equal(t, false, c.Database.Log, "Database Log should be 'false'")
	assert.Equal(t, 3000, c.ListenPort, "Listen Port should be '3000'")
	assert.Equal(t, 2000, c.Timeout, "Timeout should be '2000'")
	assert.Equal(t, "1.2.3", c.Version, "Version should be '1.2.3'")
}

func TestGetAppConfig(t *testing.T) {
	viper.Reset()
	configuration = Configuration{
		Version: "1.2.3",
	}
	app := echo.New()
	app.GET("/", GetAppConfig)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := app.NewContext(req, rec)

	if assert.NoError(t, GetAppConfig(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"version\":\"1.2.3\"}\n", rec.Body.String(), "Version should be 1.2.3")
	}
}

func TestGetConfig(t *testing.T) {
	viper.Reset()
	cfgFile := "./test-get-config.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)
	InitConfig("1.2.3", cfgFile)
	c := GetConfig()
	assert.Equal(t, "sqlite", c.Database.Driver, "Database Driver should be 'sqlite'")
	assert.Equal(t, "dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	assert.Equal(t, "user", c.Database.Username, "Database Username should be 'user'")
	assert.Equal(t, "pass", c.Database.Password, "Database Password should be 'pass'")
	assert.Equal(t, "host", c.Database.Host, "Database Host should be 'host'")
	assert.Equal(t, "1234", c.Database.Port, "Database Port should be '1234'")
	assert.Equal(t, false, c.Database.Log, "Database Log should be 'false'")
	assert.Equal(t, 3000, c.ListenPort, "Listen Port should be '3000'")
	assert.Equal(t, 2000, c.Timeout, "Timeout should be '2000'")
	assert.Equal(t, "1.2.3", c.Version, "Version should be '1.2.3'")
}
