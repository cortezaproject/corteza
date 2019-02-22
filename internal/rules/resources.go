package rules

import (
	"context"
	"strings"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
)

const (
	delimiter = ":"
)

type resources struct {
	ctx context.Context
	db  *factory.DB
}

func NewResources(ctx context.Context, db *factory.DB) ResourcesInterface {
	return (&resources{}).With(ctx, db)
}

func (r *resources) With(ctx context.Context, db *factory.DB) ResourcesInterface {
	return &resources{
		ctx: ctx,
		db:  db,
	}
}

func (r *resources) identity() uint64 {
	return auth.GetIdentityFromContext(r.ctx).Identity()
}

// IsAllowed function checks granted permission for specific resource and operation. Permission checks on
// global level are not allowed and will always return Deny.
func (r *resources) IsAllowed(resource string, operation string) Access {
	parts := strings.Split(resource, delimiter)

	// Permission check on global level is not allowed.
	if parts[len(parts)-1] == "*" {
		return Deny
	}

	access := r.checkAccess(resource, operation)
	if access == Inherit {
		// Create resource definition for global level.
		parts[len(parts)-1] = "*"

		resource = strings.Join(parts, delimiter)

		// If rule for specific resource is not set we check on global level.
		return r.checkAccess(resource, operation)
	}
	return access
}

func (r *resources) checkAccess(resource string, operation string) Access {
	user := r.identity()
	result := []Access{}
	query := []string{
		// select rules
		"select r.value from sys_rules r",
		// join members
		"inner join sys_role_member m on (m.rel_role = r.rel_role and m.rel_user=?)",
		// add conditions
		"where r.resource=? and r.operation=?",
	}
	queryString := strings.Join(query, " ")
	if err := r.db.Select(&result, queryString, user, resource, operation); err != nil {
		// @todo: log error
		return Deny
	}

	// order by deny, allow
	for _, val := range result {
		if val == Deny {
			return Deny
		}
	}
	for _, val := range result {
		if val == Allow {
			return Allow
		}
	}
	return Inherit
}

func (r *resources) Grant(roleID uint64, rules []Rule) error {
	return r.db.Transaction(func() error {
		var err error
		for _, rule := range rules {
			rule.RoleID = roleID

			switch rule.Value {
			case Inherit:
				_, err = r.db.NamedExec("delete from sys_rules where rel_role=:rel_role and resource=:resource and operation=:operation", rule)
			default:
				err = r.db.Replace("sys_rules", rule)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *resources) Read(roleID uint64) ([]Rule, error) {
	result := []Rule{}

	query := "select * from sys_rules where rel_role = ?"
	if err := r.db.Select(&result, query, roleID); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *resources) Delete(roleID uint64) error {
	query := "delete from sys_rules where rel_role = ?"
	if _, err := r.db.Exec(query, roleID); err != nil {
		return err
	}
	return nil
}
