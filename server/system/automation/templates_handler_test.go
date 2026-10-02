package automation

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

type (
	templateServiceMocked struct {
		templateService

		byHandle func(context.Context, string) (*types.Template, error)
		rendered func(uint64, types.TemplateRenderAux)
	}
)

func (svc templateServiceMocked) FindByHandle(ctx context.Context, h string) (*types.Template, error) {
	return svc.byHandle(ctx, h)
}

func (svc templateServiceMocked) Render(_ context.Context, templateID uint64, _ string, _ map[string]interface{}, _ map[string]string, aux types.TemplateRenderAux) (io.ReadSeeker, error) {
	svc.rendered(templateID, aux)
	return nil, nil
}

func TestResolveTemplateID(t *testing.T) {
	var (
		ctx = context.Background()
		svc = templateServiceMocked{
			byHandle: func(_ context.Context, h string) (*types.Template, error) {
				if h != "header" {
					return nil, fmt.Errorf("unknown template %q", h)
				}

				return &types.Template{ID: 200}, nil
			},
		}
	)

	t.Run("unset argument", func(t *testing.T) {
		req := require.New(t)
		ID, err := resolveTemplateID(ctx, svc, false, 0, "", nil)
		req.NoError(err)
		req.Zero(ID)
	})

	t.Run("by ID", func(t *testing.T) {
		req := require.New(t)
		ID, err := resolveTemplateID(ctx, svc, true, 100, "", nil)
		req.NoError(err)
		req.Equal(uint64(100), ID)
	})

	t.Run("by handle", func(t *testing.T) {
		req := require.New(t)
		ID, err := resolveTemplateID(ctx, svc, true, 0, "header", nil)
		req.NoError(err)
		req.Equal(uint64(200), ID)
	})

	t.Run("by template", func(t *testing.T) {
		req := require.New(t)
		ID, err := resolveTemplateID(ctx, svc, true, 0, "", &types.Template{ID: 300})
		req.NoError(err)
		req.Equal(uint64(300), ID)
	})

	t.Run("unknown handle", func(t *testing.T) {
		req := require.New(t)
		_, err := resolveTemplateID(ctx, svc, true, 0, "missing", nil)
		req.Error(err)
	})
}

func TestRenderPassesAuxTemplates(t *testing.T) {
	var (
		req = require.New(t)

		renderedID  uint64
		renderedAux types.TemplateRenderAux

		svc = templateServiceMocked{
			byHandle: func(_ context.Context, h string) (*types.Template, error) {
				return &types.Template{ID: 200}, nil
			},
			rendered: func(ID uint64, aux types.TemplateRenderAux) {
				renderedID, renderedAux = ID, aux
			},
		}

		h = templatesHandler{tSvc: svc}
	)

	args := &templatesRenderArgs{
		hasLookup: true,
		lookupID:  100,

		DocumentType: "application/pdf",

		hasHeaderTemplate: true,
		headerTemplateID:  300,

		hasFooterTemplate:    true,
		footerTemplateHandle: "footer",
	}

	_, err := h.render(context.Background(), args)
	req.NoError(err)

	req.Equal(uint64(100), renderedID)
	req.Equal(uint64(300), renderedAux.HeaderTemplateID)
	req.Equal(uint64(200), renderedAux.FooterTemplateID)
}

func TestRenderWithoutAuxTemplates(t *testing.T) {
	var (
		req = require.New(t)

		renderedAux types.TemplateRenderAux

		svc = templateServiceMocked{
			rendered: func(_ uint64, aux types.TemplateRenderAux) {
				renderedAux = aux
			},
		}

		h = templatesHandler{tSvc: svc}
	)

	_, err := h.render(context.Background(), &templatesRenderArgs{hasLookup: true, lookupID: 100})
	req.NoError(err)

	// nothing set, service falls back to the template's own meta
	req.Zero(renderedAux.HeaderTemplateID)
	req.Zero(renderedAux.FooterTemplateID)
}
