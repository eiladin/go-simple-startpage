package api

import (
	"context"
	"errors"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/suite"
)

type HealthcheckServiceSuite struct {
	suite.Suite
}

func (suite *HealthcheckServiceSuite) TestCheckDB() {
	cases := []struct {
		Driver string
		Name   string
		Error  bool
	}{
		{Driver: "postgres", Name: "name1", Error: true},
		{Driver: "sqlite", Name: ":memory:", Error: false},
	}

	for _, c := range cases {
		cfg := &models.Config{
			Database: models.Database{
				Driver: c.Driver,
				Name:   c.Name,
			},
		}

		var pingFunc func() error
		if c.Error {
			pingFunc = func() error { return errors.New("connection error") }
		} else {
			pingFunc = func() error { return nil }
		}

		hs := NewHealthcheckService(cfg, &mockStore{PingFunc: pingFunc})
		err := hs.checkDB(context.TODO())
		if c.Error {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *HealthcheckServiceSuite) TestRegister() {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	NewHealthcheckService(&models.Config{}, &mockStore{}).Register(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/healthz")
}

func TestHealthcheckSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckServiceSuite))
}
