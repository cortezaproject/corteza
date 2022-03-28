package rdbms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/schema"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	txRetryOnErrHandler func(int, error) bool

	Functions struct {
		// returns lower case text
		LOWER func(interface{}) exp.SQLFunctionExpression

		// returns date part of the input (YYYY-MM-DD)
		DATE func(interface{}) exp.SQLFunctionExpression
	}

	Config struct {
		DriverName     string
		DataSourceName string
		DBName         string

		Dialect goqu.DialectWrapper

		// Forces debug mode on RDBMS driver
		Debug bool

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

		// ConnTryMax maximum number of retrys for getting the connection
		ConnTryMax int

		// Disable transactions
		TxDisabled bool

		// How many times should we retry failed transaction?
		TxMaxRetries int

		// TxRetryErrHandler should return true if transaction should be retried
		//
		// Because retry algorithm varies between concrete rdbms implementations
		//
		// Handler must return true if failed transaction should be replied
		// and false if we're safe to terminate it
		TxRetryErrHandler txRetryOnErrHandler

		ErrorHandler store.ErrorHandler

		Functions Functions

		// additional per-resource filters used when searching
		// these filters can modify expression used for querying the database
		Filters extendedFilters

		// schema upgrade interface
		Upgrader schema.Upgrader
	}
)

var (
	dsnMasker = regexp.MustCompile("(.)(?:.*)(.):(.)(?:.*)(.)@")
)

// MaskedDSN replaces username & password from DSN string  dso it's usable for logging
func (c *Config) MaskedDSN() string {
	return dsnMasker.ReplaceAllString(c.DataSourceName, "$1****$2:$3****$4@")
}

func (c *Config) SetDefaults() {
	//if c.PlaceholderFormat == nil {
	//	c.PlaceholderFormat = squirrel.Question
	//}

	if c.TxMaxRetries == 0 {
		c.TxMaxRetries = TxRetryHardLimit
	}

	//if c.TxRetryErrHandler == nil {
	//	// Default transaction retry handler
	//	c.TxRetryErrHandler = TxNoRetry
	//}
	//
	//if c.ErrorHandler == nil {
	//	c.ErrorHandler = ErrHandlerFallthrough
	//}
	//
	//if c.UpsertBuilder == nil {
	//	c.UpsertBuilder = UpsertBuilder
	//}

	// ** ** ** ** ** ** ** ** ** ** ** ** ** **

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

	//if c.TriggerHandlers == nil {
	//	c.TriggerHandlers = TriggerHandlers{}
	//}
	//
	//if c.SqlSortHandler == nil {
	//	c.SqlSortHandler = SqlSortHandler
	//}

	//if c.Functions == nil {
	//	c.Functions = func(ident string, args ...interface{}) (*db.RawExpr, error) {
	//		switch strings.ToUpper(ident) {
	//		case "LOWER":
	//			if len(args) != 1 {
	//				return nil, errors.Internal("LOWER() function expects exactly 1 argument")
	//			}
	//			return db.Raw("LOWER(?)", args[0]), nil
	//		}
	//
	//		return nil, errors.Internal("unknown function %q", ident)
	//	}
	//}
	c.Functions.SetDefaults()
	c.Filters.SetDefaults()
}

// ParseExtra parses extra params (params starting with *)
// from DSN's querystring (after ?)
func (c *Config) ParseExtra() (err error) {
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

func (f *Functions) SetDefaults() {
	if f.LOWER == nil {
		f.LOWER = func(value interface{}) exp.SQLFunctionExpression {
			return goqu.Func("LOWER", value)
		}
	}

	if f.DATE == nil {
		f.DATE = func(value interface{}) exp.SQLFunctionExpression {
			return goqu.Func("DATE", value)
		}
	}

}
