package rdbms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	ConnConfig struct {
		DriverName     string
		DataSourceName string
		DBName         string

		// MaskedDSN is a DSN with sensitive data masked
		// each connection is handling that independently
		MaskedDSN string

		// MaxOpenConns sets maximum number of open connections to the database
		// defaults to same value as set in the db/sql
		MaxOpenConns int

		// ConnMaxLifetime sets the maximum amount of time a connection may be reused
		// defaults to same value as set in the db/sql
		ConnMaxLifetime time.Duration

		// MaxIdleConns sets the maximum number of connections in the idle connection pool
		// defaults to same value as set in the db/sql
		MaxIdleConns int

		// ConnTryPatience sets time window in which we do not complaining about failed connection tries
		ConnTryPatience time.Duration

		// ConnTryBackoffDelay sets backoff delay after failed try
		ConnTryBackoffDelay time.Duration

		// ConnTryTimeout sets timeout per try
		ConnTryTimeout time.Duration

		// ConnTryMax maximum number of retries for getting the connection
		ConnTryMax int
	}
)

var (
	dsnMasker = regexp.MustCompile(`(.)(?:.*)(.):(.)(?:.*)(.)@`)
)

func (c *ConnConfig) SetDefaults() {
	if c.MaskedDSN == "" {
		c.MaskedDSN = c.DBName + "://********:********@********:********/********"
	}

	if c.MaxIdleConns == 0 {
		// Same as default in the db/sql
		c.MaxIdleConns = 32
	}

	if c.MaxOpenConns == 0 {
		// Same as default in the db/sql
		c.MaxOpenConns = 256
	}

	if c.ConnMaxLifetime == 0 {
		// Same as default in the db/sql
		c.ConnMaxLifetime = 10 * time.Minute
	}

	// ** ** ** ** ** ** ** ** ** ** ** ** ** **

	if c.ConnTryPatience == 0 {
		//c.ConnTryPatience = 1 * time.Minute
	}

	if c.ConnTryBackoffDelay == 0 {
		c.ConnTryBackoffDelay = 10 * time.Second
	}

	if c.ConnTryTimeout == 0 {
		c.ConnTryTimeout = 30 * time.Second
	}

	if c.ConnTryMax == 0 {
		c.ConnTryMax = 99
	}
}

// ParseExtra parses extra params (params starting with *)
// from DSN's querystring (after ?)
func (c *ConnConfig) ParseExtra() (err error) {
	// Make sure we only got qs
	const q = "?"
	var (
		dsn = c.DataSourceName
		qs  string
	)

	if pos := strings.LastIndex(dsn, q); pos == -1 {
		return nil
	} else {
		// Trim qs from DSN, we'll re-attach the remaining params
		c.DataSourceName, qs = dsn[:pos], dsn[pos+1:]
	}

	var vv url.Values
	if vv, err = url.ParseQuery(qs); err != nil {
		return err
	}

	var (
		val string

		parseInt = func(s string) (int, error) {
			if tmp, err := strconv.ParseInt(s, 10, 32); err != nil {
				return 0, err
			} else {
				return int(tmp), nil
			}

		}
	)

	const storePrefixChar = "*"

	for key := range vv {
		val = vv.Get(key)

		if storePrefixChar != key[:1] {
			// skip non-store specific config
			continue
		}

		switch key {
		case "*connTryPatience":
			c.ConnTryPatience, err = time.ParseDuration(val)

		case "*connTryBackoffDelay":
			c.ConnTryBackoffDelay, err = time.ParseDuration(val)

		case "*connTryTimeout":
			c.ConnTryTimeout, err = time.ParseDuration(val)

		case "*connMaxTries":
			c.ConnTryMax, err = parseInt(val)

		case "*connMaxOpen":
			c.MaxOpenConns, err = parseInt(val)

		case "*connMaxLifetime":
			c.ConnMaxLifetime, err = time.ParseDuration(val)

		case "*connMaxIdle":
			c.MaxIdleConns, err = parseInt(val)

		default:
			err = fmt.Errorf("unknown key %q", key)
		}

		if err != nil {
			return fmt.Errorf("invalid store configuration for key %q: %w", key, err)
		}

		delete(vv, key)
	}

	// Encode QS back to DSN
	c.DataSourceName += q + vv.Encode()

	return nil
}
