package api

import (
	"context"
	"errors"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckDB(t *testing.T) {
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

		h := handler{Config: cfg, Store: mockStore{PingFunc: pingFunc}}
		err := h.checkDB(context.TODO())
		if c.Error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestAddHealthCheckRoutes(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := handler{ApiRoot: app}
	h.addHeathcheckRoutes()
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	assert.Contains(t, e, "GET /api/healthz")
}
