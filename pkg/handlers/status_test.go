package handlers

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/status"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStatusUseCase struct {
	mock.Mock
}

func (m *mockStatusUseCase) Get(name string) (*models.Status, error) {
	args := m.Called(name)
	data := args.Get(0).(models.Status)
	return &data, args.Error(1)
}

type StatusSuite struct {
	suite.Suite
}

func (suite *StatusSuite) TestGet() {
	app := echo.New()

	cases := []struct {
		id       uint
		param    string
		uri      string
		isUp     bool
		throwErr error
		wantErr  error
	}{
		{id: 1, param: "test-site-1", uri: "https://my.test.site", isUp: true, throwErr: nil, wantErr: nil},
		{id: 1, param: "test-site-2", uri: "https://my.fail.site", isUp: false, throwErr: nil, wantErr: nil},
		{id: 1, param: "test-site-3", uri: "https://^^invalidurl^^", isUp: false, throwErr: nil, wantErr: nil},
		{id: 1, param: "test-site-4", uri: "ssh://localhost:22224", isUp: true, throwErr: nil, wantErr: nil},
		{id: 1, param: "test-site-5", uri: "ssh://localhost:1234", isUp: false, throwErr: nil, wantErr: nil},
		{id: 1, param: "test-site-6", uri: "https://500.test.site", isUp: false, throwErr: nil, wantErr: nil},
		{id: 1, param: "", uri: "https://no-id.test.site", isUp: false, throwErr: errors.New("bad request"), wantErr: echo.ErrBadRequest},
		{id: 12345, param: "test-site-9", uri: "https://bigid.test.site", isUp: false, throwErr: status.ErrNotFound, wantErr: echo.ErrNotFound},
		{id: 1, param: "test-site-10", uri: "https://error.test.site", isUp: false, throwErr: errors.New("internal server error"), wantErr: echo.ErrInternalServerError},
		{id: 1, param: "tste-site-11", uri: "https://timeout.test.site", isUp: false, throwErr: nil, wantErr: nil},
	}

	for _, c := range cases {
		uc := new(mockStatusUseCase)
		if !errors.Is(c.wantErr, echo.ErrBadRequest) {
			uc.On("Get", c.param).Return(models.Status{IsUp: c.isUp}, c.throwErr)
		}
		ss := StatusHandler{StatusUseCase: uc}

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := app.NewContext(req, rec)
		ctx.SetPath("/:name")
		ctx.SetParamNames("name")
		ctx.SetParamValues(c.param)
		err := ss.Get(ctx)
		uc.AssertExpectations(suite.T())
		if c.wantErr != nil {
			suite.EqualError(err, c.wantErr.Error(), "%s should return %s", c.uri, c.wantErr.Error())
		} else {
			dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
			res := models.Status{}
			err := dec.Decode(&res)
			suite.NoError(err)
			suite.Equal(c.isUp, res.IsUp, "%s isUp should be %t", c.uri, c.isUp)
		}
	}
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}
