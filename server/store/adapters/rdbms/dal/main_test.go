package dal_test

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms"
	"github.com/stretchr/testify/require"
	"os"
	"sort"
	"testing"

	_ "github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/mysql"
	_ "github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/postgres"
	_ "github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
)

type (
	// @todo refactor tests that use kv and migrate to kvv
	//       no need to keep and maintain both.
	kv  map[string]any
	kvv map[string][]any
)

var (
	s *rdbms.Store
)

// @todo refactor tests to follow rdbms/tests, using connectionInfo...
// @todo should be part of the general DAL testing suite
// @todo new ENV var for enabling(def)/disabling temp. model (table) creation (DB_PERSIST_MODELS=true)
func TestMain(m *testing.M) {
	var (
		dsn = os.Getenv("DB_DSN")
		log = logger.MakeDebugLogger()
		ctx = context.Background()
		err error
		aux store.Storer
	)

	if len(dsn) == 0 {
		dsn = "sqlite3+debug://file::memory:?cache=shared&mode=memory"
	}

	// ctx = logger.ContextWithValue(context.Background(), log)
	aux, err = store.Connect(ctx, log, dsn, true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not connect: %v", err)
		os.Exit(1)
	}

	s = aux.(*rdbms.Store)

	m.Run()
}

// truncates table
func truncate(ctx context.Context, table string) error {
	table = s.Dialect.QuoteIdent(table)
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
	return err
}

func (r kv) CountValues() map[string]uint {
	out := make(map[string]uint)

	for k := range r {
		out[k]++
	}

	return out
}

func (r kv) GetValue(k string, place uint) (any, error) {
	return r[k], nil
}

func (r kv) SetValue(k string, place uint, v any) error {
	r[k] = v
	return nil
}

// String function returns string representation of the kv with sorted keys
func (r kv) String() string {
	// sort keys from map
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// build string by iterating over sorted keys and appending values
	var out string
	for i, k := range keys {
		if i > 0 {
			out += " "
		}

		out += fmt.Sprintf("%s=%v", k, r[k])
	}

	return out
}

func (r kvv) Set(k string, v ...any) kvv {
	r[k] = v
	return r
}

func (r kvv) CountValues() map[string]uint {
	out := make(map[string]uint)

	for k := range r {
		out[k] = uint(len(r[k]))
	}

	return out
}

func (r kvv) GetValue(k string, p uint) (any, error) {
	if r[k] == nil || len(r[k]) <= int(p) {
		return nil, fmt.Errorf("kvv: out of bounds")
	}

	return r[k][p], nil
}

func (r kvv) SetValue(k string, p uint, v any) error {
	if r[k] == nil {
		r[k] = make([]any, 0, 1)
	} else if len(r[k]) < int(p) {
		r[k][p] = v
	} else {
		r[k] = append(r[k], v)
	}

	return nil
}

// String function returns string representation of the kv with sorted keys
func (r kvv) String() string {
	// sort keys from map
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// build string by iterating over sorted keys and appending values
	var out string
	for i, k := range keys {
		if i > 0 {
			out += " "
		}

		out += fmt.Sprintf("%s=%v", k, r[k])
	}

	return out
}

func qlParse(req *require.Assertions, q string) *ql.ASTNode {
	n, err := ql.NewParser().Parse(q)
	req.NoError(err)
	return n
}
