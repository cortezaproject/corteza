package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	StoreEncoder struct {
		s  store.Storer
		is *importState
	}
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
)

func NewStoreEncoder(s store.Storer, is *importState) *StoreEncoder {
	return &StoreEncoder{
		s:  s,
		is: is,
	}
}

func (se *StoreEncoder) Encode(ctx context.Context, ess ...*envoy.ExecState) error {
	return store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) (err error) {
		var rID uint64

		for _, es := range ess {
			// @todo what should we do with existing resources?
			// How should we do diffing?
			rID = se.is.Existint(es.Res)
			if rID <= 0 {
				switch res := es.Res.(type) {
				case *resource.ComposeNamespace:
					rID, err = encodeComposeNamespace(ctx, s, res)

				case *resource.ComposeModule:
					rID, err = encodeComposeModule(ctx, s, res, se.is.state[res])

				case *resource.ComposeRecordSet:
					_, err = encodeComposeRecordSet(ctx, s, res, se.is.state[res])
				}

				if err != nil {
					return err
				}
			}

			for _, dr := range es.DepResources {
				se.is.AddRefMapping(dr, es.Res.ResourceType(), rID, es.Res.Identifiers().StringSlice()...)
			}

		}

		return nil
	})
}
