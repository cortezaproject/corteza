package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/pkg/errors"
)

func (e StoreEncoder) prepareRbacRule(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// @todo existing RBAC rules?
	for _, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareRbacRule with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*rbac.Rule)
		if !ok {
			panic("unexpected resource type: node expecting type of RbacRule")
		}

		// Run expressions on the nodes
		err = e.runEvals(ctx, false, n)
		if err != nil {
			return
		}

		// @todo merge conflicts if we do existing assertion

		n.Resource = res
	}

	return
}

// encodeRbacRules encodes a set of resource into the database
func (e StoreEncoder) encodeRbacRules(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeRbacRule(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeRbacRule encodes the resource into the database
func (e StoreEncoder) encodeRbacRule(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
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
	err = store.UpsertRbacRule(ctx, s, n.Resource.(*rbac.Rule))
	if err != nil {
		return
	}

	return
}
