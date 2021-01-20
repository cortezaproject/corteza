package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/federation_nodes.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchFederationNodes returns all matching rows
//
// This function calls convertFederationNodeFilter with the given
// types.NodeFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchFederationNodes(ctx context.Context, f types.NodeFilter) (types.NodeSet, types.NodeFilter, error) {
	var (
		err error
		set []*types.Node
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertFederationNodeFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryFederationNodes(ctx, q, f.Check)
		return err
	}()
}

// QueryFederationNodes queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryFederationNodes(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Node) (bool, error),
) ([]*types.Node, error) {
	var (
		set = make([]*types.Node, 0, DefaultSliceCapacity)
		res *types.Node

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalFederationNodeRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, err
			} else if !chk {
				continue
			}
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupFederationNodeByID searches for federation node by ID
//
// It returns federation node
func (s Store) LookupFederationNodeByID(ctx context.Context, id uint64) (*types.Node, error) {
	return s.execLookupFederationNode(ctx, squirrel.Eq{
		s.preprocessColumn("fdn.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupFederationNodeByBaseURLSharedNodeID searches for node by shared-node-id and base-url
func (s Store) LookupFederationNodeByBaseURLSharedNodeID(ctx context.Context, base_url string, shared_node_id uint64) (*types.Node, error) {
	return s.execLookupFederationNode(ctx, squirrel.Eq{
		s.preprocessColumn("fdn.base_url", ""):       store.PreprocessValue(base_url, ""),
		s.preprocessColumn("fdn.shared_node_id", ""): store.PreprocessValue(shared_node_id, ""),
	})
}

// LookupFederationNodeBySharedNodeID searches for node by shared-node-id
func (s Store) LookupFederationNodeBySharedNodeID(ctx context.Context, shared_node_id uint64) (*types.Node, error) {
	return s.execLookupFederationNode(ctx, squirrel.Eq{
		s.preprocessColumn("fdn.shared_node_id", ""): store.PreprocessValue(shared_node_id, ""),
	})
}

// CreateFederationNode creates one or more rows in federation_nodes table
func (s Store) CreateFederationNode(ctx context.Context, rr ...*types.Node) (err error) {
	for _, res := range rr {
		err = s.checkFederationNodeConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateFederationNodes(ctx, s.internalFederationNodeEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateFederationNode updates one or more existing rows in federation_nodes
func (s Store) UpdateFederationNode(ctx context.Context, rr ...*types.Node) error {
	return s.partialFederationNodeUpdate(ctx, nil, rr...)
}

// partialFederationNodeUpdate updates one or more existing rows in federation_nodes
func (s Store) partialFederationNodeUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Node) (err error) {
	for _, res := range rr {
		err = s.checkFederationNodeConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateFederationNodes(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("fdn.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalFederationNodeEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertFederationNode updates one or more existing rows in federation_nodes
func (s Store) UpsertFederationNode(ctx context.Context, rr ...*types.Node) (err error) {
	for _, res := range rr {
		err = s.checkFederationNodeConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertFederationNodes(ctx, s.internalFederationNodeEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationNode Deletes one or more rows from federation_nodes table
func (s Store) DeleteFederationNode(ctx context.Context, rr ...*types.Node) (err error) {
	for _, res := range rr {

		err = s.execDeleteFederationNodes(ctx, squirrel.Eq{
			s.preprocessColumn("fdn.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationNodeByID Deletes row from the federation_nodes table
func (s Store) DeleteFederationNodeByID(ctx context.Context, ID uint64) error {
	return s.execDeleteFederationNodes(ctx, squirrel.Eq{
		s.preprocessColumn("fdn.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateFederationNodes Deletes all rows from the federation_nodes table
func (s Store) TruncateFederationNodes(ctx context.Context) error {
	return s.Truncate(ctx, s.federationNodeTable())
}

// execLookupFederationNode prepares FederationNode query and executes it,
// returning types.Node (or error)
func (s Store) execLookupFederationNode(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Node, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.federationNodesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalFederationNodeRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateFederationNodes updates all matched (by cnd) rows in federation_nodes with given data
func (s Store) execCreateFederationNodes(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.federationNodeTable()).SetMap(payload))
}

// execUpdateFederationNodes updates all matched (by cnd) rows in federation_nodes with given data
func (s Store) execUpdateFederationNodes(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.federationNodeTable("fdn")).Where(cnd).SetMap(set))
}

// execUpsertFederationNodes inserts new or updates matching (by-primary-key) rows in federation_nodes with given data
func (s Store) execUpsertFederationNodes(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationNodeTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteFederationNodes Deletes all matched (by cnd) rows in federation_nodes with given data
func (s Store) execDeleteFederationNodes(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.federationNodeTable("fdn")).Where(cnd))
}

func (s Store) internalFederationNodeRowScanner(row rowScanner) (res *types.Node, err error) {
	res = &types.Node{}

	if _, has := s.config.RowScanners["federationNode"]; has {
		scanner := s.config.RowScanners["federationNode"].(func(_ rowScanner, _ *types.Node) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.SharedNodeID,
			&res.BaseURL,
			&res.Status,
			&res.Contact,
			&res.PairToken,
			&res.AuthToken,
			&res.CreatedBy,
			&res.UpdatedBy,
			&res.DeletedBy,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan federationNode db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryFederationNodes returns squirrel.SelectBuilder with set table and all columns
func (s Store) federationNodesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.federationNodeTable("fdn"), s.federationNodeColumns("fdn")...)
}

// federationNodeTable name of the db table
func (Store) federationNodeTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "federation_nodes" + alias
}

// FederationNodeColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) federationNodeColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "shared_node_id",
		alias + "base_url",
		alias + "status",
		alias + "contact",
		alias + "pair_token",
		alias + "auth_token",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false true}

// internalFederationNodeEncoder encodes fields from types.Node to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationNode
// func when rdbms.customEncoder=true
func (s Store) internalFederationNodeEncoder(res *types.Node) store.Payload {
	return store.Payload{
		"id":             res.ID,
		"name":           res.Name,
		"shared_node_id": res.SharedNodeID,
		"base_url":       res.BaseURL,
		"status":         res.Status,
		"contact":        res.Contact,
		"pair_token":     res.PairToken,
		"auth_token":     res.AuthToken,
		"created_by":     res.CreatedBy,
		"updated_by":     res.UpdatedBy,
		"deleted_by":     res.DeletedBy,
		"created_at":     res.CreatedAt,
		"updated_at":     res.UpdatedAt,
		"deleted_at":     res.DeletedAt,
	}
}

// checkFederationNodeConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkFederationNodeConstraints(ctx context.Context, res *types.Node) error {
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
