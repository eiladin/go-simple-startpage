package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockHealthcheckUseCase struct {
	mock.Mock
}

func (m *mockHealthcheckUseCase) Check() http.Handler {
	m.Called()
	return nil
}

type HealthcheckSuite struct {
	suite.Suite
}

func (suite HealthcheckSuite) TestGet() {
	uc := new(mockHealthcheckUseCase)

	uc.On("Check").Return(nil)
	hs := HealthcheckHandler{HealthcheckUseCase: uc}
	h := hs.Get()

	suite.NotNil(h)
	uc.AssertCalled(suite.T(), "Check")
}

func TestHealthcheckSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckSuite))
}
