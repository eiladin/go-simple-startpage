package status

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetSite(site *models.Site) error {
	args := m.Called(site)
	return args.Error(0)
}

type StatusSuite struct {
	suite.Suite
}

func (suite *StatusSuite) TestNew() {
	h := New(&mockRepo{}, &config.Config{})
	suite.NotNil(h)
}

func (suite *StatusSuite) TestGet() {
	cases := []struct {
		id      uint
		wantErr error
	}{
		{id: 1, wantErr: nil},
		{id: 2, wantErr: ErrNotFound},
	}

	for _, c := range cases {
		cfg := &config.Config{Timeout: 100}
		r := new(mockRepo)
		r.On("GetSite", &models.Site{ID: c.id}).Return(c.wantErr)
		ss := service{repo: r, config: cfg}

		status, err := ss.Get(c.id)
		r.AssertExpectations(suite.T())
		if c.wantErr != nil {
			suite.EqualError(err, c.wantErr.Error())
		} else {
			suite.NoError(err)
			suite.NotNil(status)
		}
	}
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}
