package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testTemplates(t *testing.T, s store.Templates) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Template {
			// minimum data set for new template
			return &types.Template{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Template) {
			req := require.New(t)
			req.NoError(s.TruncateTemplates(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateTemplate(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateTemplates(ctx))
		template := makeNew("TemplateCRUD")
		req.NoError(s.CreateTemplate(ctx, template))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, template := truncAndCreate(t)
		fetched, err := s.LookupTemplateByID(ctx, template.ID)
		req.NoError(err)
		req.Equal(template.Handle, fetched.Handle)
		req.Equal(template.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("lookup by handle", func(t *testing.T) {
		req, template := truncAndCreate(t)
		fetched, err := s.LookupTemplateByHandle(ctx, template.Handle)
		req.NoError(err)
		req.Equal(template.Handle, fetched.Handle)
		req.Equal(template.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, template := truncAndCreate(t)
		template.Handle = "TemplateCRUD+2"

		req.NoError(s.UpdateTemplate(ctx, template))

		updated, err := s.LookupTemplateByID(ctx, template.ID)
		req.NoError(err)
		req.Equal(template.Handle, updated.Handle)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, template := truncAndCreate(t)
			template.Handle = "TemplateCRUD+2"

			req.NoError(s.UpsertTemplate(ctx, template))

			updated, err := s.LookupTemplateByID(ctx, template.ID)
			req.NoError(err)
			req.Equal(template.Handle, updated.Handle)
		})

		t.Run("new", func(t *testing.T) {
			template := makeNew("upsert me")
			template.Handle = "ComposeChartCRUD+2"

			req.NoError(s.UpsertTemplate(ctx, template))

			upserted, err := s.LookupTemplateByID(ctx, template.ID)
			req.NoError(err)
			req.Equal(template.Handle, upserted.Handle)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Template", func(t *testing.T) {
			req, template := truncAndCreate(t)
			req.NoError(s.DeleteTemplate(ctx, template))
			_, err := s.LookupTemplateByID(ctx, template.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, template := truncAndCreate(t)
			req.NoError(s.DeleteTemplateByID(ctx, template.ID))
			_, err := s.LookupTemplateByID(ctx, template.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Template{
			makeNew("one-one"),
			makeNew("one-two"),
			makeNew("two-one"),
			makeNew("two-two"),
			makeNew("two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateTemplates(ctx))
		req.NoError(s.CreateTemplate(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchTemplates(ctx, types.TemplateFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchTemplates(ctx, types.TemplateFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchTemplates(ctx, types.TemplateFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchTemplates(ctx, types.TemplateFilter{Handle: "two-one"})
		req.NoError(err)
		req.Len(set, 1)

		_ = f // dummy
	})
}
