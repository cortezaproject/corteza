package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func NewTemplateFromResource(res *resource.Template, cfg *EncoderConfig) resourceState {
	return &template{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

// Prepare prepares the template to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *template) Prepare(ctx context.Context, pl *payload) (err error) {
	// Try to get the original template
	n.t, err = findTemplateS(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.t != nil {
		n.res.Res.ID = n.t.ID
	}
	return nil
}

// Encode encodes the template to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *template) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.t != nil && n.t.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.t.ID
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

	// Userstamps
	us := n.res.Userstamps()
	if us != nil {
		if us.OwnedBy != nil {
			res.OwnerID = us.OwnedBy.UserID
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh template
	if !exists {
		return store.CreateTemplate(ctx, pl.s, res)
	}

	// Update existing template
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeTemplates(n.t, res)

	case resource.MergeRight:
		res = mergeTemplates(res, n.t)
	}

	err = store.UpdateTemplate(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
