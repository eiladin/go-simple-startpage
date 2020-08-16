package config

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func createConfigFile(cfgFile string) {
	content := []byte(`
db_driver: "sqlite"
db_name: "dbname.db"
db_username: "user"
db_password: "pass"
db_host: "host"
db_port: "1234"
db_log: "false"
listen_port: "3000"
timeout: "2000"
`)

	ioutil.WriteFile(cfgFile, content, 0644)
}

func TestEnvConfig(t *testing.T) {
	cfgFile := "./not-found.yaml"
	os.Setenv("GSS_DB_DRIVER", "driver")
	os.Setenv("GSS_DB_NAME", "name")
	os.Setenv("GSS_DB_USERNAME", "username")
	os.Setenv("GSS_DB_PASSWORD", "password")
	os.Setenv("GSS_DB_HOST", "host")
	os.Setenv("GSS_DB_PORT", "1")
	os.Setenv("GSS_DB_LOG", "false")
	os.Setenv("GSS_LISTEN_PORT", "2")
	os.Setenv("GSS_TIMEOUT", "3")
	c := InitConfig("1.2.3", cfgFile)
	assert.Equal(t, "driver", c.DBDriver)
	assert.Equal(t, "name", c.DBName)
	assert.Equal(t, "username", c.DBUsername)
	assert.Equal(t, "password", c.DBPassword)
	assert.Equal(t, "host", c.DBHost)
	assert.Equal(t, "1", c.DBPort)
	assert.Equal(t, false, c.DBLog)
	assert.Equal(t, 2, c.ListenPort)
	assert.Equal(t, 3, c.Timeout)
	os.Unsetenv("GSS_DB_DRIVER")
	os.Unsetenv("GSS_DB_NAME")
	os.Unsetenv("GSS_DB_USERNAME")
	os.Unsetenv("GSS_DB_PASSWORD")
	os.Unsetenv("GSS_DB_HOST")
	os.Unsetenv("GSS_DB_PORT")
	os.Unsetenv("GSS_DB_LOG")
	os.Unsetenv("GSS_LISTEN_PORT")
	os.Unsetenv("GSS_TIMEOUT")
}

func TestDefaultConfig(t *testing.T) {
	cfgFile := "./config.yaml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", "")
	assert.Equal(t, "sqlite", c.DBDriver)
	assert.Equal(t, "dbname.db", c.DBName)
	assert.Equal(t, "user", c.DBUsername)
	assert.Equal(t, "pass", c.DBPassword)
	assert.Equal(t, "host", c.DBHost)
	assert.Equal(t, "1234", c.DBPort)
	assert.Equal(t, false, c.DBLog)
	assert.Equal(t, 3000, c.ListenPort)
	assert.Equal(t, 2000, c.Timeout)
	assert.Equal(t, "1.2.3", c.Version)
}

func TestConfigFile(t *testing.T) {
	cfgFile := "./test-config-file.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", cfgFile)
	assert.Equal(t, "sqlite", c.DBDriver)
	assert.Equal(t, "dbname.db", c.DBName)
	assert.Equal(t, "user", c.DBUsername)
	assert.Equal(t, "pass", c.DBPassword)
	assert.Equal(t, "host", c.DBHost)
	assert.Equal(t, "1234", c.DBPort)
	assert.Equal(t, false, c.DBLog)
	assert.Equal(t, 3000, c.ListenPort)
	assert.Equal(t, 2000, c.Timeout)
	assert.Equal(t, "1.2.3", c.Version)
}

func TestGetAppConfig(t *testing.T) {
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
		assert.Equal(t, "{\"version\":\"1.2.3\"}\n", rec.Body.String())
	}
}

func TestGetConfig(t *testing.T) {
	cfgFile := "./test-get-config.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)
	InitConfig("1.2.3", cfgFile)
	c := GetConfig()
	assert.Equal(t, "sqlite", c.DBDriver)
	assert.Equal(t, "dbname.db", c.DBName)
	assert.Equal(t, "user", c.DBUsername)
	assert.Equal(t, "pass", c.DBPassword)
	assert.Equal(t, "host", c.DBHost)
	assert.Equal(t, "1234", c.DBPort)
	assert.Equal(t, false, c.DBLog)
	assert.Equal(t, 3000, c.ListenPort)
	assert.Equal(t, 2000, c.Timeout)
	assert.Equal(t, "1.2.3", c.Version)
}
