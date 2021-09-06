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

type UseCaseSuite struct {
	suite.Suite
}

func (suite *UseCaseSuite) TestNew() {
	s := New(&mockRepo{})
	suite.NotNil(s)
}

func (suite *UseCaseSuite) TestCheckDB() {
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
		hs := service{repo: repo}

		err := hs.checkDB(context.TODO())
		if c.Error != nil {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *UseCaseSuite) TestCheck() {
	hs := service{repo: &mockRepo{}}
	handler := hs.Check()
	suite.NotNil(handler)
}
func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
