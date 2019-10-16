package rh

import (
	"gopkg.in/Masterminds/squirrel.v1"
)

type (
	// Waiting for PR to be merged:
	// https://github.com/Masterminds/squirrel/pull/206
	//
	// then we can move to squirrel.Fn(...)
	squirrelFunction struct {
		name  string
		fargs []squirrel.Sqlizer
	}
)

func SquirrelFunction(name string, args ...squirrel.Sqlizer) *squirrelFunction {
	return &squirrelFunction{name: name, fargs: args}
}

func (f squirrelFunction) ToSql() (sql string, args []interface{}, err error) {
	var (
		aSql  string
		aArgs []interface{}
	)

	sql = f.name + "("
	args = make([]interface{}, 0)
	for a := 0; a < len(f.fargs); a++ {
		if a > 0 {
			sql += ", "
		}

		aSql, aArgs, err = f.fargs[a].ToSql()
		if err != nil {
			return
		}

		sql += aSql
		args = append(args, aArgs...)
	}
	sql += ")"

	return
}
