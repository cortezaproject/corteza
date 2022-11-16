package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	template struct {
		cfg *EncoderConfig

		res *resource.Template
		t   *types.Template
	}
)

// mergeTemplates merges b into a, prioritising a
func mergeTemplates(a, b *types.Template) *types.Template {
	c := *a

	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Language == "" {
		c.Language = b.Language
	}
	if c.Type == "" {
		c.Type = b.Type
	}
	if c.Meta.Short == "" {
		c.Meta.Short = b.Meta.Short
	}
	if c.Meta.Description == "" {
		c.Meta.Description = b.Meta.Description
	}
	if c.Template == "" {
		c.Template = b.Template
	}

	if c.OwnerID == 0 {
		c.OwnerID = b.OwnerID
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}
	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}
	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return &c
}

// findTemplate looks for the template in the resources & the store
//
// Provided resources are prioritized.
func findTemplate(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (u *types.Template, err error) {
	u = resource.FindTemplate(rr, ii)
	if u != nil {
		return u, nil
	}

	return findTemplateStore(ctx, s, makeGenericFilter(ii))
}

// findTemplateStore looks for the template in the store
func findTemplateStore(ctx context.Context, s store.Storer, gf genericFilter) (t *types.Template, err error) {
	if gf.id > 0 {
		t, err = store.LookupTemplateByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if t != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		// Handle
		t, err = store.LookupTemplateByHandle(ctx, s, i)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if t != nil {
			return
		}
	}

	return nil, nil
}
