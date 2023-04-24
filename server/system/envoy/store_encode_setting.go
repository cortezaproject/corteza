package envoy

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

func (e StoreEncoder) prepareSetting(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// @todo existing settings?
	for _, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareSetting with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.SettingValue)
		if !ok {
			panic("unexpected resource type: node expecting type of SettingValue")
		}

		// Run expressions on the nodes
		err = e.runEvals(ctx, false, n)
		if err != nil {
			return
		}

		// @todo merge conflicts if we do existing assertion

		if res.UpdatedAt.IsZero() {
			res.UpdatedAt = time.Now()
		}

		n.Resource = res
	}

	return
}

// encodeSettings encodes a set of resource into the database
func (e StoreEncoder) encodeSettings(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeSetting(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeSetting encodes the resource into the database
func (e StoreEncoder) encodeSetting(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// SettingValues don't have any references so there is no need to handle them

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertSettingValue(ctx, s, n.Resource.(*systemTypes.SettingValue))
		if err != nil {
			return
		}
	}

	return
}
