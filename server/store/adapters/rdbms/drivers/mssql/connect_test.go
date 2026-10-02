package mssql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_DBName(t *testing.T) {
	req := require.New(t)

	// database in the query param, instance in the path
	c, err := NewConfig("sqlserver://sa:secret@localhost/INSTANCENAME?database=corteza")
	req.NoError(err)
	req.Equal("corteza", c.DBName)
	req.Contains(c.DataSourceName, "localhost/INSTANCENAME")
	req.Contains(c.DataSourceName, "database=corteza")

	// database in the query param, no instance
	c, err = NewConfig("sqlserver://sa:secret@localhost:1433?database=corteza")
	req.NoError(err)
	req.Equal("corteza", c.DBName)

	// no database param; the path is all we have
	c, err = NewConfig("sqlserver://sa:secret@localhost/corteza")
	req.NoError(err)
	req.Equal("corteza", c.DBName)

	// nothing at all; lookups use the current database
	c, err = NewConfig("sqlserver://sa:secret@localhost")
	req.NoError(err)
	req.Equal("", c.DBName)
}

func TestDbObject(t *testing.T) {
	req := require.New(t)
	req.Equal("corteza.INFORMATION_SCHEMA.COLUMNS", dbObject("corteza", "INFORMATION_SCHEMA.COLUMNS"))
	req.Equal("INFORMATION_SCHEMA.COLUMNS", dbObject("", "INFORMATION_SCHEMA.COLUMNS"))
}
