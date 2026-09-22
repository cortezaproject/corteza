package corredor

import (
	"sync"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestService_Status(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		svc := &service{opt: options.CorredorOpt{Enabled: false}, sScriptsL: &sync.Mutex{}}
		enabled, connected, refreshedAt := svc.Status()
		require.False(t, enabled)
		require.False(t, connected)
		require.Nil(t, refreshedAt)
	})

	t.Run("enabled without a connection", func(t *testing.T) {
		svc := &service{opt: options.CorredorOpt{Enabled: true}, sScriptsL: &sync.Mutex{}}
		enabled, connected, refreshedAt := svc.Status()
		require.True(t, enabled)
		require.False(t, connected)
		require.Nil(t, refreshedAt)
	})

	t.Run("reports the last fetch", func(t *testing.T) {
		ts := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
		svc := &service{opt: options.CorredorOpt{Enabled: true}, sScriptsL: &sync.Mutex{}, sScriptsTS: ts}
		_, _, refreshedAt := svc.Status()
		require.NotNil(t, refreshedAt)
		require.Equal(t, ts, *refreshedAt)
	})
}
