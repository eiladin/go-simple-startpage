package status

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetSite(site *network.Site) error {
	args := m.Called(site)
	return args.Error(0)
}

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) TestNew() {
	h := New(&mockRepo{}, &config.Config{})
	suite.NotNil(h)
}

func (suite *HandlerSuite) TestGet() {
	cases := []struct {
		name    string
		wantErr error
	}{
		{name: "test-site-1", wantErr: nil},
		{name: "test-site-2", wantErr: ErrNotFound},
	}

	for _, c := range cases {
		cfg := &config.Config{Timeout: 100}
		r := new(mockRepo)
		r.On("GetSite", &network.Site{Name: c.name}).Return(c.wantErr)
		ss := handler{repo: r, config: cfg}

		status, err := ss.Get(c.name)
		r.AssertExpectations(suite.T())
		if c.wantErr != nil {
			suite.EqualError(err, c.wantErr.Error())
		} else {
			suite.NoError(err)
			suite.NotNil(status)
		}
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
