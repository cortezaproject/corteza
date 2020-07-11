package db

import (
	"context"
	"github.com/titpetric/factory"
)

func Healthcheck() func(ctx context.Context) error {
	const name = "default"
	return func(ctx context.Context) error {
		db, err := factory.Database.Get(name)
		if err != nil {
			return err
		}

		err = db.PingContext(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}
