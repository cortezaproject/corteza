package service

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/ctxcache"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
)

type (
	countingLookupStore struct {
		store.Storer
		calls int
	}
)

func (s *countingLookupStore) LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*types.Namespace, error) {
	s.calls++
	return &types.Namespace{ID: 1, Slug: slug}, nil
}

func (s *countingLookupStore) LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error) {
	s.calls++
	return &types.Module{ID: id, Handle: "mod"}, nil
}

func (s *countingLookupStore) SearchComposeModuleFields(ctx context.Context, f types.ModuleFieldFilter) (types.ModuleFieldSet, types.ModuleFieldFilter, error) {
	s.calls++
	return types.ModuleFieldSet{{ID: 10, Name: "f1"}}, f, nil
}

func (s *countingLookupStore) SearchLabels(ctx context.Context, f labelTypes.LabelFilter) (labelTypes.LabelSet, labelTypes.LabelFilter, error) {
	s.calls++
	return labelTypes.LabelSet{{Kind: f.Kind, Name: "l"}}, f, nil
}

func TestCachingStore(t *testing.T) {
	var (
		req = require.New(t)
		ctx = ctxcache.With(context.Background(), ctxcache.New(time.Minute))
	)

	t.Run("lookups repeated in a run read the store once", func(t *testing.T) {
		var (
			raw = &countingLookupStore{}
			s   = cachingStore{Storer: raw}
		)

		for i := 0; i < 3; i++ {
			_, err := s.LookupComposeNamespaceBySlug(ctx, "ns")
			req.NoError(err)
			_, err = s.LookupComposeModuleByID(ctx, 2)
			req.NoError(err)
			_, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{2}})
			req.NoError(err)
			_, _, err = s.SearchLabels(ctx, labelTypes.LabelFilter{Kind: types.Module{}.LabelResourceKind(), ResourceID: []uint64{2}})
			req.NoError(err)
		}

		req.Equal(4, raw.calls)
	})

	t.Run("callers get copies", func(t *testing.T) {
		var (
			s = cachingStore{Storer: &countingLookupStore{}}
		)

		m, err := s.LookupComposeModuleByID(ctx, 3)
		req.NoError(err)
		m.Handle = "changed"

		ff, _, err := s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{3}})
		req.NoError(err)
		ff[0].Name = "changed"

		m, err = s.LookupComposeModuleByID(ctx, 3)
		req.NoError(err)
		req.Equal("mod", m.Handle)

		ff, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{3}})
		req.NoError(err)
		req.Equal("f1", ff[0].Name)
	})

	t.Run("lookups outside a run are not cached", func(t *testing.T) {
		var (
			raw = &countingLookupStore{}
			s   = cachingStore{Storer: raw}
		)

		_, err := s.LookupComposeModuleByID(context.Background(), 2)
		req.NoError(err)
		_, err = s.LookupComposeModuleByID(context.Background(), 2)
		req.NoError(err)
		req.Equal(2, raw.calls)
	})

	t.Run("labels of other resources are not cached", func(t *testing.T) {
		var (
			raw = &countingLookupStore{}
			s   = cachingStore{Storer: raw}
		)

		for i := 0; i < 2; i++ {
			_, _, err := s.SearchLabels(ctx, labelTypes.LabelFilter{Kind: types.Page{}.LabelResourceKind(), ResourceID: []uint64{2}})
			req.NoError(err)
		}

		req.Equal(2, raw.calls)
	})
}
