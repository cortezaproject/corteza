package rdbms

import (
	"context"
	"fmt"
	automationModels "github.com/cortezaproject/corteza/server/automation/model"
	composeModels "github.com/cortezaproject/corteza/server/compose/model"
	federationModels "github.com/cortezaproject/corteza/server/federation/model"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	systemModels "github.com/cortezaproject/corteza/server/system/model"
	"go.uber.org/zap"
)

func (s *Store) Upgrade(ctx context.Context) (err error) {
	for _, fix := range fixes {
		if err = fix(ctx, s); err != nil {
			return
		}
	}

	err = createTablesFromModels(
		ctx,
		s.log(ctx),
		s.DataDefiner,
		systemModels.Models(),
		composeModels.Models(),
		automationModels.Models(),
		federationModels.Models(),
	)

	if err != nil {
		return err
	}

	return
}

func createTablesFromModels(ctx context.Context, log *zap.Logger, dd ddl.DataDefiner, sets ...dal.ModelSet) (err error) {
	var (
		tbl *ddl.Table
	)

	for _, mm := range sets {
		for _, m := range mm {
			log.Debug("verifying primary store table", zap.String("table", m.Ident))

			if tbl, err = dd.ConvertModel(m); err != nil {
				return fmt.Errorf("can not convert model %q to table: %w", m.Ident, err)
			}

			_, err = dd.TableLookup(ctx, m.Ident)
			if err != nil {
				if !errors.IsNotFound(err) {
					return fmt.Errorf("can not do a table lookup: %w", err)
				}

				if err = dd.TableCreate(ctx, tbl); err != nil {
					return fmt.Errorf("can not create table from model %q: %w", m.Ident, err)
				}
			}

			for _, idx := range tbl.Indexes {
				if idx.Ident == ddl.PRIMARY_KEY {
					// @todo move this decision to drivers!
					continue
				}

				_, err = dd.IndexLookup(ctx, idx.Ident, idx.TableIdent)
				if err != nil && !errors.IsNotFound(err) {
					return
				} else if errors.IsNotFound(err) {
					if err = dd.IndexCreate(ctx, tbl.Ident, idx); err != nil {
						return fmt.Errorf("can not create index %q on table %q: %w", idx.Ident, tbl.Ident, err)
					}
				}
			}
		}
	}

	return nil
}

func dropTable(ctx context.Context, s *Store, table string) error {
	_, err := s.DataDefiner.TableLookup(ctx, table)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}

	return s.DataDefiner.TableDrop(ctx, table)
}

// addColumn adds column on a table but only if table exists!
//
// If table does not exist adding column can be skipped
// We can assume that 2nd step of the upgrade process will include the column
func addColumn(ctx context.Context, s *Store, table string, attr *dal.Attribute) error {
	tbl, err := s.DataDefiner.TableLookup(ctx, table)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}

	if tbl.ColumnByIdent(attr.StoreIdent()) != nil {
		return nil
	}

	s.log(ctx).Info(fmt.Sprintf("extending %q table with %q column", table, attr.StoreIdent()))

	col, err := s.DataDefiner.ConvertAttribute(attr)
	if err != nil {
		return err
	}

	return s.DataDefiner.ColumnAdd(ctx, table, col)
}

// dropColumns removes columns from a table but only if table exists!
//
// If table does not exist column removing can be skipped
// We can assume that 2nd step of the upgrade process will omit the column
func dropColumns(ctx context.Context, s *Store, table string, cc ...string) error {
	tbl, err := s.DataDefiner.TableLookup(ctx, table)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}

	for _, c := range cc {
		if tbl.ColumnByIdent(c) == nil {
			// column does not exist, nothing to do
			continue
		}

		s.log(ctx).Info(fmt.Sprintf("dropping %q column from %q", c, table))
		if err := s.DataDefiner.ColumnDrop(ctx, table, c); err != nil {
			return err
		}
	}
	return nil
}

// renameColumn renames columns from a table but only if table exists!
//
// If table does not exist column renaming can be skipped
// We can assume that 2nd step of the upgrade process will have columns properly nameed
func renameColumn(ctx context.Context, s *Store, table string, from, to string) error {
	tbl, err := s.DataDefiner.TableLookup(ctx, table)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}

	if tbl.ColumnByIdent(from) == nil {
		// from column does not exist, nothing to do
		return nil
	}

	if tbl.ColumnByIdent(to) != nil {
		// to column already exists, nothing to do
		return nil
	}

	s.log(ctx).Info(fmt.Sprintf("renaming %q column on table %q to %q", from, table, to))
	if err := s.DataDefiner.ColumnRename(ctx, table, from, to); err != nil {
		return err
	}

	return nil
}

// tableNames returns table names that Corteza creates
func tableNames() (tnames []string) {
	cortezaModels := append(systemModels.Models(), composeModels.Models()...)
	cortezaModels = append(cortezaModels, automationModels.Models()...)
	cortezaModels = append(cortezaModels, federationModels.Models()...)

	for _, m := range cortezaModels {
		tnames = append(tnames, m.Ident)
	}

	return tnames
}
