package healthcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Check() http.Handler {
	m.Called()
	return nil
}

type ServiceSuite struct {
	suite.Suite
}

func (suite ServiceSuite) TestGet() {
	handler := new(mockHandler)

	handler.On("Check").Return(nil)
	hs := Service{Handler: handler}
	h := hs.Get()

	suite.NotNil(h)
	handler.AssertCalled(suite.T(), "Check")
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
