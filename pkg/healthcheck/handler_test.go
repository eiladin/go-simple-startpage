package healthcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockUseCase struct {
	mock.Mock
}

func (m *mockUseCase) Check() http.Handler {
	m.Called()
	return nil
}

type HandlerSuite struct {
	suite.Suite
}

func (suite HandlerSuite) TestGet() {
	uc := new(mockUseCase)

	uc.On("Check").Return(nil)
	hs := Handler{UseCase: uc}
	h := hs.Get()

	suite.NotNil(h)
	uc.AssertCalled(suite.T(), "Check")
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
