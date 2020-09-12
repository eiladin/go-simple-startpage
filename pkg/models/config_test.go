package models

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) SetupTest() {
	viper.Reset()
}

func createConfigFile(t *testing.T, cfgFile string) {
	content := []byte(`database:
  driver: "sqlite"
  name: "dbname.db"
  log: "false"
listen_port: "8080"
timeout: "5000"
`)

	assert.NoError(t, ioutil.WriteFile(cfgFile, content, 0644))
}

func createErrorConfigFile(t *testing.T, cfgFile string) {
	content := []byte(`database:
driver: "sqlite1"
  name: "dbname.db"
  log: "true"
listen_port: "8080"
timeout: "3000"
`)

	assert.NoError(t, ioutil.WriteFile(cfgFile, content, 0644))
}

func (suite *ConfigSuite) TestEnvConfig() {
	cfgFile := "./not-found.yaml"
	os.Setenv("GSS_DATABASE_NAME", "name1")
	os.Setenv("GSS_DATABASE_LOG", "false")
	os.Setenv("GSS_LISTEN_PORT", "5")
	os.Setenv("GSS_TIMEOUT", "6")
	os.Setenv("GSS_ENVIRONMENT", "Production")
	c := NewConfig("1.2.3", cfgFile)
	suite.Equal("name1", c.Database.Name)
	suite.Equal(false, c.Database.Log)
	suite.Equal(5, c.ListenPort)
	suite.Equal(6, c.Timeout)
	suite.Equal("Production", c.Environment)
	os.Unsetenv("GSS_DATABASE_NAME")
	os.Unsetenv("GSS_DATABASE_LOG")
	os.Unsetenv("GSS_LISTEN_PORT")
	os.Unsetenv("GSS_TIMEOUT")
	os.Unsetenv("GSS_ENVIRONMENT")
}

func (suite *ConfigSuite) TestIsProduction() {
	cases := []struct {
		Environment string
		Expected    bool
	}{
		{"Production", true},
		{"", false},
		{"Dev", false},
		{"PRODUCTION", true},
		{"PrOdUcTiOn", true},
	}

	for _, c := range cases {
		os.Unsetenv("GSS_ENVIRONMENT")
		if c.Environment != "" {
			os.Setenv("GSS_ENVIRONMENT", c.Environment)
		}
		cfg := NewConfig("test", "not-found")
		suite.Equal(c.Expected, cfg.IsProduction(), "IsProduction should be %t", c.Expected)
	}
	os.Unsetenv("GSS_ENVIRONMENT")
}

func (suite *ConfigSuite) TestConfigFile() {
	viper.Reset()
	cfgFile := "./test-config-file.yml"
	createConfigFile(suite.T(), cfgFile)
	defer os.RemoveAll(cfgFile)

	c := NewConfig("1.2.3", cfgFile)
	suite.Equal("dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	suite.Equal(false, c.Database.Log, "Database Log should be 'false'")
	suite.Equal(8080, c.ListenPort, "Listen Port should be '8080'")
	suite.Equal(5000, c.Timeout, "Timeout should be '5000'")
	suite.Equal("1.2.3", c.Version, "Version should be '1.2.3'")
}

func (suite *ConfigSuite) TestDefaultConfigFile() {
	viper.Reset()
	cfgFile := "./config.yml"
	createConfigFile(suite.T(), cfgFile)
	defer os.RemoveAll(cfgFile)

	c := NewConfig("1.2.3", "")
	suite.Equal("dbname.db", c.Database.Name, "Database Name should be 'dbname.db'")
	suite.Equal(false, c.Database.Log, "Database Log should be 'false'")
	suite.Equal(8080, c.ListenPort, "Listen Port should be '8080'")
	suite.Equal(5000, c.Timeout, "Timeout should be '5000'")
	suite.Equal("1.2.3", c.Version, "Version should be '1.2.3'")
}

func (suite *ConfigSuite) TestConfigFileErr() {
	cfgFile := "./test-config-file-error.yml"
	createErrorConfigFile(suite.T(), cfgFile)
	defer os.RemoveAll(cfgFile)
	c := NewConfig("1.2.3", cfgFile)
	suite.Equal("simple-startpage.db", c.Database.Name, "Database Name should be 'simple-startpage.db'")
	suite.Equal(false, c.Database.Log, "Database Log should be 'false'")
	suite.Equal(3000, c.ListenPort, "Listen Port should be '3000'")
	suite.Equal(2000, c.Timeout, "Timeout should be '2000'")
	suite.Equal("1.2.3", c.Version, "Version should be '1.2.3'")
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
