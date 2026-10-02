package mysql

import (
	"testing"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms"
	"github.com/stretchr/testify/require"
)

func TestProcDataSourceName(t *testing.T) {
	var (
		req = require.New(t)
		c   *rdbms.ConnConfig
		err error
	)

	c, err = NewConfig("mysql://uid:@/dbname?parseTime=true")
	req.NoError(err)
	req.Contains(c.DataSourceName, "parseTime=true")
	req.Equal(c.DBName, "dbname")

	c, err = NewConfig("mysql+foo://uid:@/dbname")
	req.NoError(err)
	req.Contains(c.DataSourceName, "parseTime=true")
	req.Equal(c.DriverName, "mysql+foo")
	req.Equal(c.DBName, "dbname")

	// store params are applied to the pool and kept away from the driver
	c, err = NewConfig("mysql://uid:@/dbname?parseTime=true&*connMaxOpen=20&*connMaxIdle=5")
	req.NoError(err)
	req.Contains(c.DataSourceName, "parseTime=true")
	req.NotContains(c.DataSourceName, "connMaxOpen")
	req.Equal(20, c.MaxOpenConns)
	req.Equal(5, c.MaxIdleConns)

	_, err = NewConfig("mysql://uid:@/dbname?*connMaxOpen=many")
	req.Error(err)
}
