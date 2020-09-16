package providers

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProviderSuite struct {
	suite.Suite
}

type mockStore struct {
	mock.Mock
}

func (m *mockStore) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockStore) CreateNetwork(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

func (m *mockStore) GetNetwork(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

func (m *mockStore) GetSite(site *models.Site) error {
	args := m.Called(site)
	return args.Error(0)
}

func (suite *ProviderSuite) TestInitProvider() {
	s := new(mockStore)
	s.On("Ping").Return(nil)
	s.On("CreateNetwork", nil).Return(nil)
	s.On("GetNetwork", nil).Return(nil)
	s.On("GetSite", nil).Return(nil)
	prv := InitProvider(&config.Config{}, &mockStore{})
	suite.NotNil(prv.Config)
	suite.NotNil(prv.Healthcheck)
	suite.NotNil(prv.Network)
	suite.NotNil(prv.Status)
	s.AssertNotCalled(suite.T(), "Ping")
	s.AssertNotCalled(suite.T(), "CreateNetwork", nil)
	s.AssertNotCalled(suite.T(), "GetNetwork", nil)
	s.AssertNotCalled(suite.T(), "GetSite", nil)

}

func TestProviderSuite(t *testing.T) {
	suite.Run(t, new(ProviderSuite))
}
