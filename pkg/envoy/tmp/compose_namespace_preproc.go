package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeNamespacePreproc struct {
		is *importState
		s  store.Storer
	}
)

func NewComposeNamespacePreproc(is *importState, s store.Storer) envoy.Processor {
	return &composeNamespacePreproc{
		is: is,
		s:  s,
	}
}

func (p *composeNamespacePreproc) Process(ctx context.Context, state *envoy.ExecState) error {
	// @todo can we/should we have the same pattern as with decoder's CanDecode?
	res, is := state.Res.(*resource.ComposeNamespace)
	if !is {
		return nil
	}

	// Check if the given namespace exists
	// Namespaces are top-level resources, so they don't depend on anything;
	// no other checks are needed

	ns, err := loadComposeNamespace(ctx, p.s, res.SearchQuery())
	if err != nil {
		return err
	}

	if ns != nil {
		p.is.AddExisting(res, ns.ID)
	}

	return nil
}
