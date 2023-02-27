package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
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
	for fieldLabel, ref := range n.References {
		rn := tree.ParentForRef(n, ref)
		if rn == nil {
			err = fmt.Errorf("missing node for ref %v", ref)
			return
		}

		auxID = rn.Resource.GetID()
		if auxID == 0 {
			err = fmt.Errorf("related resource doesn't provide an ID")
			return
		}

		err = n.Resource.SetValue(fieldLabel, 0, auxID)
		if err != nil {
			return
		}
	}

	// Flush to the DB
	err = store.UpsertRbacRule(ctx, s, n.Resource.(*rbac.Rule))
	if err != nil {
		return
	}

	return
}
