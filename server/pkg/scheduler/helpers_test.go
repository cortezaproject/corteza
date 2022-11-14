package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnInterval(t *testing.T) {
	cases := []struct {
		name  string
		now   string
		ii    []string
		match bool
		err   bool
	}{
		{
			"empty",
			time.Now().Format(time.RFC3339),
			[]string{},
			false,
			false,
		},
		{
			"nil",
			time.Now().Format(time.RFC3339),
			nil,
			false,
			false,
		},
		{
			"next minute",
			"2019-10-10T10:11:00Z",
			[]string{"* * * * *"},
			true,
			false,
		},
		{
			"every 5 minutes",
			"2019-10-10T10:15:00Z",
			[]string{"*/5 * * * *"},
			true,
			false,
		},
		{
			"not full minute",
			"2019-10-10T10:15:50Z",
			[]string{"* * * * *"},
			false,
			false,
		},
		{
			"invalid format",
			"2019-10-10T10:15:50Z",
			[]string{":P"},
			false,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, c.now)
			assert.NoError(t, err)

			match, err := onInterval(now, c.ii...)
			if !c.err {
				assert.NoError(t, err)
			}

			assert.Equal(t, c.match, match)
		})
	}

	// touch exported func
	OnInterval("* * * * *")
	OnInterval(":P")
}

func TestOnTimestamp(t *testing.T) {
	cases := []struct {
		name  string
		now   string
		ii    []string
		match bool
		err   bool
	}{
		{
			"empty",
			time.Now().Format(time.RFC3339),
			[]string{},
			false,
			false,
		},
		{
			"nil",
			time.Now().Format(time.RFC3339),
			nil,
			false,
			false,
		},
		{
			"must match",
			"2019-10-10T10:11:00Z",
			[]string{"2019-10-10T10:11:00Z"},
			true,
			false,
		},
		{
			"must not match",
			"2019-10-10T10:15:00Z",
			[]string{"2020-10-10T10:15:00Z"},
			false,
			false,
		},
		{
			"invalid format",
			"2019-10-10T10:15:00Z",
			[]string{":P"},
			false,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, c.now)
			assert.NoError(t, err)

			match, err := onTimestamp(now, c.ii...)
			if !c.err {
				assert.NoError(t, err)
			}

			assert.Equal(t, c.match, match)
		})
	}

	// touch exported func
	OnTimestamp("2019-10-10T10:15:00Z")
	OnTimestamp(":P")
}
