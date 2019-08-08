package permissions

import (
	"fmt"
	"strings"
)

type (
	AccessCheck struct {
		prefix    string
		pkColName string

		resource          Resource
		operation         Operation
		roles             []uint64
		checkExplicitDeny bool
	}
)

func InitAccessCheckFilter(operation Operation, roles []uint64, checkExplicitDeny bool) AccessCheck {
	var ac = AccessCheck{
		operation:         operation,
		roles:             roles,
		checkExplicitDeny: checkExplicitDeny,
		pkColName:         "id",
	}

	return ac
}

func (ac *AccessCheck) BindToEnv(resource Resource, prefix string) *AccessCheck {
	ac.resource = resource
	ac.prefix = prefix
	return ac
}

func (ac *AccessCheck) SetPrimaryKeyName(col string) *AccessCheck {
	ac.pkColName = col
	return ac
}

func (ac AccessCheck) HasOperation() bool {
	return ac.operation != ""
}

// ToSql converts access check to SQL (with args) that will help with filtering
//
// Satisfies squirrel.Sqlizer interface
func (ac AccessCheck) ToSql() (sql string, args []interface{}, err error) {
	if len(ac.roles) == 0 {
		sql = "false"
		return
	}

	sql = fmt.Sprintf(
		`EXISTS (SELECT 1
                            FROM %s_permission_rules
                           WHERE resource = CONCAT(?, %s)
                             AND operation = ?
                             AND access = ?
                             AND rel_role IN (%s))`,

		ac.prefix,
		ac.pkColName,

		// Generate placeholder for every role we have
		strings.Repeat(",?", len(ac.roles))[1:],
	)

	args = []interface{}{
		ac.resource,
		ac.operation,
	}

	if ac.checkExplicitDeny {
		// User has permissions to read on wildcard (*) resource
		// so we need to check if there is any explicit denies
		args = append(args, Deny)
		sql = fmt.Sprintf("NOT %s", sql)
	} else {
		// User is explicitly denied to read on wildcard (*) resource
		// check for all that have explicit allow
		args = append(args, Allow)

	}

	for _, roleID := range ac.roles {
		args = append(args, roleID)
	}

	return sql, args, nil
}
