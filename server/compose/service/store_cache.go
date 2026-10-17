package service

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/ctxcache"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	// cachingStore serves namespace, module, module field and label lookups
	// from the cache the context carries, when it carries one
	//
	// Record steps of a workflow load the namespace and the module with its
	// fields and labels for every record they touch; a workflow run carries a
	// cache so a loop reads them once. Callers get copies, they may change
	// what they get.
	cachingStore struct {
		store.Storer
	}
)

func (s cachingStore) LookupComposeNamespaceByID(ctx context.Context, id uint64) (*types.Namespace, error) {
	ns, err := ctxcache.Get(ctx, fmt.Sprintf("compose:namespace:id:%d", id), func() (*types.Namespace, error) {
		return s.Storer.LookupComposeNamespaceByID(ctx, id)
	})
	return copyNamespace(ns), err
}

func (s cachingStore) LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*types.Namespace, error) {
	ns, err := ctxcache.Get(ctx, "compose:namespace:slug:"+slug, func() (*types.Namespace, error) {
		return s.Storer.LookupComposeNamespaceBySlug(ctx, slug)
	})
	return copyNamespace(ns), err
}

func (s cachingStore) LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error) {
	m, err := ctxcache.Get(ctx, fmt.Sprintf("compose:module:id:%d", id), func() (*types.Module, error) {
		return s.Storer.LookupComposeModuleByID(ctx, id)
	})
	return copyModule(m), err
}

func (s cachingStore) LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Module, error) {
	m, err := ctxcache.Get(ctx, fmt.Sprintf("compose:module:handle:%d:%s", namespaceID, handle), func() (*types.Module, error) {
		return s.Storer.LookupComposeModuleByNamespaceIDHandle(ctx, namespaceID, handle)
	})
	return copyModule(m), err
}

func (s cachingStore) SearchComposeModuleFields(ctx context.Context, f types.ModuleFieldFilter) (types.ModuleFieldSet, types.ModuleFieldFilter, error) {
	type result struct {
		set types.ModuleFieldSet
		f   types.ModuleFieldFilter
	}

	r, err := ctxcache.Get(ctx, fmt.Sprintf("compose:module-fields:%v:%d:%d", f.ModuleID, f.Deleted, f.Limit), func() (r result, err error) {
		r.set, r.f, err = s.Storer.SearchComposeModuleFields(ctx, f)
		return
	})
	if err != nil {
		return nil, f, err
	}

	return r.set.Clone(), r.f, nil
}

func (s cachingStore) SearchLabels(ctx context.Context, f labelTypes.LabelFilter) (labelTypes.LabelSet, labelTypes.LabelFilter, error) {
	switch {
	case f.Kind != types.Namespace{}.LabelResourceKind() &&
		f.Kind != types.Module{}.LabelResourceKind() &&
		f.Kind != types.ModuleField{}.LabelResourceKind(),
		len(f.Filter) > 0, f.Name != "", len(f.Value) > 0:
		return s.Storer.SearchLabels(ctx, f)
	}

	type result struct {
		set labelTypes.LabelSet
		f   labelTypes.LabelFilter
	}

	r, err := ctxcache.Get(ctx, fmt.Sprintf("labels:%s:%v:%d", f.Kind, f.ResourceID, f.Limit), func() (r result, err error) {
		r.set, r.f, err = s.Storer.SearchLabels(ctx, f)
		return
	})
	if err != nil {
		return nil, f, err
	}

	return append(labelTypes.LabelSet{}, r.set...), r.f, nil
}

func copyNamespace(ns *types.Namespace) *types.Namespace {
	if ns == nil {
		return nil
	}

	c := *ns
	return &c
}

func copyModule(m *types.Module) *types.Module {
	if m == nil {
		return nil
	}

	return m.Clone()
}
