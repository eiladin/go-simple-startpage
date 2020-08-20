package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/helpers"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/stretchr/testify/assert"
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
		cfg := &config.Config{
			Database: config.Database{
				Driver:   c.Driver,
				Name:     c.Dbname,
				Username: c.Username,
				Password: c.Password,
				Host:     c.Host,
				Port:     c.Port,
			},
		}

		dsn := getDSN(cfg)
		assert.Equal(t, c.Expected, dsn.Name())
	}
}

func TestOpenErr(t *testing.T) {
	os.Setenv("GSS_DATABASE_DRIVER", "postgres")
	config.InitConfig("1.2.3", "./not-found.yaml")
	var b bytes.Buffer
	origFatalf := helpers.Fatalf
	helpers.Fatalf = func(format string, args ...interface{}) {
		fmt.Fprintf(&b, format, args)
	}
	defer func() { helpers.Fatalf = origFatalf }()
	InitDB()
	assert.Contains(t, b.String(), "failed to connect to")
	os.Unsetenv("GSS_DATABASE_DRIVER")
}

func TestDBFunctions(t *testing.T) {
	os.Setenv("GSS_DATABASE_DRIVER", "sqlite")
	os.Setenv("GSS_DATABASE_NAME", ":memory:")
	config.InitConfig("1.2.3", "./not-found.yaml")
	conn := InitDB()
	MigrateDB(conn)
	db := DB{DB: conn}

	net := model.Network{
		Network: "test",
		Links: []model.Link{
			{Name: "test-link-1"},
			{Name: "test-link-2"},
		},
		Sites: []model.Site{
			{FriendlyName: "test-site-1"},
			{FriendlyName: "test-site-2"},
		},
	}
	db.CreateNetwork(&net)
	// CreateNetwork assertions
	assert.Equal(t, uint(1), net.ID)
	assert.Equal(t, uint(1), net.Sites[0].ID)
	assert.Equal(t, uint(2), net.Sites[1].ID)
	assert.Equal(t, uint(1), net.Links[0].ID)
	assert.Equal(t, uint(2), net.Links[1].ID)

	findNet := model.Network{
		ID: 1,
	}
	db.GetNetwork(&findNet)
	// GetNetwork assertions
	assert.Equal(t, "test", findNet.Network)
	assert.Equal(t, "test-site-1", findNet.Sites[0].FriendlyName)

	findSite := model.Site{
		ID: 1,
	}
	db.GetSite(&findSite)
	// GetSite assertions
	assert.Equal(t, "test-site-1", findSite.FriendlyName)

	os.Unsetenv("GSS_DATABASE_DRIVER")
	os.Unsetenv("GSS_DATABASE_NAME")
}
