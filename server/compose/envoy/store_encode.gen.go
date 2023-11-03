package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/pkg/errors"
)

type (
	// StoreEncoder is responsible for encoding Corteza resources into the
	// database via the Storer or the DAL interface
	//
	// @todo consider having a different encoder for the DAL resources
	StoreEncoder struct{}
)

// Prepare performs some initial processing on the resource before it can be encoded
//
// Preparation runs validation, default value initialization, matching with
// already existing instances, ...
//
// The prepare function receives a set of nodes grouped by the resource type.
// This enables some batching optimization and simplifications when it comes to
// matching with existing resources.
//
// Prepare does not receive any placeholder nodes which are used solely
// for dependency resolution.
func (e StoreEncoder) Prepare(ctx context.Context, p envoyx.EncodeParams, rt string, nn envoyx.NodeSet) (err error) {
	s, err := e.grabStorer(p)
	if err != nil {
		return
	}

	switch rt {
	case types.ChartResourceType:
		return e.prepareChart(ctx, p, s, nn)
	case types.ModuleResourceType:
		return e.prepareModule(ctx, p, s, nn)
	case types.ModuleFieldResourceType:
		return e.prepareModuleField(ctx, p, s, nn)
	case types.NamespaceResourceType:
		return e.prepareNamespace(ctx, p, s, nn)
	case types.PageResourceType:
		return e.preparePage(ctx, p, s, nn)
	case types.PageLayoutResourceType:
		return e.preparePageLayout(ctx, p, s, nn)

	default:
		return e.prepare(ctx, p, s, rt, nn)
	}

	return
}

