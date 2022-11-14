package envoy

import (
	"context"
)

type (
	Provider interface {
		NextInverted(ctx context.Context) (*ResourceState, error)
	}
)

// @todo errors!
func Encode(ctx context.Context, p Provider, e Encoder) error {
	// @todo add support for multiple encoders at the same time.
	// The issue occurs with routines and error handling...
	return e.Encode(ctx, p)
}
