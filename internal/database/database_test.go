package database

import (
	"errors"
	"os"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseSuite struct {
	suite.Suite
}

func (suite *DatabaseSuite) TestMigrationError() {
	cfg := &config.Database{
		Driver: "fake",
	}
	store, err := New(cfg)
	defer os.Remove("simple-startpage.db")
	suite.Nil(store)
	suite.Contains(err.Error(), migrationFailedErr(""), "A migration failed error should be raised")
}

func (suite *DatabaseSuite) TestNetworkNotFoundError() {
	cfg := &config.Database{
		Driver: "sqlite",
		Name:   ":memory:",
	}
	db, err := New(cfg)
	suite.NoError(err)
	net := network.Network{}
	err = db.GetNetwork(&net)
	suite.EqualError(err, store.ErrNotFound.Error())
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
		cfg := &config.Database{
			Driver:   c.Driver,
			Name:     c.Dbname,
			Username: c.Username,
			Password: c.Password,
			Host:     c.Host,
			Port:     c.Port,
		}

		dsn := getDSN(cfg)
		suite.Equal(c.Expected, dsn.Name(), "DSN Name should be %s", c.Expected)
	}
}

func (suite *DatabaseSuite) TestOpenError() {
	c := config.Database{
		Driver: "postgres",
	}
	_, err := New(&c)
	suite.Contains(err.Error(), connectionRefusedErr(""), "A connectionRefusedError should be raised")
}

func (suite *DatabaseSuite) TestMigrateDB() {
	conn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)
	err = migrateDB(conn, &config.Database{Driver: "sqlite"})
	suite.NoError(err)
}

func (suite *DatabaseSuite) TestHandleError() {
	cases := []struct {
		Err      error
		Expected error
	}{
		{Err: errors.New("unknown error"), Expected: errors.New("unknown error")},
		{Err: gorm.ErrRecordNotFound, Expected: store.ErrNotFound},
	}

	for _, c := range cases {
		err := handleError(c.Err)
		suite.EqualError(err, c.Expected.Error())
	}
}

func (suite *DatabaseSuite) TestGetGormConfig() {
	getGormConfig(&config.Database{
		Log: true,
	})
}

func (suite *DatabaseSuite) TestPing() {
	c := config.Database{
		Driver: "sqlite",
		Name:   ":memory:",
	}
	db, err := New(&c)
	suite.NoError(err)
	suite.NoError(db.Ping())
}

func (suite *DatabaseSuite) TestDBFunctions() {
	c := config.Database{
		Driver: "sqlite",
		Name:   ":memory:",
	}
	db, err := New(&c)
	suite.NoError(err)

	net := network.Network{
		Network: "test",
		Links: []network.Link{
			{Name: "test-link-1"},
			{Name: "test-link-2"},
		},
		Sites: []network.Site{
			{
				Name: "test-site-1",
				DBTags: []network.DBTag{
					{Value: "tag-1"},
				},
				Tags: []string{"tag-2"},
			},
			{
				Name: "test-site-2",
				DBTags: []network.DBTag{
					{Value: "tag-3"},
					{Value: "tag-4"},
				},
			},
		},
	}
	suite.NoError(db.CreateNetwork(&net))
	// CreateNetwork assertions
	suite.Equal(uint(1), net.ID, "Network ID should be '1'")
	suite.Equal(uint(1), net.Sites[0].ID, "Site ID should be '1'")
	suite.Equal(uint(2), net.Sites[1].ID, "Site ID should be '2'")
	suite.Equal(uint(1), net.Links[0].ID, "Link ID should be '1'")
	suite.Equal(uint(2), net.Links[1].ID, "Link ID should be '2'")

	findNet := network.Network{ID: 1}
	suite.NoError(db.GetNetwork(&findNet))
	// GetNetwork assertions
	suite.Equal("test", findNet.Network, "Network should be 'test'")
	suite.Equal("test-site-1", findNet.Sites[0].Name, "Site Name should be 'test-site-1'")
	suite.Equal("test-link-1", findNet.Links[0].Name, "Link Name should be 'test-link-1'")
	suite.Equal(2, len(findNet.Sites[0].Tags), "Site 1 should have 2 tags")
	suite.Equal(2, len(findNet.Sites[1].Tags), "Site 2 should have 2 tags")
	suite.Equal("tag-1", findNet.Sites[0].Tags[0], "Site 1 should contain 'tag-1'")
	suite.Equal("tag-2", findNet.Sites[0].Tags[1], "Site 1 should contain 'tag-2'")
	suite.Equal("tag-3", findNet.Sites[1].Tags[0], "Site 2 should contain 'tag-3'")
	suite.Equal("tag-4", findNet.Sites[1].Tags[1], "Site 2 should contain 'tag-4'")

	findSite := network.Site{Name: "test-site-1"}
	suite.NoError(db.GetSite(&findSite))
	// GetSite assertions
	suite.Equal("test-site-1", findSite.Name, "Site Name should be 'test-site-1'")

	missingSite := network.Site{ID: 3}
	err = db.GetSite(&missingSite)
	suite.EqualError(err, store.ErrNotFound.Error())
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
