package config

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockUseCase struct {
	mock.Mock
}

func (m *mockUseCase) Get() (*Config, error) {
	args := m.Called()
	data := args.Get(0).(Config)
	return &data, args.Error(1)
}

type HandlerSuite struct {
	suite.Suite
}

func (suite HandlerSuite) TestGet() {
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
		uc := new(mockUseCase)
		uc.On("Get").Return(Config{}, c.throwErr)

		cs := Handler{UseCase: uc}
		err := cs.Get(ctx)
		if c.wantErr == nil {
			suite.NoError(err)
		} else {
			suite.EqualError(err, c.wantErr.Error())
		}
		uc.AssertExpectations(suite.T())
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
