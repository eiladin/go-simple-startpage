package status

import (
	"errors"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/models"
)

var (
	ErrNotFound = errors.New("Status not found")
)

type IStatus interface {
	Get(string) (*models.Status, error)
}

type repository interface {
	GetSite(*models.Site) error
}

// Compile-time proof of interface implementation.
var _ IStatus = (*service)(nil)

type service struct {
	repo   repository
	config *config.Config
}

func New(repo repository, cfg *config.Config) IStatus {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

func (c *service) Get(name string) (*models.Status, error) {
	site := models.Site{Name: name}

	if err := c.repo.GetSite(&site); err != nil {
		return nil, ErrNotFound
	}

	res := models.NewStatus(c.config.Timeout, &site)
	return &res, nil
}
