package healthcheck

import (
	"context"
	"errors"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Ping() error {
	args := m.Called()
	return args.Error(0)
}

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) TestNew() {
	s := NewHandler(&mockRepo{})
	suite.NotNil(s)
}

func (suite *HandlerSuite) TestCheckDB() {
	cases := []struct {
		Database config.Database
		Error    error
	}{
		{Error: errors.New("connection error")},
		{Error: nil},
	}

	for _, c := range cases {
		repo := new(mockRepo)
		repo.On("Ping").Return(c.Error)
		hs := handler{repo: repo}

		err := hs.checkDB(context.TODO())
		if c.Error != nil {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *HandlerSuite) TestCheck() {
	hs := handler{repo: &mockRepo{}}
	res := hs.Check()
	suite.NotNil(res)
}
func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
