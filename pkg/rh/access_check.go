package rh

import (
	"fmt"

	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	expAccessCheck struct {
		ac   permissions.AccessCheck
		expr squirrel.Sqlizer
	}
)

// ExpAccessCheck wraps expression into access control check
//
// If access-check is enabled and user can not perform an operation
// condition will fail
func ExpAccessCheck(ac permissions.AccessCheck, expr squirrel.Sqlizer) expAccessCheck {
	if ac.HasOperation() && !ac.HasResource() {
		// This is not runtime error
		// developer should make sure access control has resource!
		panic("can not wrap expression with access check without resource")
	}

	return expAccessCheck{
		ac:   ac,
		expr: expr,
	}
}

func (eac expAccessCheck) ToSql() (string, []interface{}, error) {
	if !eac.ac.HasOperation() {
		return eac.expr.ToSql()
	}

	acSql, acArgs, err := eac.ac.ToSql()
	if err != nil {
		return "", nil, err
	}

	sql, args, err := eac.expr.ToSql()
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf(`IF(%s, %s, false)`, acSql, sql), append(acArgs, args...), nil
}
