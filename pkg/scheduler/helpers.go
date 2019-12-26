package scheduler

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorhill/cronexpr"
)

// OnInterval parses all given strings as crontab expressions (ii) and returns true if any of them matches current time
func OnInterval(ee ...string) bool {
	var (
		// This function will likely to be called exactly on a minute (00 sec) or a few milliseconds after
		// Round it up to the smallest unit that cronexpr package supports
		currTime = now().Truncate(time.Second)

		// For cron expression reference we need to subtract 1ns
		// this will cause Next() fn to include next nanosecond (if it matches)
		cronRef = currTime.Add(-time.Nanosecond)
	)

	// At least one of the given expressions should match
	for _, e := range ee {
		exp, err := cronexpr.Parse(e)
		if err != nil {
			sentry.CaptureException(err)
			return false
		}

		if currTime.Equal(exp.Next(cronRef)) {
			return true
		}
	}

	return false
}

// OnTimestamp parses all given strings as RFC3339 timestamps and returns true if any of them matches current time
func OnTimestamp(tt ...string) bool {
	var (
		// This function will likely to be called exactly on a minute (00 sec) or a few milliseconds after
		// Round it up to the smallest unit that cronexpr package supports
		currTime = now().Truncate(time.Second)
	)

	for _, t := range tt {
		ts, err := time.Parse(time.RFC3339, t)

		if err != nil {
			sentry.CaptureException(err)
			return false
		}

		if currTime.Equal(ts.Round(time.Second)) {
			return true
		}
	}

	return false
}
