package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/rbac_rules.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchRbacRules returns all matching rows
//
// This function calls convertRbacRuleFilter with the given
// rbac.RuleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRbacRules(ctx context.Context, f rbac.RuleFilter) (rbac.RuleSet, rbac.RuleFilter, error) {
	var (
		err error
		set []*rbac.Rule
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.rbacRulesSelectBuilder()

		set, err = s.QueryRbacRules(ctx, q, nil)
		return err
	}()
}

// QueryRbacRules queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryRbacRules(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*rbac.Rule) (bool, error),
) ([]*rbac.Rule, error) {
	var (
		set = make([]*rbac.Rule, 0, DefaultSliceCapacity)
		res *rbac.Rule

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalRbacRuleRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateRbacRule creates one or more rows in rbac_rules table
func (s Store) CreateRbacRule(ctx context.Context, rr ...*rbac.Rule) (err error) {
	for _, res := range rr {
		err = s.checkRbacRuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateRbacRules(ctx, s.internalRbacRuleEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateRbacRule updates one or more existing rows in rbac_rules
func (s Store) UpdateRbacRule(ctx context.Context, rr ...*rbac.Rule) error {
	return s.partialRbacRuleUpdate(ctx, nil, rr...)
}

// partialRbacRuleUpdate updates one or more existing rows in rbac_rules
func (s Store) partialRbacRuleUpdate(ctx context.Context, onlyColumns []string, rr ...*rbac.Rule) (err error) {
	for _, res := range rr {
		err = s.checkRbacRuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateRbacRules(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rls.rel_role", ""): store.PreprocessValue(res.RoleID, ""), s.preprocessColumn("rls.resource", ""): store.PreprocessValue(res.Resource, ""), s.preprocessColumn("rls.operation", ""): store.PreprocessValue(res.Operation, ""),
			},
			s.internalRbacRuleEncoder(res).Skip("rel_role", "resource", "operation").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertRbacRule updates one or more existing rows in rbac_rules
func (s Store) UpsertRbacRule(ctx context.Context, rr ...*rbac.Rule) (err error) {
	for _, res := range rr {
		err = s.checkRbacRuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertRbacRules(ctx, s.internalRbacRuleEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRbacRule Deletes one or more rows from rbac_rules table
func (s Store) DeleteRbacRule(ctx context.Context, rr ...*rbac.Rule) (err error) {
	for _, res := range rr {

		err = s.execDeleteRbacRules(ctx, squirrel.Eq{
			s.preprocessColumn("rls.rel_role", ""): store.PreprocessValue(res.RoleID, ""), s.preprocessColumn("rls.resource", ""): store.PreprocessValue(res.Resource, ""), s.preprocessColumn("rls.operation", ""): store.PreprocessValue(res.Operation, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRbacRuleByRoleIDResourceOperation Deletes row from the rbac_rules table
func (s Store) DeleteRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error {
	return s.execDeleteRbacRules(ctx, squirrel.Eq{
		s.preprocessColumn("rls.rel_role", ""):  store.PreprocessValue(roleID, ""),
		s.preprocessColumn("rls.resource", ""):  store.PreprocessValue(resource, ""),
		s.preprocessColumn("rls.operation", ""): store.PreprocessValue(operation, ""),
	})
}

// TruncateRbacRules Deletes all rows from the rbac_rules table
func (s Store) TruncateRbacRules(ctx context.Context) error {
	return s.Truncate(ctx, s.rbacRuleTable())
}

// execLookupRbacRule prepares RbacRule query and executes it,
// returning rbac.Rule (or error)
func (s Store) execLookupRbacRule(ctx context.Context, cnd squirrel.Sqlizer) (res *rbac.Rule, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.rbacRulesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalRbacRuleRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateRbacRules updates all matched (by cnd) rows in rbac_rules with given data
func (s Store) execCreateRbacRules(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.rbacRuleTable()).SetMap(payload))
}

// execUpdateRbacRules updates all matched (by cnd) rows in rbac_rules with given data
func (s Store) execUpdateRbacRules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.rbacRuleTable("rls")).Where(cnd).SetMap(set))
}

// execUpsertRbacRules inserts new or updates matching (by-primary-key) rows in rbac_rules with given data
func (s Store) execUpsertRbacRules(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.rbacRuleTable(),
		set,
		s.preprocessColumn("rel_role", ""),
		s.preprocessColumn("resource", ""),
		s.preprocessColumn("operation", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteRbacRules Deletes all matched (by cnd) rows in rbac_rules with given data
func (s Store) execDeleteRbacRules(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.rbacRuleTable("rls")).Where(cnd))
}

func (s Store) internalRbacRuleRowScanner(row rowScanner) (res *rbac.Rule, err error) {
	res = &rbac.Rule{}

	if _, has := s.config.RowScanners["rbacRule"]; has {
		scanner := s.config.RowScanners["rbacRule"].(func(_ rowScanner, _ *rbac.Rule) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.RoleID,
			&res.Resource,
			&res.Operation,
			&res.Access,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan rbacRule db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryRbacRules returns squirrel.SelectBuilder with set table and all columns
func (s Store) rbacRulesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.rbacRuleTable("rls"), s.rbacRuleColumns("rls")...)
}

// rbacRuleTable name of the db table
func (Store) rbacRuleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "rbac_rules" + alias
}

// RbacRuleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) rbacRuleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_role",
		alias + "resource",
		alias + "operation",
		alias + "access",
	}
}

// {true true false false false false}

// internalRbacRuleEncoder encodes fields from rbac.Rule to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeRbacRule
// func when rdbms.customEncoder=true
func (s Store) internalRbacRuleEncoder(res *rbac.Rule) store.Payload {
	return store.Payload{
		"rel_role":  res.RoleID,
		"resource":  res.Resource,
		"operation": res.Operation,
		"access":    res.Access,
	}
}

// checkRbacRuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkRbacRuleConstraints(ctx context.Context, res *rbac.Rule) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	return nil
}
