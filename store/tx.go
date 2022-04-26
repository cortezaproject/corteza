package store

import (
	"context"
)

func Tx(ctx context.Context, s Storer, fn func(context.Context, Storer) error) error {
	return s.Tx(ctx, fn)
}
