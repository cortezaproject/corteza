package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The config scans from bytes, from a string (SQL Server) and from NULL
func TestNotificationConfigScan(t *testing.T) {
	for _, src := range []any{[]byte(`{"recordID":"1"}`), `{"recordID":"1"}`, nil} {
		var nc NotificationConfig
		require.NoError(t, nc.Scan(src), "%T", src)
	}
}
