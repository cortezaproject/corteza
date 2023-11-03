package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/pkg/errors"
)

func (e StoreEncoder) prepareResourceTranslation(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// @todo existing resource translations?
	for _, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareResourceTranslation with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.ResourceTranslation)
		if !ok {
			panic("unexpected resource type: node expecting type of ResourceTranslation")
		}

		// Run expressions on the nodes
		err = e.runEvals(ctx, false, n)
		if err != nil {
			return
		}

		res.ID = id.Next()

		// @todo merge conflicts if we do existing assertion

		n.Resource = res
	}

	return
}

// encodeResourceTranslations encodes a set of resource into the database
func (e StoreEncoder) encodeResourceTranslations(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeResourceTranslation(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeResourceTranslation encodes the resource into the database
func (e StoreEncoder) encodeResourceTranslation(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %s", n.ResourceType, n.Resource))
		return
	}

	// Flush to the DB
	err = store.UpsertResourceTranslation(ctx, s, n.Resource.(*types.ResourceTranslation))
	if err != nil {
		return
	}

	return
}
