package postgres

import (
	"testing"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/stretchr/testify/require"
)

func TestProcDataSourceName(t *testing.T) {
	var (
		req = require.New(t)
		c   *rdbms.Config
		err error
	)

	c, err = ProcDataSourceName("postgres://uid:@/dbname")
	req.NoError(err)
	req.Equal(c.DBName, "dbname")

	c, err = ProcDataSourceName("postgres+foo://uid:@/dbname")
	req.NoError(err)
	req.Equal(c.DriverName, "postgres+foo")
	req.Equal(c.DBName, "dbname")
}
