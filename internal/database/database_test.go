package database

import (
	"errors"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/httperror"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseSuite struct {
	suite.Suite
}

func (suite *DatabaseSuite) TestGetDSN() {
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
		suite.Equal(c.Expected, dsn.Name(), "DSN Name should be %s", c.Expected)
	}
}

func (suite *DatabaseSuite) TestOpenError() {
	c := config.Config{
		Database: config.Database{
			Driver: "postgres",
		},
	}
	_, err := New(&c)
	suite.Contains(err.Error(), connectionRefusedErr(""), "A connectionRefusedError should be raised")
}

func (suite *DatabaseSuite) TestMigrateDB() {
	conn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)
	err = migrateDB(conn)
	suite.NoError(err)
}

func (suite *DatabaseSuite) TestHandleError() {
	cases := []struct {
		Err      error
		Expected error
	}{
		{Err: errors.New("unknown error"), Expected: errors.New("unknown error")},
		{Err: gorm.ErrRecordNotFound, Expected: httperror.ErrNotFound},
	}

	for _, c := range cases {
		err := handleError(c.Err)
		suite.EqualError(err, c.Expected.Error())
	}
}

func (suite *DatabaseSuite) TestPing() {
	c := config.Config{
		Database: config.Database{
			Driver: "sqlite",
			Name:   ":memory:",
		},
	}
	db, err := New(&c)
	suite.NoError(err)
	suite.NoError(db.Ping())
}

func (suite *DatabaseSuite) TestDBFunctions() {
	c := config.Config{
		Database: config.Database{
			Driver: "sqlite",
			Name:   ":memory:",
		},
	}
	db, err := New(&c)
	suite.NoError(err)

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
	suite.NoError(db.CreateNetwork(&net))
	// CreateNetwork assertions
	suite.Equal(uint(1), net.ID, "Network ID should be '1'")
	suite.Equal(uint(1), net.Sites[0].ID, "Site ID should be '1'")
	suite.Equal(uint(2), net.Sites[1].ID, "Site ID should be '2'")
	suite.Equal(uint(1), net.Links[0].ID, "Link ID should be '1'")
	suite.Equal(uint(2), net.Links[1].ID, "Link ID should be '2'")

	findNet := models.Network{ID: 1}
	suite.NoError(db.GetNetwork(&findNet))
	// GetNetwork assertions
	suite.Equal("test", findNet.Network, "Network should be 'test'")
	suite.Equal("test-site-1", findNet.Sites[0].FriendlyName, "Site FriendlyName should be 'test-site-1'")
	suite.Equal("test-link-1", findNet.Links[0].Name, "Link Name should be 'test-link-1'")

	findSite := models.Site{ID: 1}
	suite.NoError(db.GetSite(&findSite))
	// GetSite assertions
	suite.Equal("test-site-1", findSite.FriendlyName, "Site FriendlyName should be 'test-site-1'")
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
