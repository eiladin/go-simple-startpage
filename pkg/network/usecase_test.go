package network

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateNetwork(net *Network) error {
	args := m.Called(net)
	return args.Error(0)
}

func (m *mockRepo) GetNetwork(net *Network) error {
	args := m.Called(net)
	return args.Error(0)
}

type UseCaseSuite struct {
	suite.Suite
}

func (suite *UseCaseSuite) TestNew() {
	s := New(&mockRepo{})
	suite.NotNil(s)
}

func (suite *UseCaseSuite) TestCreate() {
	cases := []struct {
		input *Network
		err   error
	}{
		{input: &Network{}, err: nil},
		{input: &Network{Network: "error"}, err: errors.New("Create Error")},
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

func (suite *UseCaseSuite) TestGet() {
	cases := []struct {
		network     Network
		err         error
		expected    *Network
		expectedErr error
	}{
		{
			network:     Network{},
			err:         nil,
			expected:    &Network{},
			expectedErr: nil,
		},
		{
			network:     Network{},
			err:         errors.New("get error"),
			expected:    nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, c := range cases {
		repo := new(mockRepo)
		repo.On("GetNetwork", &Network{}).Return(c.err)
		ns := service{repo: repo}
		net, err := ns.Get()

		suite.Equal(c.expectedErr, err)
		suite.Equal(c.expected, net)
		repo.AssertExpectations(suite.T())
	}
}

func (suite UseCaseSuite) TestSortSitesByName() {
	sites := []Site{
		{ID: 1, Name: "z"},
		{ID: 2, Name: "a"},
	}

	sortSitesByName(sites)

	suite.Equal(uint(2), sites[0].ID)
	suite.Equal("a", sites[0].Name)
	suite.Equal(uint(1), sites[1].ID)
	suite.Equal("z", sites[1].Name)
}

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
