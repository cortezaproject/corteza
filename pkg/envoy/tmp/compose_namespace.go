package tmp

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

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

func encodeComposeNamespace(ctx context.Context, s store.Storer, ns *resource.ComposeNamespace) (uint64, error) {
	var (
		res = ns.Res
	)

	res.ID = nextID()
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	// err := store.CreateComposeNamespace(ctx, s, res)
	// if err != nil {
	// 	return nil, err
	// }

	return res.ID, nil
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
