package network

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type mockNetworkService struct {
	CreateNetworkFunc func(*interfaces.Network)
	FindNetworkFunc   func(*interfaces.Network)
}

func (m *mockNetworkService) CreateNetwork(net *interfaces.Network) {
	m.CreateNetworkFunc(net)
}

func (m *mockNetworkService) FindNetwork(net *interfaces.Network) {
	m.FindNetworkFunc(net)
}

func TestNewNetwork(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	body := `{ "network": "test-network" }`

	ctx.Fasthttp.Request.Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Fasthttp.Request.SetBody([]byte(body))
	ctx.Fasthttp.Request.Header.SetContentLength(len(body))

	defer app.ReleaseCtx(ctx)
	var store mockNetworkService
	store.CreateNetworkFunc = func(net *interfaces.Network) {
		net.ID = 12345
	}
	handler := Handler{NetworkService: &store}
	handler.NewNetwork(ctx)

	assert.Equal(t, `{"id":12345}`, string(ctx.Fasthttp.Response.Body()))
}

func TestNewNetworkError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	body := `{ "network": "test-network" }`

	ctx.Fasthttp.Request.SetBody([]byte(body))
	ctx.Fasthttp.Request.Header.SetContentLength(len(body))

	defer app.ReleaseCtx(ctx)
	var store mockNetworkService
	store.CreateNetworkFunc = func(net *interfaces.Network) {
		net.ID = 12345
	}
	handler := Handler{NetworkService: &store}
	handler.NewNetwork(ctx)

	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Fasthttp.Response.StatusCode())
}

func TestFindNetwork(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	var store mockNetworkService
	store.FindNetworkFunc = func(net *interfaces.Network) {
		net.ID = 12345
		net.Network = "test-network"
	}
	handler := Handler{NetworkService: &store}
	handler.GetNetwork(ctx)

	assert.Equal(t, `{"network":"test-network","links":null,"sites":null}`, string(ctx.Fasthttp.Response.Body()))
}
