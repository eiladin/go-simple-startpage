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
	Get(uint) (*models.Status, error)
}

type repository interface {
	GetSite(*models.Site) error
}

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

func (c *service) Get(id uint) (*models.Status, error) {
	site := models.Site{ID: id}

	if err := c.repo.GetSite(&site); err != nil {
		return nil, ErrNotFound
	}

	res := models.NewStatus(c.config.Timeout, &site)
	return &res, nil
}
