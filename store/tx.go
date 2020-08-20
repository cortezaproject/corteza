package store

import "context"

func Tx(ctx context.Context, s Storable, fn func(context.Context, Storable) error) error {
	return s.Tx(ctx, fn)
}
