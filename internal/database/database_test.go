package database

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/config"
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
		cfg := &config.Configuration{
			Database: config.DatabaseConfiguration{
				Driver:   c.Driver,
				Dbname:   c.Dbname,
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
