package rdbms

import (
	"context"
	"fmt"
)

func (s *Store) Upgrade(ctx context.Context) (err error) {
	if err = UpgradeBeforeTableCreation(ctx, s); err != nil {
		return
	}

	if err = UpgradeCreateTables(ctx, s); err != nil {
		return
	}

	if err = UpgradeAfterTableCreation(ctx, s); err != nil {
		return
	}

	return
}

// UpgradeBeforeTableCreation all actions that need to happen before tables are created
//
// Important note!
// Corteza needs to be updated to the latest patch release under 2022.3.x before upgrading to 2022.9!
func UpgradeBeforeTableCreation(ctx context.Context, s *Store) (err error) {
	// @todo figure out how we will detect if database is on the latest version
	return
}

// UpgradeCreateTables creates all tables needed by RDBMS store for Corteza to function properly
func UpgradeCreateTables(ctx context.Context, s *Store) (err error) {
	for _, t := range Tables() {
		if err = s.SchemaAPI.CreateTable(ctx, s.DB, t); err != nil {
			return fmt.Errorf("could not create table %s: %w", t.Name, err)
		}
	}

	return
}

func UpgradeAfterTableCreation(ctx context.Context, s *Store) (err error) {
	return
}
