package config

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (suite ConfigSuite) TestGet() {
	s := New(&config.Config{
		Version: "test",
	})

	c, err := s.Get()
	suite.NoError(err)
	suite.Equal("test", c.Version)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
