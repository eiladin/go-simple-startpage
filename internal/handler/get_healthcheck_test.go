package handler

import (
	"context"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
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
		h := handler{Config: &models.Config{
			Database: models.Database{
				Driver: c.Driver,
				Name:   c.Name,
			},
		}}
		err := h.checkDB(context.TODO())
		if c.Error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestAddGetHealthcheckRoute(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := handler{}
	h.AddGetHealthcheckRoute(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	assert.Contains(t, e, "GET /api/healthz")
}
