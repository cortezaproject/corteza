package mysql

import (
	"github.com/cortezaproject/corteza-server/store/rdbms"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcDataSourceName(t *testing.T) {
	var (
		req = require.New(t)
		c   *rdbms.Config
		err error
	)

	c, err = ProcDataSourceName("mysql://uid:@/dbname?parseTime=true")
	req.NoError(err)
	req.Contains(c.DataSourceName, "parseTime=true")
	req.Equal(c.DBName, "dbname")

	c, err = ProcDataSourceName("mysql+foo://uid:@/dbname")
	req.NoError(err)
	req.Contains(c.DataSourceName, "parseTime=true")
	req.Equal(c.DriverName, "mysql+foo")
	req.Equal(c.DBName, "dbname")
}
