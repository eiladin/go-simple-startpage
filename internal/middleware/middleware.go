package middleware

import (
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func loggerConfig(production bool) middleware.LoggerConfig {
	loggerConfig := middleware.DefaultLoggerConfig
	if !production {
		loggerConfig = middleware.LoggerConfig{
			Output: os.Stdout,
			Format: "method=${method}, uri=${uri}, status=${status} ${error}\n",
		}
	}
	return loggerConfig
}

func csrfSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().Host, "localhost")
}

func gzipSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().URL.Path, "/swagger")
}

func staticSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().URL.Path, "/api") ||
		strings.Contains(ctx.Request().URL.Path, "/swagger")
}

func GetMiddleware(c *model.Config) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.LoggerWithConfig(loggerConfig(c.IsProduction())),
		middleware.CORS(),
		middleware.RequestID(),
		middleware.Secure(),
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			Skipper:      csrfSkipper,
			CookieSecure: true,
		}),
		middleware.Recover(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: gzipSkipper,
		}),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper: staticSkipper,
			Index:   "index.html",
			Root:    "ui/dist",
			Browse:  false,
			HTML5:   true,
		}),
	}
}
