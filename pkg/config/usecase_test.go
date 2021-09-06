package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UseCaseSuite struct {
	suite.Suite
}

func (suite UseCaseSuite) TestGet() {
	s := New(&Config{
		Version: "test",
	})

	c, err := s.Get()
	suite.NoError(err)
	suite.Equal("test", c.Version)
}

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
