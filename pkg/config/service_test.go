package config

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Get() (*Config, error) {
	args := m.Called()
	data := args.Get(0).(Config)
	return &data, args.Error(1)
}

type ServiceSuite struct {
	suite.Suite
}

func (suite ServiceSuite) TestGet() {
	cases := []struct {
		throwErr error
		wantErr  error
	}{
		{throwErr: nil, wantErr: nil},
		{throwErr: errors.New("unhandled error"), wantErr: echo.ErrInternalServerError},
	}
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	for _, c := range cases {
		handler := new(mockHandler)
		handler.On("Get").Return(Config{}, c.throwErr)

		cs := Service{Handler: handler}
		err := cs.Get(ctx)
		if c.wantErr == nil {
			suite.NoError(err)
		} else {
			suite.EqualError(err, c.wantErr.Error())
		}
		handler.AssertExpectations(suite.T())
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
