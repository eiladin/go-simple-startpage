package healthcheck

import (
	"context"
	"errors"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStore struct {
	mock.Mock
	PingFunc func() error
}

func (m *mockStore) Ping() error                            { return m.PingFunc() }
func (m *mockStore) CreateNetwork(net *model.Network) error { return nil }
func (m *mockStore) GetNetwork(net *model.Network) error    { return nil }
func (m *mockStore) GetSite(site *model.Site) error         { return nil }

type HealthcheckServiceSuite struct {
	suite.Suite
}

func (suite *HealthcheckServiceSuite) TestCheckDB() {
	cases := []struct {
		Database model.Database
		Error    error
	}{
		{Database: model.Database{Driver: "postgres", Name: "name1"}, Error: errors.New("connection error")},
		{Database: model.Database{Driver: "sqlite", Name: ":memory:"}, Error: nil},
	}

	store := &mockStore{}
	cfg := &model.Config{}
	hs := NewHealthcheckService(cfg, store)

	for _, c := range cases {
		cfg.Database = c.Database
		store.PingFunc = func() error { return c.Error }

		err := hs.checkDB(context.TODO())
		if c.Error != nil {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *HealthcheckServiceSuite) TestRegister() {
	app := echo.New()
	NewHealthcheckService(&model.Config{}, &mockStore{}).Register(app)
	e := []string{}
	for _, r := range app.Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/healthz")
}

func TestHealthcheckSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckServiceSuite))
}
