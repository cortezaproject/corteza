package envoy

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/filter"
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
	seen := make(map[string]uint64)

	for _, n := range nn {
		err = e.encodeResourceTranslation(ctx, p, s, n, tree, seen)
		if err != nil {
			return
		}
	}

	return
}

// encodeResourceTranslation encodes the resource into the database
func (e StoreEncoder) encodeResourceTranslation(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser, seen map[string]uint64) (err error) {
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

	res := n.Resource.(*types.ResourceTranslation)

	// Reuse the ID of the matching translation so the upsert updates it instead
	// of tripping the unique (lang, resource, key) constraint; translations can
	// repeat inside the import batch or already exist in the database
	key := strings.ToLower(res.Lang.String() + "|" + res.Resource + "|" + res.K)
	if ID, ok := seen[key]; ok {
		res.ID = ID
	} else {
		if err = e.setExistingResourceTranslationID(ctx, s, res); err != nil {
			return
		}

		seen[key] = res.ID
	}

	// Flush to the DB
	err = store.UpsertResourceTranslation(ctx, s, res)
	if err != nil {
		return
	}

	return
}

// setExistingResourceTranslationID looks up the stored translation with the
// same (lang, resource, key) and reuses its ID; deleted translations are
// matched too as the unique constraint covers them as well
func (e StoreEncoder) setExistingResourceTranslationID(ctx context.Context, s store.Storer, res *types.ResourceTranslation) (err error) {
	ee, _, err := store.SearchResourceTranslations(ctx, s, types.ResourceTranslationFilter{
		Lang:     res.Lang.String(),
		Resource: res.Resource,
		Deleted:  filter.StateInclusive,
	})
	if err != nil && err != store.ErrNotFound {
		return err
	}

	for _, ex := range ee {
		if strings.EqualFold(ex.K, res.K) {
			res.ID = ex.ID
			break
		}
	}

	return nil
}
