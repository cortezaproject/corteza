package tmp

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func moduleFilterFromGeneric(gf genericFilter) types.ModuleFilter {
	f := types.ModuleFilter{
		Name:   gf.Name,
		Handle: gf.Ref,
	}
	if gf.ID > 0 {
		f.Query = fmt.Sprintf("moduleID=%d", gf.ID)
	}

	return f
}

func encodeComposeModule(ctx context.Context, s store.Storer, mod *resource.ComposeModule, rm resMap) (uint64, error) {
	var (
		res = mod.Res
	)

	res.ID = nextID()
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	// Namespace...
	// A module can exist under a single namespace, so this is good enough for now
	nsID := uint64(0)
	for _, v := range rm["compose:namespace"] {
		nsID = v
		break
	}
	res.NamespaceID = nsID

	// @todo fields

	return res.ID, nil
}

func loadComposeModule(ctx context.Context, s store.Storer, f types.ModuleFilter) (*types.Module, error) {
	mdd, f, err := store.SearchComposeModules(ctx, s, f)
	if err != nil {
		return nil, err
	}

	if len(mdd) > 0 {
		return mdd[0], nil
	}

	return nil, nil
}
