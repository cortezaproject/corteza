package schema

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type (
	Schema struct{}

	// schemaUpgrader provides procedures to upgrade rdbms store tables
	Upgrader interface {
		SetLogger(*zap.Logger)
		Before(context.Context) error
		CreateTable(context.Context, *Table) error
		After(context.Context) error
	}
)

// Upgrade upgrades given rdbms schema by adding all missing tables and applying changes to existing
func Upgrade(ctx context.Context, sud Upgrader) (err error) {
	if sud == nil {
		return fmt.Errorf("can not upgrade database schema, upgrade interface not set")
	}

	if err = sud.Before(ctx); err != nil {
		return fmt.Errorf("could not run \"before\" upgrade procedures: %w", err)
	}

	for _, t := range tables() {
		if err = sud.CreateTable(ctx, t); err != nil {
			return fmt.Errorf("could not create table %s: %w", t.Name, err)
		}
	}

	if err = sud.After(ctx); err != nil {
		return fmt.Errorf("could not run \"after\" upgrade procedures: %w", err)
	}

	return nil
}
