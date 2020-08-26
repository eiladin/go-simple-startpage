package database

import (
	"errors"
	"os"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetDSN(t *testing.T) {
	cases := []struct {
		Driver   string
		Dbname   string
		Username string
		Password string
		Host     string
		Port     string
		Expected string
	}{
		{Driver: "sqlite", Dbname: "test.db", Username: "testuser", Password: "testpass", Host: "localhost", Port: "1234", Expected: "sqlite"},
		{Driver: "postgres", Dbname: "test.db", Username: "testuser", Password: "testpass", Host: "localhost", Port: "1234", Expected: "postgres"},
		{Driver: "mysql", Dbname: "test.db", Username: "testuser", Password: "testpass", Host: "localhost", Port: "1234", Expected: "mysql"},
		{Driver: "notfound", Dbname: "test.db", Username: "testuser", Password: "testpass", Host: "localhost", Port: "1234", Expected: "sqlite"},
	}

	for _, c := range cases {
		cfg := &models.Config{
			Database: models.Database{
				Driver:   c.Driver,
				Name:     c.Dbname,
				Username: c.Username,
				Password: c.Password,
				Host:     c.Host,
				Port:     c.Port,
			},
		}

		dsn := getDSN(cfg)
		assert.Equal(t, c.Expected, dsn.Name(), "DSN Name should be %s", c.Expected)
	}
}

func TestOpenError(t *testing.T) {
	os.Setenv("GSS_DATABASE_DRIVER", "postgres")
	config.New("1.2.3", "./not-found.yaml")
	_, err := DB{}.New()
	assert.Contains(t, err.Error(), connectionRefusedErr(""), "A connectionRefusedError should be raised")
	os.Unsetenv("GSS_DATABASE_DRIVER")
}

func TestHandleError(t *testing.T) {
	cases := []struct {
		Err      error
		Expected error
	}{
		{Err: errors.New("unknown error"), Expected: errors.New("unknown error")},
		{Err: gorm.ErrRecordNotFound, Expected: store.ErrNotFound},
	}

	for _, c := range cases {
		err := handleError(c.Err)
		assert.EqualError(t, err, c.Expected.Error())
	}
}

func TestDBFunctions(t *testing.T) {
	os.Setenv("GSS_DATABASE_DRIVER", "sqlite")
	os.Setenv("GSS_DATABASE_NAME", ":memory:")
	config.New("1.2.3", "./not-found.yaml")
	db, err := DB{}.New()
	assert.NoError(t, err)

	net := models.Network{
		Network: "test",
		Links: []models.Link{
			{Name: "test-link-1"},
			{Name: "test-link-2"},
		},
		Sites: []models.Site{
			{FriendlyName: "test-site-1"},
			{FriendlyName: "test-site-2"},
		},
	}
	assert.NoError(t, db.CreateNetwork(&net))
	// CreateNetwork assertions
	assert.Equal(t, uint(1), net.ID, "Network ID should be '1'")
	assert.Equal(t, uint(1), net.Sites[0].ID, "Site ID should be '1'")
	assert.Equal(t, uint(2), net.Sites[1].ID, "Site ID should be '2'")
	assert.Equal(t, uint(1), net.Links[0].ID, "Link ID should be '1'")
	assert.Equal(t, uint(2), net.Links[1].ID, "Link ID should be '2'")

	findNet := models.Network{
		ID: 1,
	}
	assert.NoError(t, db.GetNetwork(&findNet))
	// GetNetwork assertions
	assert.Equal(t, "test", findNet.Network, "Network should be 'test'")
	assert.Equal(t, "test-site-1", findNet.Sites[0].FriendlyName, "Site FriendlyName should be 'test-site-1'")

	findSite := models.Site{
		ID: 1,
	}
	assert.NoError(t, db.GetSite(&findSite))
	// GetSite assertions
	assert.Equal(t, "test-site-1", findSite.FriendlyName, "Site FriendlyName should be 'test-site-1'")

	os.Unsetenv("GSS_DATABASE_DRIVER")
	os.Unsetenv("GSS_DATABASE_NAME")
}
