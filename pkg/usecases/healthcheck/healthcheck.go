package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
)

type IHealthcheck interface {
	Check() http.Handler
}

type repository interface {
	Ping() error
}

// Compile-time proof of interface implementation.
var _ IHealthcheck = (*service)(nil)

type service struct {
	repo repository
}

func New(repo repository) IHealthcheck {
	return &service{
		repo: repo,
	}
}

func (c *service) checkDB(ctx context.Context) error {
	if err := c.repo.Ping(); err != nil {
		return fmt.Errorf("unable to connect to database %w", err)
	}
	return nil
}

func (c *service) Check() http.Handler {
	return healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(c.checkDB)),
	)

}
