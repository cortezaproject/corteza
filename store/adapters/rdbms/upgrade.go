package rdbms

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

func (s *Store) Upgrade(ctx context.Context) (err error) {
	var (
		tableExists bool
	)

	for _, t := range Tables() {

		tableExists, err = s.SchemaAPI.TableExists(ctx, s.DB, t.Name)
		if err != nil {
			return fmt.Errorf("could not check table %q existance: %w", t.Name, err)
		}

		if !tableExists {
			s.log(ctx).Debug("creating table", zap.String("table", t.Name))
			if err = s.SchemaAPI.CreateTable(ctx, s.DB, t); err != nil {
				return fmt.Errorf("could not create table %q: %w", t.Name, err)
			}
		}
	}

	fixes := []func(context.Context, *Store) error{
		fix202209_extendComposeModuleForPrivacyAndDAL,
		fix202209_extendComposeModuleFieldsForPrivacyAndDAL,
	}

	for _, fix := range fixes {
		if err = fix(ctx, s); err != nil {
			return
		}
	}

	return
}
