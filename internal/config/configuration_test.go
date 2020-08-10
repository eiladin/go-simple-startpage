package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func createConfigFile(cfgFile string) {
	content := `database:
  driver: "sqlite"
  dbname: "dbname.db"
  username: "user"
  password: "pass"
  host: "host"
  port: "1234"
  log: "false"

server:
  port: "3000"

healthCheck:
  timeout: "2000"
`

	cfg := []byte(content)
	ioutil.WriteFile(cfgFile, cfg, 0644)
}

func TestDefaultConfig(t *testing.T) {
	cfgFile := "./config.yaml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", "")
	assert.Equal(t, "sqlite", c.Database.Driver)
	assert.Equal(t, "dbname.db", c.Database.Dbname)
	assert.Equal(t, "user", c.Database.Username)
	assert.Equal(t, "pass", c.Database.Password)
	assert.Equal(t, "host", c.Database.Host)
	assert.Equal(t, "1234", c.Database.Port)
	assert.Equal(t, false, c.Database.Log)
	assert.Equal(t, "3000", c.Server.Port)
	assert.Equal(t, 2000, c.HealthCheck.Timeout)
	assert.Equal(t, "1.2.3", c.Version)
}

func TestConfigFile(t *testing.T) {
	cfgFile := "./test-config-file.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)

	c := InitConfig("1.2.3", cfgFile)
	assert.Equal(t, "sqlite", c.Database.Driver)
	assert.Equal(t, "dbname.db", c.Database.Dbname)
	assert.Equal(t, "user", c.Database.Username)
	assert.Equal(t, "pass", c.Database.Password)
	assert.Equal(t, "host", c.Database.Host)
	assert.Equal(t, "1234", c.Database.Port)
	assert.Equal(t, false, c.Database.Log)
	assert.Equal(t, "3000", c.Server.Port)
	assert.Equal(t, 2000, c.HealthCheck.Timeout)
	assert.Equal(t, "1.2.3", c.Version)
}

func TestGetAppConfig(t *testing.T) {
	cfgFile := "./test-get-app-config.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)
	InitConfig("1.2.3", cfgFile)
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	GetAppConfig(ctx)
	assert.Equal(t, `{"version":"1.2.3"}`, string(ctx.Fasthttp.Response.Body()))
}

func TestGetConfig(t *testing.T) {
	cfgFile := "./test-get-config.yml"
	createConfigFile(cfgFile)
	defer os.RemoveAll(cfgFile)
	InitConfig("1.2.3", cfgFile)
	c := GetConfig()
	assert.Equal(t, "sqlite", c.Database.Driver)
	assert.Equal(t, "dbname.db", c.Database.Dbname)
	assert.Equal(t, "user", c.Database.Username)
	assert.Equal(t, "pass", c.Database.Password)
	assert.Equal(t, "host", c.Database.Host)
	assert.Equal(t, "1234", c.Database.Port)
	assert.Equal(t, false, c.Database.Log)
	assert.Equal(t, "3000", c.Server.Port)
	assert.Equal(t, 2000, c.HealthCheck.Timeout)
	assert.Equal(t, "1.2.3", c.Version)
}
