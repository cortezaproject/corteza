package rh

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/ql"
)

func ParseOrder(order string, valid ...string) (out []string, err error) {
	var (
		// Sort parser
		sp = ql.NewParser()

		// Sort columns
		sc ql.Columns

		whitelist = map[string]bool{}
	)

	for _, col := range valid {
		if i := strings.Index(col, "."); i > -1 {
			whitelist[col[i+1:]] = true
		}

		whitelist[col] = true
	}

	sp.OnIdent = func(i ql.Ident) (ql.Ident, error) {
		if !whitelist[i.Value] {
			return i, errors.Errorf("unknown order-by column %q", i.Value)
		}

		i.Value += " "
		return i, nil
	}

	if sc, err = sp.ParseColumns(order); err != nil {
		return
	}

	out = sc.Strings()

	return
}
