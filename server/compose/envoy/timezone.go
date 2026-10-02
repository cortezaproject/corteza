package envoy

import (
	"fmt"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/spf13/cast"
)

// resolveTimezone reads the "timezone" param and returns the loaded location.
// Returns a nil location when no timezone is set. An unresolvable timezone is
// reported as an error so the export fails loudly instead of silently skipping
// the conversion.
func resolveTimezone(p envoyx.EncodeParams) (*time.Location, error) {
	tz := cast.ToString(p.Params["timezone"])
	if tz == "" {
		return nil, nil
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("invalid export timezone %q: %w", tz, err)
	}
	return loc, nil
}

// formatInTimezone reparses v as a datetime and reformats it in loc.
// On parse failure returns v unchanged so values are never dropped.
func formatInTimezone(v string, loc *time.Location) string {
	if loc == nil || v == "" {
		return v
	}

	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999 -0700 MST", "2006-01-02 15:04:05 -0700 MST", "2006-01-02T15:04:05Z"}
	for _, ly := range layouts {
		if t, err := time.Parse(ly, v); err == nil {
			return t.In(loc).Format(time.RFC3339)
		}
	}

	// Fall back to zone-less layouts; assume values are stored in UTC.
	naiveLayouts := []string{"2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05", "2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05"}
	for _, ly := range naiveLayouts {
		if t, err := time.ParseInLocation(ly, v, time.UTC); err == nil {
			return t.In(loc).Format(time.RFC3339)
		}
	}

	return v
}
