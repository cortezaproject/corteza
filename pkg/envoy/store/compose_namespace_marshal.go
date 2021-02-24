package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func newComposeNamespaceFromResource(res *resource.ComposeNamespace, cfg *EncoderConfig) resourceState {
	return &composeNamespace{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

// Prepare prepares the composeNamespace to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *composeNamespace) Prepare(ctx context.Context, pl *payload) (err error) {
	// Try to get the original namespace
	n.ns, err = findComposeNamespaceS(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.ns != nil {
		n.res.Res.ID = n.ns.ID
	}
	return nil
}

// Encode encodes the composeNamespace to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *composeNamespace) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.ns != nil && n.ns.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.ns.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh namespace
	if !exists {
		return store.CreateComposeNamespace(ctx, pl.s, res)
	}

	// Update existing namespace
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeComposeNamespaces(n.ns, res)

	case resource.MergeRight:
		res = mergeComposeNamespaces(res, n.ns)
	}

	err = store.UpdateComposeNamespace(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
