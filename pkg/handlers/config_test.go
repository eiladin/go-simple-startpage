package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockConfigUseCase struct {
	mock.Mock
}

func (m *mockConfigUseCase) Get() (*config.Config, error) {
	args := m.Called()
	data := args.Get(0).(config.Config)
	return &data, args.Error(1)
}

type ConfigSuite struct {
	suite.Suite
}

func (suite ConfigSuite) TestGet() {
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
		uc := new(mockConfigUseCase)
		uc.On("Get").Return(config.Config{}, c.throwErr)

		cs := ConfigHandler{ConfigUseCase: uc}
		err := cs.Get(ctx)
		if c.wantErr == nil {
			suite.NoError(err)
		} else {
			suite.EqualError(err, c.wantErr.Error())
		}
		uc.AssertExpectations(suite.T())
	}
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