// Encode encodes the given Corteza resources into the primary store
//
// Encoding should not do any additional processing apart from matching with
// dependencies and runtime validation
//
// The Encode function is called for every resource type where the resource
// appears at the root of the dependency tree.
// All of the root-level resources for that resource type are passed into the function.
// The encoding function must traverse the branches to encode all of the dependencies.
//
// This flow is used to simplify the flow of how resources are encoded into YAML
// (and other documents) as well as to simplify batching.
//
// Encode does not receive any placeholder nodes which are used solely
// for dependency resolution.
func (e StoreEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	s, err := e.grabStorer(p)
	if err != nil {
		return
	}

	switch rt {
	case types.ChartResourceType:
		return e.encodeCharts(ctx, p, s, nodes, tree)

	case types.ModuleResourceType:
		return e.encodeModules(ctx, p, s, nodes, tree)

	case types.ModuleFieldResourceType:
		return e.encodeModuleFields(ctx, p, s, nodes, tree)

	case types.NamespaceResourceType:
		return e.encodeNamespaces(ctx, p, s, nodes, tree)

	case types.PageResourceType:
		return e.encodePages(ctx, p, s, nodes, tree)

	case types.PageLayoutResourceType:
		return e.encodePageLayouts(ctx, p, s, nodes, tree)
	default:
		return e.encode(ctx, p, s, rt, nodes, tree)
	}
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource chart
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepareChart prepares the resources of the given type for encoding
func (e StoreEncoder) prepareChart(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.Chart, len(nn))
	err = e.matchupCharts(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing Charts")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareChart with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.Chart)
		if !ok {
			panic("unexpected resource type: node expecting type of chart")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.setChartDefaults(res)
			if err != nil {
				return err
			}

			err = e.validateChart(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodeCharts encodes a set of resource into the database
func (e StoreEncoder) encodeCharts(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeChart(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeChart encodes the resource into the database
func (e StoreEncoder) encodeChart(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			if auxID == 0 {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertComposeChart(ctx, s, n.Resource.(*types.Chart))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert Chart")
			return
		}
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {

			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	return
}

// matchupCharts returns an index with indicates what resources already exist
func (e StoreEncoder) matchupCharts(ctx context.Context, s store.Storer, uu map[int]types.Chart, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchComposeCharts(ctx, s, types.ChartFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.Chart, len(aa))
	strMap := make(map[string]*types.Chart, len(aa))

	for _, a := range aa {
		strMap[a.Handle] = a
		idMap[a.ID] = a

	}

	var aux *types.Chart
	var ok bool
	for i, n := range nn {

		scope := scopes[i]
		if scope == nil {
			continue
		}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource module
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepareModule prepares the resources of the given type for encoding
func (e StoreEncoder) prepareModule(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.Module, len(nn))
	err = e.matchupModules(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing Modules")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareModule with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.Module)
		if !ok {
			panic("unexpected resource type: node expecting type of module")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.setModuleDefaults(res)
			if err != nil {
				return err
			}

			err = e.validateModule(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodeModules encodes a set of resource into the database
func (e StoreEncoder) encodeModules(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeModule(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}
	err = e.postModulesEncode(ctx, p, s, tree, nn)

	return
}

// encodeModule encodes the resource into the database
func (e StoreEncoder) encodeModule(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			if auxID == 0 {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertComposeModule(ctx, s, n.Resource.(*types.Module))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert Module")
			return
		}
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	nested := make(envoyx.NodeSet, 0, 10)

	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {

			case types.ModuleFieldResourceType:
				err = e.encodeModuleFields(ctx, p, s, nn, tree)
				if err != nil {
					return
				}

				nested = append(nested, nn...)

			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	err = e.encodeModuleExtend(ctx, p, s, n, nested, tree)
	if err != nil {
		err = errors.Wrap(err, "post encode logic failed with errors")
		return
	}

	return
}

// matchupModules returns an index with indicates what resources already exist
func (e StoreEncoder) matchupModules(ctx context.Context, s store.Storer, uu map[int]types.Module, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchComposeModules(ctx, s, types.ModuleFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.Module, len(aa))
	strMap := make(map[string]*types.Module, len(aa))

	for _, a := range aa {
		strMap[a.Handle] = a
		idMap[a.ID] = a

	}

	var aux *types.Module
	var ok bool
	for i, n := range nn {

		scope := scopes[i]
		if scope == nil {
			continue
		}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource moduleField
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepareModuleField prepares the resources of the given type for encoding
func (e StoreEncoder) prepareModuleField(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.ModuleField, len(nn))
	err = e.matchupModuleFields(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing ModuleFields")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareModuleField with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.ModuleField)
		if !ok {
			panic("unexpected resource type: node expecting type of moduleField")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.setModuleFieldDefaults(res)
			if err != nil {
				return err
			}

			err = e.validateModuleField(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodeModuleFields encodes a set of resource into the database
func (e StoreEncoder) encodeModuleFields(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeModuleField(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeModuleField encodes the resource into the database
func (e StoreEncoder) encodeModuleField(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			if auxID == 0 {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	// Custom resource sanitization before saving.
	// This can be used to cleanup arbitrary config fields.
	e.sanitizeModuleFieldBeforeSave(n.Resource.(*types.ModuleField))

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertComposeModuleField(ctx, s, n.Resource.(*types.ModuleField))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert ModuleField")
			return
		}
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {

			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	return
}

// matchupModuleFields returns an index with indicates what resources already exist
func (e StoreEncoder) matchupModuleFields(ctx context.Context, s store.Storer, uu map[int]types.ModuleField, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.ModuleField, len(aa))
	strMap := make(map[string]*types.ModuleField, len(aa))

	for _, a := range aa {
		idMap[a.ID] = a
		strMap[a.Name] = a

	}

	var aux *types.ModuleField
	var ok bool
	for i, n := range nn {

		scope := scopes[i]
		if scope == nil {
			continue
		}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource namespace
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepareNamespace prepares the resources of the given type for encoding
func (e StoreEncoder) prepareNamespace(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.Namespace, len(nn))
	err = e.matchupNamespaces(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing Namespaces")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareNamespace with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.Namespace)
		if !ok {
			panic("unexpected resource type: node expecting type of namespace")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.setNamespaceDefaults(res)
			if err != nil {
				return err
			}

			err = e.validateNamespace(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodeNamespaces encodes a set of resource into the database
func (e StoreEncoder) encodeNamespaces(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeNamespace(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeNamespace encodes the resource into the database
func (e StoreEncoder) encodeNamespace(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			if auxID == 0 {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertComposeNamespace(ctx, s, n.Resource.(*types.Namespace))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert Namespace")
			return
		}
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {

			case types.ChartResourceType:
				err = e.encodeCharts(ctx, p, s, nn, tree)
				if err != nil {
					return
				}

			case types.ModuleResourceType:
				err = e.encodeModules(ctx, p, s, nn, tree)
				if err != nil {
					return
				}

			case types.PageResourceType:
				err = e.encodePages(ctx, p, s, nn, tree)
				if err != nil {
					return
				}

			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	return
}

// matchupNamespaces returns an index with indicates what resources already exist
func (e StoreEncoder) matchupNamespaces(ctx context.Context, s store.Storer, uu map[int]types.Namespace, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchComposeNamespaces(ctx, s, types.NamespaceFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.Namespace, len(aa))
	strMap := make(map[string]*types.Namespace, len(aa))

	for _, a := range aa {
		idMap[a.ID] = a
		strMap[a.Slug] = a

	}

	var aux *types.Namespace
	var ok bool
	for i, n := range nn {

		scope := scopes[i]
		if scope == nil {
			continue
		}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource page
// // // // // // // // // // // // // // // // // // // // // // // // //

// preparePage prepares the resources of the given type for encoding
func (e StoreEncoder) preparePage(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.Page, len(nn))
	err = e.matchupPages(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing Pages")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call preparePage with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.Page)
		if !ok {
			panic("unexpected resource type: node expecting type of page")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.setPageDefaults(res)
			if err != nil {
				return err
			}

			err = e.validatePage(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodePages encodes a set of resource into the database
func (e StoreEncoder) encodePages(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodePage(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodePage encodes the resource into the database
func (e StoreEncoder) encodePage(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			if auxID == 0 {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertComposePage(ctx, s, n.Resource.(*types.Page))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert Page")
			return
		}
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {

			case types.PageLayoutResourceType:
				err = e.encodePageLayouts(ctx, p, s, nn, tree)
				if err != nil {
					return
				}

			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	return
}

// matchupPages returns an index with indicates what resources already exist
func (e StoreEncoder) matchupPages(ctx context.Context, s store.Storer, uu map[int]types.Page, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchComposePages(ctx, s, types.PageFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.Page, len(aa))
	strMap := make(map[string]*types.Page, len(aa))

	for _, a := range aa {
		strMap[a.Handle] = a
		idMap[a.ID] = a

	}

	var aux *types.Page
	var ok bool
	for i, n := range nn {

		scope := scopes[i]
		if scope == nil {
			continue
		}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource pageLayout
// // // // // // // // // // // // // // // // // // // // // // // // //

// preparePageLayout prepares the resources of the given type for encoding
func (e StoreEncoder) preparePageLayout(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.PageLayout, len(nn))
	err = e.matchupPageLayouts(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing PageLayouts")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call preparePageLayout with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.PageLayout)
		if !ok {
			panic("unexpected resource type: node expecting type of pageLayout")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.setPageLayoutDefaults(res)
			if err != nil {
				return err
			}

			err = e.validatePageLayout(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodePageLayouts encodes a set of resource into the database
func (e StoreEncoder) encodePageLayouts(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodePageLayout(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodePageLayout encodes the resource into the database
func (e StoreEncoder) encodePageLayout(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			auxID = safeParentID(tree, n, ref)
			if auxID == 0 {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	// Flush to the DB
	if !n.Evaluated.Skip {
		err = store.UpsertComposePageLayout(ctx, s, n.Resource.(*types.PageLayout))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert PageLayout")
			return
		}
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {

			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	return
}

// matchupPageLayouts returns an index with indicates what resources already exist
func (e StoreEncoder) matchupPageLayouts(ctx context.Context, s store.Storer, uu map[int]types.PageLayout, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchComposePageLayouts(ctx, s, types.PageLayoutFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.PageLayout, len(aa))
	strMap := make(map[string]*types.PageLayout, len(aa))

	for _, a := range aa {
		strMap[a.Handle] = a
		idMap[a.ID] = a

	}

	var aux *types.PageLayout
	var ok bool
	for i, n := range nn {

		scope := scopes[i]
		if scope == nil {
			continue
		}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility functions
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e *StoreEncoder) grabStorer(p envoyx.EncodeParams) (s store.Storer, err error) {
	auxs, ok := p.Params[paramsKeyStorer]
	if !ok {
		err = errors.Errorf("store encoder expects a store conforming to store.Storer interface")
		return
	}

	s, ok = auxs.(store.Storer)
	if !ok {
		err = errors.Errorf("store encoder expects a store conforming to store.Storer interface")
		return
	}

	return
}

func (e *StoreEncoder) runEvals(ctx context.Context, existing bool, n *envoyx.Node) (err error) {
	// Skip if
	if n.Config.SkipIfEval == nil {
		return
	}

	aux, err := expr.EmptyVars().Cast(map[string]any{
		"missing": !existing,
	})
	if err != nil {
		return
	}

	n.Evaluated.Skip, err = n.Config.SkipIfEval.Test(ctx, aux.(*expr.Vars))
	return
}

func (e StoreEncoder) getScopeNodes(ctx context.Context, s store.Storer, nn envoyx.NodeSet) (scopes envoyx.NodeSet, err error) {
	// Get all requested scopes
	scopes = make(envoyx.NodeSet, len(nn))

	err = func() (err error) {
		for i, n := range nn {
			if n.Scope.ResourceType == "" {
				continue
			}

			// For now the scope can only point to namespace so this will do
			var nn envoyx.NodeSet
			nn, err = e.decodeNamespace(ctx, s, e.makeNamespaceFilter(nil, nil, envoyx.ResourceFilter{Identifiers: n.Scope.Identifiers}))
			if err != nil {
				return
			}
			if len(nn) > 1 {
				err = fmt.Errorf("ambiguous scope %v: matches multiple resources", n.Scope)
				return
			}

			// when encoding, it could be missing
			if len(nn) == 0 {
				return
			}

			scopes[i] = nn[0]
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode node scopes")
		return
	}

	return
}

func (e StoreEncoder) decodeNamespace(ctx context.Context, s store.Storer, f types.NamespaceFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposeNamespaces(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = NamespaceToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func (e StoreEncoder) makeNamespaceFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.NamespaceFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.NamespaceID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Slug = hh[0]
	}

	if scope == nil {
		return
	}

	if scope.ResourceType == "" {
		return
	}

	// Overwrite it
	out.NamespaceID = id.Strings(scope.Resource.GetID())

	return
}

func safeParentID(tt envoyx.Traverser, n *envoyx.Node, ref envoyx.Ref) (out uint64) {
	rn := tt.ParentForRef(n, ref)
	if rn == nil {
		return
	}

	return rn.Resource.GetID()
}
