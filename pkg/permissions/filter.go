package permissions

import (
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// ResourceFilter is Helper for *Filter structs
	//
	// It is used to provide filtering on db level and meant to be used
	// mainly for checking for read operations.
	//
	// It creates a complex SQL syntax for permission checking depending on the
	// permissions of the user:
	//
	//  - if user is superuser no extra checks are made, simple TRUE sql is returned
	//  - if user is member of one or more roles a query is assembled that checks
	//    for allow & deny permissions for each resource
	//  - if one of the roles has wildcard ALLOW / DENY rule this is then the final check
	//  - we check everyone role rules for each resource
	//  - if everyone role has wildcard ALLOW / DENY rule this is then the final check
	//  - fallback access check is added at the end
	//
	// Resulting SQL check SHOULD reflect rules check ("overall flow" in the header of
	// ruleset_checks.go file)
	//
	ResourceFilter struct {
		dbTable   string
		pkColName string

		resource  Resource
		operation Operation

		chk interface {
			Check(res Resource, op Operation, roles ...uint64) (v Access)
		}

		fallback Access

		superuser bool
		roles     []uint64
	}
)

func NewSuperuserFilter() *ResourceFilter {
	return &ResourceFilter{superuser: true}
}

func (rf *ResourceFilter) Build(pkColName string) *ResourceFilter {
	rf.pkColName = pkColName
	return rf
}

func (rf ResourceFilter) ToSql() (sql string, args []interface{}, err error) {
	if rf.superuser {
		return "TRUE", nil, nil
	}

	// selects first rule for res+op+role
	// rules are ordered by access - denies first
	// end query will return 1 row with 1 column - FALSE if user has at least one DENY rule
	// and TRUE if there is at least one ALLOW
	//
	// Final query is then wrapped in simple CASE statement that casts NULL (no rules)
	// to TRUE. So: no rule == inherit
	base := squirrel.
		Select(fmt.Sprintf("access = %d", Allow)).
		From(rf.dbTable).
		Where(squirrel.Eq{"operation": rf.operation}).
		Where(squirrel.Expr(fmt.Sprintf("resource = CONCAT(?, %s)", rf.pkColName), rf.resource)).
		OrderBy("access").
		Limit(1)

	var (
		checks = []squirrel.Sqlizer{}

		expTRUE  = squirrel.Expr("TRUE")
		expFALSE = squirrel.Expr("FALSE")

		check = func(rr ...uint64) squirrel.Sqlizer {
			return squirrel.And{base.Where(squirrel.Eq{"rel_role": rr})}
		}

		build = func(ss ...squirrel.Sqlizer) (sql string, args []interface{}, err error) {
			return rh.SquirrelFunction("COALESCE", append(checks, ss...)...).ToSql()
		}
	)

	if len(rf.roles) > 0 {
		// Add per-resource check for all roles
		checks = append(checks, check(rf.roles...))

		if rf.chk != nil {
			switch rf.chk.Check(rf.resource.AppendWildcard(), rf.operation, rf.roles...) {
			// Explicit deny/allow on wildcard:
			// Add false/true to check-list and return it immediately
			case Deny:
				return build(expFALSE)
			case Allow:
				return build(expTRUE)
			}
		}
	}

	// Add per-resource check for Everyone
	checks = append(checks, check(EveryoneRoleID))

	if rf.chk != nil {
		switch rf.chk.Check(rf.resource.AppendWildcard(), rf.operation, rf.roles...) {
		// Explicit deny/allow on wildcard:
		// Add false/true to check-list and return it immediately
		case Deny:
			return build(expFALSE)
		case Allow:
			return build(expTRUE)
		}
	}

	// Fallback access
	if rf.fallback == Deny {
		return build(expFALSE)
	} else {
		return build(expTRUE)
	}
}
