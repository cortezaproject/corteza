package plugin

// Collection of boot-lifecycle related functions
// that exec plugin functions

import (
	"context"

	sdk "github.com/cortezaproject/corteza/server/sdk/plugin"
	"go.uber.org/zap"
)

func (pp Set) Setup(log *zap.Logger) error {
	for _, p := range pp {
		d, is := p.def.(sdk.Setup)
		if !is {
			continue
		}

		err := d.Setup(log)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pp Set) Initialize(ctx context.Context, log *zap.Logger) error {
	for _, p := range pp {
		d, is := p.def.(sdk.Initialize)
		if !is {
			continue
		}

		err := d.Initialize(ctx, log)
		if err != nil {
			return err
		}
	}

	return nil
}
