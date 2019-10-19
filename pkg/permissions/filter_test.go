package permissions

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestResourceFilter_Build(t *testing.T) {
	rf := ResourceFilter{
		dbTable:   "ptbl",
		pkColName: "pkcol",
		resource:  "res:",
		operation: "read",
		chk:       nil,
	}

	req := require.New(t)

	rf.fallback = Allow
	req.Equal(
		`COALESCE((SELECT access = 1 FROM ptbl WHERE operation = 'read' AND resource = CONCAT('res:', pkcol) AND rel_role IN ('1') ORDER BY access LIMIT 1), TRUE)`,
		squirrel.DebugSqlizer(rf),
	)

	rf.roles = []uint64{123}
	req.Equal(
		`COALESCE((SELECT access = 1 FROM ptbl WHERE operation = 'read' AND resource = CONCAT('res:', pkcol) AND rel_role IN ('123') ORDER BY access LIMIT 1), (SELECT access = 1 FROM ptbl WHERE operation = 'read' AND resource = CONCAT('res:', pkcol) AND rel_role IN ('1') ORDER BY access LIMIT 1), TRUE)`,
		squirrel.DebugSqlizer(rf),
	)

	rf.chk = &ServiceDenyAll{}
	req.Equal(
		`COALESCE((SELECT access = 1 FROM ptbl WHERE operation = 'read' AND resource = CONCAT('res:', pkcol) AND rel_role IN ('123') ORDER BY access LIMIT 1), FALSE)`,
		squirrel.DebugSqlizer(rf),
	)

	rf.chk = &ServiceAllowAll{}
	req.Equal(
		`COALESCE((SELECT access = 1 FROM ptbl WHERE operation = 'read' AND resource = CONCAT('res:', pkcol) AND rel_role IN ('123') ORDER BY access LIMIT 1), TRUE)`,
		squirrel.DebugSqlizer(rf),
	)

	rf.superuser = true
	req.Equal(
		`TRUE`,
		squirrel.DebugSqlizer(rf),
	)

}
