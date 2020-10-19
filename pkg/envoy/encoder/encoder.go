package encoder

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/types"
)

type (
	Encoder interface {
		Encode(context.Context, ...types.Node) error
	}

	Provider interface {
		Next(context.Context) (node types.Node, parentNodes []types.Node, childNodes []types.Node, err error)
	}
)

func Encode(ctx context.Context, p Provider, ee ...Encoder) error {
	for {
		node, pp, _, err := p.Next(ctx)

		// If both node AND error are provided, that means its a dep. conflict
		if node != nil && err != nil {
			// @todo handle dep. conflicts
			continue
		} else if err != nil {
			return err
		}

		if node == nil {
			// No more nodes...
			return nil
		}

		// Upgradable nodes need to be processed based on their parent nodes
		if un, is := node.(types.NodeUpdater); is {
			un.Update(pp...)
		}

		for _, e := range ee {
			err = e.Encode(ctx, node)
			if err != nil {
				return err
			}
		}
	}
}
