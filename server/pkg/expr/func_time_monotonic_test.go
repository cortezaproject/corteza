package expr

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TL;DR: now() returns a wall-clock time, with no monotonic reading attached.
// Example: an automation writes now() into a record. Everything downstream
// stringifies it, and time.Time.String() appends the monotonic reading as a
// trailing "m=±<seconds>" field — which is process-local noise that no date
// parser reads back, so it must not be there to begin with.
func TestNow_CarriesNoMonotonicReading(t *testing.T) {
	got := now()

	require.NotContains(t, fmt.Sprintf("%v", got), " m=")
	require.NotContains(t, got.String(), " m=")

	// The wall clock itself must be untouched.
	require.WithinDuration(t, time.Now(), got, time.Minute)

	// Round(0) on an already-stripped time is a no-op: what comes back is
	// exactly what a caller comparing against the wall clock would expect.
	require.True(t, got.Equal(got.Round(0)))
}
