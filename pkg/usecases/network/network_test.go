package network

import (
	"errors"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateNetwork(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

func (m *mockRepo) GetNetwork(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

type NetworkSuite struct {
	suite.Suite
}

func (suite *NetworkSuite) TestNew() {
	s := New(&mockRepo{})
	suite.NotNil(s)
}

func (suite *NetworkSuite) TestCreate() {
	cases := []struct {
		input *models.Network
		err   error
	}{
		{input: &models.Network{}, err: nil},
		{input: &models.Network{Network: "error"}, err: errors.New("Create Error")},
	}

	for _, c := range cases {
		repo := new(mockRepo)
		repo.On("CreateNetwork", c.input).Return(c.err)
		ns := service{repo: repo}
		err := ns.Create(c.input)
		suite.Equal(err, c.err)
		repo.AssertExpectations(suite.T())
	}
}

func (suite *NetworkSuite) TestGet() {
	cases := []struct {
		network     models.Network
		err         error
		expected    *models.Network
		expectedErr error
	}{
		{
			network:     models.Network{},
			err:         nil,
			expected:    &models.Network{},
			expectedErr: nil,
		},
		{
			network:     models.Network{},
			err:         errors.New("get error"),
			expected:    nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, c := range cases {
		repo := new(mockRepo)
		repo.On("GetNetwork", &models.Network{}).Return(c.err)
		ns := service{repo: repo}
		net, err := ns.Get()

		suite.Equal(c.expectedErr, err)
		suite.Equal(c.expected, net)
		repo.AssertExpectations(suite.T())
	}
}

func (suite NetworkSuite) TestSortSitesByName() {
	sites := []models.Site{
		{ID: 1, FriendlyName: "z"},
		{ID: 2, FriendlyName: "a"},
	}

	sortSitesByName(sites)

	suite.Equal(uint(2), sites[0].ID)
	suite.Equal("a", sites[0].FriendlyName)
	suite.Equal(uint(1), sites[1].ID)
	suite.Equal("z", sites[1].FriendlyName)
}

func TestNetworkSuite(t *testing.T) {
	suite.Run(t, new(NetworkSuite))
}
