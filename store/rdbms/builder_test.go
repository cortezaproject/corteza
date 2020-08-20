package rdbms

import (
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuilder(t *testing.T) {
	var (
		req = require.New(t)
	)

	upsert, err := UpsertBuilder("tbl", store.Payload{"c1": "v1", "c2": "v2"}, "c1")
	req.NoError(err)
	sql, args, err := upsert.ToSql()
	req.NoError(err)
	req.Contains(sql, "ON CONFLICT UPDATE SET")
	req.Contains(sql, "INSERT INTO tbl")
	req.Equal([]interface{}{"v1", "v2", "v2", "v1"}, args)

}
