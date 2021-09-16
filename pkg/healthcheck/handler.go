package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
)

type IHandler interface {
	Check() http.Handler
}

type repository interface {
	Ping() error
}

// Compile-time proof of interface implementation.
var _ IHandler = (*handler)(nil)

type handler struct {
	repo repository
}

func NewHandler(repo repository) IHandler {
	return &handler{repo: repo}
}

func (c *handler) checkDB(ctx context.Context) error {
	if err := c.repo.Ping(); err != nil {
		return fmt.Errorf("unable to connect to database %w", err)
	}
	return nil
}

func (c *handler) Check() http.Handler {
	return healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(c.checkDB)),
	)
}
