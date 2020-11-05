package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	storeEncoder struct {
		s  store.Storer
		es *encoderState
	}

	encodingContext struct {
		exists  bool
		partial bool
	}
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
)

// NewStoreEncoder initializes and returns a fresh store encoder
//
// @todo add support for merge options (skip, overwrite, leftMerge, rightMerge)
func NewStoreEncoder(s store.Storer, es *encoderState) envoy.Encoder {
	if es == nil {
		es = NewEncoderState()
	}

	return &storeEncoder{
		s:  s,
		es: es,
	}
}

// Encode encodes the given resource
//
// @todo improve the transaction with channels
func (se *storeEncoder) Encode(ctx context.Context, ee ...*envoy.ExecState) error {
	return store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) (err error) {
		var rState resRefs
		var state resRefs

		for _, e := range ee {
			state = se.es.Get(e.Res)
			ectx := &encodingContext{
				exists:  se.es.Exists(e.Res),
				partial: e.Conflicting,
			}

			switch res := e.Res.(type) {
			case *resource.ComposeNamespace:
				rState, err = encodeComposeNamespace(ctx, ectx, s, state, res)

			case *resource.ComposeModule:
				rState, err = encodeComposeModule(ctx, ectx, s, state, res)

			case *resource.ComposeRecord:
				rState, err = encodeComposeRecord(ctx, ectx, s, state, res)
			}

			if err != nil {
				return err
			}

			for _, dr := range e.DepResources {
				se.es.Merge(dr, rState)
			}

			// @todo this is only relevant for resources conflicting with themselves
			if e.Conflicting {
				se.es.Merge(e.Res, rState)
			}
		}

		return nil
	})
}
