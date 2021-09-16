package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func (suite HandlerSuite) TestGet() {
	s := NewHandler(&Config{
		Version: "test",
	})

	c, err := s.Get()
	suite.NoError(err)
	suite.Equal("test", c.Version)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
