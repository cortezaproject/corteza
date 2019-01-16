package rules

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/system/types"
)

type Access string

const (
	Allow   Access = "yes"
	Deny           = "no"
	Inherit        = ""
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

func (r *resources) CheckAccessMulti(resource string, operation string) error {
	user := r.identity()
	result := []Access{}
	query := []string{
		// select rules
		"select r.value from sys_rules r",
		// join members
		"inner join sys_team_member m on (m.rel_team = r.rel_team and m.rel_user=?)",
		// add conditions
		"where r.resource LIKE ? and r.operation=?",
	}
	resource = strings.Replace(resource, "*", "%", -1)
	queryString := strings.Join(query, " ")
	if err := r.db.Select(&result, queryString, user, resource, operation); err != nil {
		return err
	}

	// order by deny, allow
	for _, val := range result {
		if val == Deny {
			return errors.New("Access not allowed")
		}
	}
	for _, val := range result {
		if val == Allow {
			return nil
		}
	}
	return errors.New("Access not allowed")
}

func (r *resources) CheckAccess(resource string, operation string) error {
	user := r.identity()
	result := []Access{}
	query := []string{
		// select rules
		"select r.value from sys_rules r",
		// join members
		"inner join sys_team_member m on (m.rel_team = r.rel_team and m.rel_user=?)",
		// add conditions
		"where r.resource=? and r.operation=?",
	}
	queryString := strings.Join(query, " ")
	if err := r.db.Select(&result, queryString, user, resource, operation); err != nil {
		return err
	}

	// order by deny, allow
	for _, val := range result {
		if val == Deny {
			return errors.New("Access not allowed")
		}
	}
	for _, val := range result {
		if val == Allow {
			return nil
		}
	}
	return errors.New("Access not allowed")
}

func (r *resources) Grant(resource string, teamID uint64, operations []string, value Access) error {
	row := types.Rules{
		TeamID:   teamID,
		Resource: resource,
		Value:    string(value),
	}
	var err error
	for _, operation := range operations {
		row.Operation = operation
		switch value {
		case Inherit:
			_, err = r.db.NamedExec("delete from sys_rules where rel_team=:rel_team and resource=:resource and operation=:operation", row)
		default:
			err = r.db.Replace("sys_rules", row)
		}
		if err != nil {
			break
		}
	}
	return err
}
