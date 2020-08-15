package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_charts.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/jmoiron/sqlx"
)

// SearchComposeCharts returns all matching rows
//
// This function calls convertComposeChartFilter with the given
// types.ChartFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeCharts(ctx context.Context, f types.ChartFilter) (types.ChartSet, types.ChartFilter, error) {
	q, err := s.convertComposeChartFilter(f)
	if err != nil {
		return nil, f, err
	}

	q = ApplyPaging(q, f.PageFilter)

	scap := f.PerPage
	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Chart, 0, scap)
		res *types.Chart
	)

	return set, f, func() error {
		if f.Count, err = Count(ctx, s.db, q); err != nil || f.Count == 0 {
			return err
		}
		rows, err := s.Query(ctx, q)
		if err != nil {
			return err
		}

		for rows.Next() {
			if res, err = s.internalComposeChartRowScanner(rows, rows.Err()); err != nil {
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

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
func (s Store) LookupComposeChartByID(ctx context.Context, id uint64) (*types.Chart, error) {
	return s.ComposeChartLookup(ctx, squirrel.Eq{
		"cch.id": id,
	})
}

// LookupComposeChartByHandle searches for compose chart by handle (case-insensitive)
func (s Store) LookupComposeChartByHandle(ctx context.Context, handle string) (*types.Chart, error) {
	return s.ComposeChartLookup(ctx, squirrel.Eq{
		"cch.handle": handle,
	})
}

// CreateComposeChart creates one or more rows in compose_chart table
func (s Store) CreateComposeChart(ctx context.Context, rr ...*types.Chart) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ComposeChartTable()).SetMap(s.internalComposeChartEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateComposeChart updates one or more existing rows in compose_chart
func (s Store) UpdateComposeChart(ctx context.Context, rr ...*types.Chart) error {
	return s.PartialUpdateComposeChart(ctx, nil, rr...)
}

// PartialUpdateComposeChart updates one or more existing rows in compose_chart
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateComposeChart(ctx context.Context, onlyColumns []string, rr ...*types.Chart) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateComposeCharts(
				ctx,
				squirrel.Eq{s.preprocessColumn("cch.id", ""): s.preprocessValue(res.ID, "")},
				s.internalComposeChartEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeChart removes one or more rows from compose_chart table
func (s Store) RemoveComposeChart(ctx context.Context, rr ...*types.Chart) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeChartTable("cch")).Where(squirrel.Eq{s.preprocessColumn("cch.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeChartByID removes row from the compose_chart table
func (s Store) RemoveComposeChartByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeChartTable("cch")).Where(squirrel.Eq{s.preprocessColumn("cch.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateComposeCharts removes all rows from the compose_chart table
func (s Store) TruncateComposeCharts(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ComposeChartTable())
}

// ExecUpdateComposeCharts updates all matched (by cnd) rows in compose_chart with given data
func (s Store) ExecUpdateComposeCharts(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ComposeChartTable("cch")).Where(cnd).SetMap(set))
}

// ComposeChartLookup prepares ComposeChart query and executes it,
// returning types.Chart (or error)
func (s Store) ComposeChartLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Chart, error) {
	return s.internalComposeChartRowScanner(s.QueryRow(ctx, s.QueryComposeCharts().Where(cnd)))
}

func (s Store) internalComposeChartRowScanner(row rowScanner, err error) (*types.Chart, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Chart{}
	if _, has := s.config.RowScanners["composeChart"]; has {
		scanner := s.config.RowScanners["composeChart"].(func(rowScanner, *types.Chart) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Name,
			&res.Config,
			&res.NamespaceID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposeChart: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeCharts returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryComposeCharts() squirrel.SelectBuilder {
	return s.Select(s.ComposeChartTable("cch"), s.ComposeChartColumns("cch")...)
}

// ComposeChartTable name of the db table
func (Store) ComposeChartTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_chart" + alias
}

// ComposeChartColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ComposeChartColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "name",
		alias + "config",
		alias + "rel_namespace",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// internalComposeChartEncoder encodes fields from types.Chart to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeChart
// func when rdbms.customEncoder=true
func (s Store) internalComposeChartEncoder(res *types.Chart) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"handle":        res.Handle,
		"name":          res.Name,
		"config":        res.Config,
		"rel_namespace": res.NamespaceID,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}
