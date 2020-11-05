package tmp

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeNamespacePreproc struct {
		es *encoderState
		s  store.Storer
	}
)

func NewComposeNamespacePreproc(es *encoderState, s store.Storer) envoy.Processor {
	return &composeNamespacePreproc{
		es: es,
		s:  s,
	}
}

func (p *composeNamespacePreproc) Process(ctx context.Context, state *envoy.ExecState) error {
	res, is := state.Res.(*resource.ComposeNamespace)
	if !is {
		return nil
	}

	// Check if the given namespace exists
	// Namespaces are top-level resources, so they don't depend on anything;
	// no other checks are needed

	ns, err := loadComposeNamespace(ctx, p.s, namespaceFilterFromGeneric(GenericFilter(res.Identifiers())))
	if err != nil {
		return err
	}

	if ns != nil {
		p.es.SetExists(res)
		p.es.Set(res, res.ResourceType(), ns.ID, res.Identifiers().StringSlice()...)
	}

	return nil
}

func encodeComposeNamespace(ctx context.Context, ectx *encodingContext, s store.Storer, state resRefs, res *resource.ComposeNamespace) (resRefs, error) {
	var (
		ns     = res.Res
		rState = make(resRefs)
	)

	ns.ID = state.Get(res)
	if ns.ID <= 0 {
		ns.ID = nextID()
	}
	if ns.CreatedAt.IsZero() {
		ns.CreatedAt = time.Now()
	}

	rState.Set(resource.COMPOSE_NAMESPACE_RESOURCE_TYPE, ns.ID, res.Identifiers().StringSlice()...)

	if !ectx.partial && !ectx.exists {
		err := store.CreateComposeNamespace(ctx, s, ns)
		if err != nil {
			return nil, err
		}
	}

	return rState, nil
}

// Utils

func namespaceFilterFromGeneric(gf genericFilter) types.NamespaceFilter {
	f := types.NamespaceFilter{
		Name: gf.Name,
		Slug: gf.Ref,
	}
	if gf.ID > 0 {
		f.Query = fmt.Sprintf("namespaceID=%d", gf.ID)
	}

	return f
}

func findNamespace(ctx context.Context, s store.Storer, rr []resource.Interface, ii resource.Identifiers) (ns *types.Namespace, err error) {
	// Try to find it in the parent resources
	var nsRes *resource.ComposeNamespace
	walkResources(rr, func(r resource.Interface) error {
		nsR, ok := r.(*resource.ComposeNamespace)
		if !ok {
			return nil
		}

		if nsR.Identifiers().HasAny(r.Identifiers().StringSlice()...) {
			nsRes = nsR
		}
		return nil
	})

	// Found it
	if nsRes != nil {
		return nsRes.Res, nil
	}

	// Go in the store
	f := namespaceFilterFromGeneric(GenericFilter(ii))
	return loadComposeNamespace(ctx, s, f)
}

func loadComposeNamespace(ctx context.Context, s store.Storer, f types.NamespaceFilter) (*types.Namespace, error) {
	nss, f, err := store.SearchComposeNamespaces(ctx, s, f)
	if err != nil {
		return nil, err
	}

	if len(nss) > 0 {
		return nss[0], nil
	}

	return nil, nil
}
