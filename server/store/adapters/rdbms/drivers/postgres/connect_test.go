package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	var (
		req = require.New(t)
	)

	c, err := NewConfig("postgres://uid:pwd@localhost/dbname?sslmode=disable")
	req.NoError(err)
	req.Equal("postgres", c.DriverName)
	req.Equal("dbname", c.DBName)
	req.Equal("postgres://uid:pwd@localhost/dbname?sslmode=disable", c.DataSourceName)
	req.Equal(256, c.MaxOpenConns)
	req.Equal(32, c.MaxIdleConns)

	// store params are applied to the pool and kept away from the driver
	c, err = NewConfig("postgres://uid:pwd@localhost/dbname?sslmode=disable&*connMaxOpen=20&*connMaxIdle=5")
	req.NoError(err)
	req.Equal("postgres://uid:pwd@localhost/dbname?sslmode=disable", c.DataSourceName)
	req.Equal(20, c.MaxOpenConns)
	req.Equal(5, c.MaxIdleConns)

	_, err = NewConfig("postgres://uid:pwd@localhost/dbname?*connMaxOpen=many")
	req.Error(err)

	_, err = NewConfig("mysql://uid:pwd@localhost/dbname")
	req.Error(err)
}
