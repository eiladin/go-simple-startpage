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

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) TestNew() {
	s := NewHandler(&mockRepo{})
	suite.NotNil(s)
}

func (suite *HandlerSuite) TestCreate() {
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
		handler := handler{repo: repo}
		err := handler.Create(c.input)
		suite.Equal(err, c.err)
		repo.AssertExpectations(suite.T())
	}
}

func (suite *HandlerSuite) TestGet() {
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
		handler := handler{repo: repo}
		net, err := handler.Get()

		suite.Equal(c.expectedErr, err)
		suite.Equal(c.expected, net)
		repo.AssertExpectations(suite.T())
	}
}

func (suite HandlerSuite) TestSortSitesByName() {
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

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
