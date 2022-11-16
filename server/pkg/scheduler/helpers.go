package scheduler

import (
	"time"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/getsentry/sentry-go"
	"github.com/gorhill/cronexpr"
	"go.uber.org/zap"
)

// OnInterval parses all given strings as crontab expressions (ii) and returns true if any of them matches current time
func OnInterval(ii ...string) bool {
	match, err := onInterval(now(), ii...)
	if err != nil {
		logger.Default().Error("failed to parse interval value", zap.Strings("value", ii), zap.Error(err))
		sentry.CaptureException(err)
	}
	return match
}

func onInterval(now time.Time, ii ...string) (bool, error) {
	var (
		// This function will likely to be called exactly on a minute (00 sec) or a few milliseconds after
		// Round it up to the smallest unit that cronexpr package supports
		currTime = now.Truncate(time.Second)

		// For cron expression reference we need to subtract 1ns
		// this will cause Next() fn to include next nanosecond (if it matches)
		cronRef = currTime.Add(-time.Nanosecond)
	)

	// At least one of the given expressions should match
	for _, i := range ii {
		if len(i) == 0 {
			// skip empty values
			continue
		}

		exp, err := cronexpr.Parse(i)
		if err != nil {
			return false, err
		}

		return currTime.Equal(exp.Next(cronRef)), nil
	}

	return false, nil
}

// OnTimestamp parses all given strings as RFC3339 timestamps and returns true if any of them matches current time
func OnTimestamp(tt ...string) bool {
	match, err := onTimestamp(now(), tt...)
	if err != nil {
		logger.Default().Error("failed to parse timestamp value", zap.Strings("value", tt), zap.Error(err))
		sentry.CaptureException(err)
	}
	return match
}

func onTimestamp(now time.Time, tt ...string) (bool, error) {
	var (
		// This function will likely to be called exactly on a minute (00 sec) or a few milliseconds after
		// Round it up to the smallest unit that cronexpr package supports
		currTime = now.Truncate(time.Second)
	)

	for _, t := range tt {
		if len(t) == 0 {
			// skip empty values
			continue
		}

		ts, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return false, err
		}

		return currTime.Equal(ts.Round(time.Second)), nil
	}

	return false, nil
}
