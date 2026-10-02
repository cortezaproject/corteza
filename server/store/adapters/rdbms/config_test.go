package rdbms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConnConfig_ParseExtra(t *testing.T) {
	tcc := []struct {
		name string
		dsn  string
		out  string
		err  string
		chk  func(*require.Assertions, *ConnConfig)
	}{
		{
			name: "no query string",
			dsn:  "postgres://u:p@h/db",
			out:  "postgres://u:p@h/db",
		},
		{
			name: "no store params",
			dsn:  "postgres://u:p@h/db?sslmode=disable",
			out:  "postgres://u:p@h/db?sslmode=disable",
		},
		{
			name: "pool params are applied and stripped",
			dsn:  "postgres://u:p@h/db?*connMaxOpen=20&sslmode=disable&*connMaxIdle=5&*connMaxLifetime=5m",
			out:  "postgres://u:p@h/db?sslmode=disable",
			chk: func(req *require.Assertions, c *ConnConfig) {
				req.Equal(20, c.MaxOpenConns)
				req.Equal(5, c.MaxIdleConns)
				req.Equal(5*time.Minute, c.ConnMaxLifetime)
			},
		},
		{
			name: "only store params leaves a clean DSN",
			dsn:  "postgres://u:p@h/db?*connMaxOpen=20",
			out:  "postgres://u:p@h/db",
			chk: func(req *require.Assertions, c *ConnConfig) {
				req.Equal(20, c.MaxOpenConns)
			},
		},
		{
			name: "retry params",
			dsn:  "postgres://u:p@h/db?*connTryPatience=1m&*connTryBackoffDelay=2s&*connTryTimeout=3s&*connMaxTries=4",
			out:  "postgres://u:p@h/db",
			chk: func(req *require.Assertions, c *ConnConfig) {
				req.Equal(time.Minute, c.ConnTryPatience)
				req.Equal(2*time.Second, c.ConnTryBackoffDelay)
				req.Equal(3*time.Second, c.ConnTryTimeout)
				req.Equal(4, c.ConnTryMax)
			},
		},
		{
			name: "unknown store param",
			dsn:  "postgres://u:p@h/db?*connFoo=1",
			err:  `invalid store configuration for key "*connFoo"`,
		},
		{
			name: "invalid value",
			dsn:  "postgres://u:p@h/db?*connMaxOpen=many",
			err:  `invalid store configuration for key "*connMaxOpen"`,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				c   = &ConnConfig{DataSourceName: tc.dsn}
				err = c.ParseExtra()
			)

			if tc.err != "" {
				req.ErrorContains(err, tc.err)
				return
			}

			req.NoError(err)
			req.Equal(tc.out, c.DataSourceName)

			if tc.chk != nil {
				tc.chk(req, c)
			}
		})
	}
}
