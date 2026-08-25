package values

import (
	"fmt"
	"testing"
	"time"

	"github.com/crusttech/human/server/compose/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

// TL;DR: a time taken from time.Now() reads as a date, however it was
// stringified on the way in.
// Example: a TAQ stamps a DateTime field with now(). The value is cast to a
// string map before it reaches the record, and time.Time.String() appends the
// monotonic clock reading — so the field was refused as "not a date/time".
func TestSanitizeDatetime_MonotonicReading(t *testing.T) {
	// A time carrying a monotonic reading, exactly what time.Now() returns.
	mono := time.Now()
	require.Contains(t, fmt.Sprintf("%v", mono), " m=", "test needs a monotonic reading to be meaningful")

	viaCast, err := cast.ToStringE(mono)
	require.NoError(t, err)

	for name, in := range map[string]interface{}{
		"time.Time":           mono,
		"pointer":             &mono,
		"already stringified": viaCast,
	} {
		t.Run(name, func(t *testing.T) {
			require.NotEmpty(t, sDatetime(in, false, false), "should read as a date/time")
		})
	}

	t.Run("value matches the wall clock", func(t *testing.T) {
		want := mono.Round(0).UTC().Format(datetimeInternalFormatFull)
		require.Equal(t, want, sDatetime(mono, false, false))
		require.Equal(t, want, sDatetime(viaCast, false, false))
	})

	// sanitizeStrict is what turns an unreadable date into a save-time error,
	// so it is the surface a person or a TAQ actually hits.
	t.Run("sanitizeStrict accepts it", func(t *testing.T) {
		f := &types.ModuleField{Kind: "DateTime"}
		out, err := sanitizeStrict(f, viaCast)
		require.NoError(t, err)
		require.NotEmpty(t, out)
	})

	t.Run("onlyDate still narrows to the date", func(t *testing.T) {
		f := &types.ModuleField{Kind: "DateTime"}
		f.Options = types.ModuleFieldOptions{"onlyDate": true}
		out, err := sanitizeStrict(f, viaCast)
		require.NoError(t, err)
		require.Equal(t, mono.Round(0).UTC().Format(datetimeInternalFormatDate), out)
	})
}

// TL;DR: a genuinely unreadable value is still refused.
// Example: the monotonic strip must not turn the sanitizer into one that
// accepts anything.
func TestSanitizeDatetime_StillRejectsNonsense(t *testing.T) {
	f := &types.ModuleField{Kind: "DateTime"}
	for _, in := range []string{"not a date", "m=+1.5", "2026-13-45"} {
		_, err := sanitizeStrict(f, in)
		require.Error(t, err, "%q must not read as a date/time", in)
	}
}
