package objstore

import (
	"context"
)

func Healthcheck(s Store) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return s.Healthcheck(ctx)
	}
}
