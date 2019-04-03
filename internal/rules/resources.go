package rules

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
)

const (
	everyoneRoleId uint64 = 1
	defaultAccess  Access = Deny
)

type (
	resources struct {
		ctx context.Context
		db  *factory.DB
	}

	// CheckAccessFunc function.
	CheckAccessFunc func() Access
)

func NewResources(ctx context.Context, db *factory.DB) ResourcesInterface {
	return (&resources{}).With(ctx, db)
}

func (rr *resources) With(ctx context.Context, db *factory.DB) ResourcesInterface {
	return &resources{
		ctx: ctx,
		db:  db,
	}
}

func (rr *resources) identity() uint64 {
	return auth.GetIdentityFromContext(rr.ctx).Identity()
}

// IsAllowed function checks granted permission for specific resource and operation. Permission checks on
// global level are not allowed and will always return Deny.
func (rr *resources) Check(r Resource, operation string, fallbacks ...CheckAccessFunc) Access {
	// Number one check, do we have a valid identity?
	if !auth.GetIdentityFromContext(rr.ctx).Valid() {
		return Deny
	}

	// Allow anything if we're in a testing context
	if v := rr.ctx.Value("testing"); v != nil {
		return Allow
	}

	if !r.IsValid() {
		// Make sure we do not let through wildcard or undefined resources
		return Deny
	}

	// Resource-specific check
	checks := []CheckAccessFunc{
		func() Access { return rr.checkAccess(r, operation) },
		func() Access { return rr.checkAccessEveryone(r, operation) },
	}

	if r.IsAppendable() {
		wc := r.AppendWildcard()

		checks = append(
			checks,
			func() Access { return rr.checkAccess(wc, operation) },
			func() Access { return rr.checkAccessEveryone(wc, operation) },
		)
	}

	checks = append(checks, fallbacks...)

	for _, check := range checks {
		if access := check(); access != Inherit {
			return access
		}
	}
	return defaultAccess
}

//
func (rr *resources) checkAccess(resource Resource, operation string) Access {
	var result = make([]Access, 0)

	user := rr.identity()

	query := "" +
		"SELECT r.value " +
		"  FROM sys_rules r" +
		"       INNER JOIN sys_role_member m ON (m.rel_role = r.rel_role AND m.rel_user = ?)" +
		" WHERE r.resource = ? AND  r.operation = ?"

	if err := rr.db.Select(&result, query, user, resource, operation); err != nil {
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

func (rr *resources) checkAccessEveryone(resource Resource, operation string) Access {
	var result = make([]Access, 0)

	query := "" +
		"SELECT r.value " +
		"  FROM sys_rules r" +
		" WHERE r.rel_role = ?  AND r.resource = ? AND  r.operation = ?"

	if err := rr.db.Select(&result, query, everyoneRoleId, resource, operation); err != nil {
		// @todo: log error
		return Deny
	}

	if len(result) > 0 {
		return result[0]
	}

	return Inherit
}

func (rr *resources) Grant(roleID uint64, rules []Rule) error {
	return rr.db.Transaction(func() error {
		var err error
		for _, rule := range rules {
			rule.RoleID = roleID

			switch rule.Value {
			case Inherit:
				_, err = rr.db.NamedExec("delete from sys_rules where rel_role=:rel_role and resource=:resource and operation=:operation", rule)
			default:
				err = rr.db.Replace("sys_rules", rule)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (rr *resources) Read(roleID uint64) ([]Rule, error) {
	result := []Rule{}

	query := "select * from sys_rules where rel_role = ?"
	if err := rr.db.Select(&result, query, roleID); err != nil {
		return nil, err
	}
	return result, nil
}

func (rr *resources) Delete(roleID uint64) error {
	query := "delete from sys_rules where rel_role = ?"
	if _, err := rr.db.Exec(query, roleID); err != nil {
		return err
	}
	return nil
}
