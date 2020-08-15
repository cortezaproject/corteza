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
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/jmoiron/sqlx"
)

// SearchRbacRules returns all matching rows
//
// This function calls convertRbacRuleFilter with the given
// permissions.RuleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRbacRules(ctx context.Context, f permissions.RuleFilter) (permissions.RuleSet, permissions.RuleFilter, error) {
	q := s.QueryRbacRules()

	scap := DefaultSliceCapacity

	var (
		set = make([]*permissions.Rule, 0, scap)
		res *permissions.Rule
	)

	return set, f, func() error {
		rows, err := s.Query(ctx, q)
		if err != nil {
			return err
		}

		for rows.Next() {
			if res, err = s.internalRbacRuleRowScanner(rows, rows.Err()); err != nil {
				if cerr := rows.Close(); cerr != nil {
					return fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
				}

				return err
			}

			set = append(set, res)
		}

		return rows.Close()
	}()
}

// CreateRbacRule creates one or more rows in rbac_rules table
func (s Store) CreateRbacRule(ctx context.Context, rr ...*permissions.Rule) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.RbacRuleTable()).SetMap(s.internalRbacRuleEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateRbacRule updates one or more existing rows in rbac_rules
func (s Store) UpdateRbacRule(ctx context.Context, rr ...*permissions.Rule) error {
	return s.PartialUpdateRbacRule(ctx, nil, rr...)
}

// PartialUpdateRbacRule updates one or more existing rows in rbac_rules
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateRbacRule(ctx context.Context, onlyColumns []string, rr ...*permissions.Rule) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateRbacRules(
				ctx,
				squirrel.Eq{s.preprocessColumn("rls.rel_role", ""): s.preprocessValue(res.RoleID, ""),
					s.preprocessColumn("rls.resource", ""):  s.preprocessValue(res.Resource, ""),
					s.preprocessColumn("rls.operation", ""): s.preprocessValue(res.Operation, ""),
				},
				s.internalRbacRuleEncoder(res).Skip("rel_role", "resource", "operation").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRbacRule removes one or more rows from rbac_rules table
func (s Store) RemoveRbacRule(ctx context.Context, rr ...*permissions.Rule) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.RbacRuleTable("rls")).Where(squirrel.Eq{s.preprocessColumn("rls.rel_role", ""): s.preprocessValue(res.RoleID, ""),
				s.preprocessColumn("rls.resource", ""):  s.preprocessValue(res.Resource, ""),
				s.preprocessColumn("rls.operation", ""): s.preprocessValue(res.Operation, ""),
			}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRbacRuleByRoleIDResourceOperation removes row from the rbac_rules table
func (s Store) RemoveRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.RbacRuleTable("rls")).Where(squirrel.Eq{s.preprocessColumn("rls.rel_role", ""): s.preprocessValue(roleID, ""),

		s.preprocessColumn("rls.resource", ""): s.preprocessValue(resource, ""),

		s.preprocessColumn("rls.operation", ""): s.preprocessValue(operation, ""),
	}))
}

// TruncateRbacRules removes all rows from the rbac_rules table
func (s Store) TruncateRbacRules(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.RbacRuleTable())
}

// ExecUpdateRbacRules updates all matched (by cnd) rows in rbac_rules with given data
func (s Store) ExecUpdateRbacRules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.RbacRuleTable("rls")).Where(cnd).SetMap(set))
}

// RbacRuleLookup prepares RbacRule query and executes it,
// returning permissions.Rule (or error)
func (s Store) RbacRuleLookup(ctx context.Context, cnd squirrel.Sqlizer) (*permissions.Rule, error) {
	return s.internalRbacRuleRowScanner(s.QueryRow(ctx, s.QueryRbacRules().Where(cnd)))
}

func (s Store) internalRbacRuleRowScanner(row rowScanner, err error) (*permissions.Rule, error) {
	if err != nil {
		return nil, err
	}

	var res = &permissions.Rule{}
	if _, has := s.config.RowScanners["rbacRule"]; has {
		scanner := s.config.RowScanners["rbacRule"].(func(rowScanner, *permissions.Rule) error)
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
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for RbacRule: %w", err)
	} else {
		return res, nil
	}
}

// QueryRbacRules returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryRbacRules() squirrel.SelectBuilder {
	return s.Select(s.RbacRuleTable("rls"), s.RbacRuleColumns("rls")...)
}

// RbacRuleTable name of the db table
func (Store) RbacRuleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "rbac_rules" + alias
}

// RbacRuleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) RbacRuleColumns(aa ...string) []string {
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

// internalRbacRuleEncoder encodes fields from permissions.Rule to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeRbacRule
// func when rdbms.customEncoder=true
func (s Store) internalRbacRuleEncoder(res *permissions.Rule) store.Payload {
	return store.Payload{
		"rel_role":  res.RoleID,
		"resource":  res.Resource,
		"operation": res.Operation,
		"access":    res.Access,
	}
}
