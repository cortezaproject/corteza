package envoy

import (
	"context"
)

type (
	Encoder interface {
		Encode(ctx context.Context, ss ...*ExecState) error
	}

	Provider interface {
		Next(ctx context.Context) (*ExecState, error)
	}
)

func Encode(ctx context.Context, p Provider, ee ...Encoder) error {
	for {
		state, err := p.Next(ctx)
		if err != nil {
			return err
		} else if state == nil {
			return nil
		}

		for _, e := range ee {
			// Any dep conflicts should be handled by the Encoder
			err = e.Encode(ctx, state)
			if err != nil {
				return err
			}
		}
	}
}
