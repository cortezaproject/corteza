package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"io"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type (
	// YamlEncoder is responsible for encoding Corteza resources into
	// a YAML supported format
	YamlEncoder struct{}
)

const (
	paramsKeyWriter = "writer"
)

// Encode encodes the given Corteza resources into some YAML supported format
//
// Encoding should not do any additional processing apart from matching with
// dependencies and runtime validation
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
func (e YamlEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (err error) {
	var (
		out *yaml.Node
		aux *yaml.Node
	)
	_ = aux

	w, err := e.getWriter(p)
	if err != nil {
		return
	}

	switch rt {
	case types.ChartResourceType:
		aux, err = e.encodeCharts(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "chart", aux)
		if err != nil {
			return
		}
	case types.ModuleResourceType:
		aux, err = e.encodeModules(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "module", aux)
		if err != nil {
			return
		}
	case types.ModuleFieldResourceType:
		aux, err = e.encodeModuleFields(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "moduleField", aux)
		if err != nil {
			return
		}
	case types.NamespaceResourceType:
		aux, err = e.encodeNamespaces(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "namespace", aux)
		if err != nil {
			return
		}
	case types.PageResourceType:
		aux, err = e.encodePages(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "page", aux)
		if err != nil {
			return
		}
	case types.PageLayoutResourceType:
		aux, err = e.encodePageLayouts(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "pageLayout", aux)
		if err != nil {
			return
		}

	default:
		out, err = e.encode(ctx, out, p, rt, nodes, tt)
		if err != nil {
			return
		}
	}

	// Don't output nil values since that will produce broken yaml docs
	if out == nil {
		return
	}

	return yaml.NewEncoder(w).Encode(out)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource chart
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeCharts(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeChart(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encodeChart focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeChart(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Chart)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxConfig, err := e.encodeChartConfigC(ctx, p, tt, node, res, res.Config)
	if err != nil {
		return
	}
	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxNamespaceID, err := e.encodeRef(p, res.NamespaceID, "NamespaceID", node, tt)
	if err != nil {
		return
	}
	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"config", auxConfig,
		"createdAt", auxCreatedAt,
		"deletedAt", auxDeletedAt,
		"handle", res.Handle,
		"chartID", res.ID,
		"name", res.Name,
		"namespaceID", auxNamespaceID,
		"updatedAt", auxUpdatedAt,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource module
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeModules(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeModule(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encodeModule focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeModule(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Module)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes

	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxNamespaceID, err := e.encodeRef(p, res.NamespaceID, "NamespaceID", node, tt)
	if err != nil {
		return
	}
	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"config", res.Config,
		"createdAt", auxCreatedAt,
		"deletedAt", auxDeletedAt,
		"handle", res.Handle,
		"moduleID", res.ID,
		"meta", res.Meta,
		"name", res.Name,
		"namespaceID", auxNamespaceID,
		"updatedAt", auxUpdatedAt,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	// When processing module fields, we need to filter out the ones that
	// don't belong to this module
	//
	// @todo offload this to dependency resolution; this is a hack
	children := tt.ChildrenForResourceType(node, types.ModuleFieldResourceType)

	var proc envoyx.NodeSet
	selfRef := node.ToRef()
	for _, c := range children {
		if c.References["ModuleID"].Equals(selfRef) {
			proc = append(proc, c)
		}
	}

	aux, err = e.encodeModuleFields(ctx, p, proc, tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"moduleField", aux,
	)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource moduleField
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeModuleFields(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeModuleField(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encodeModuleField focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeModuleField(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.ModuleField)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes

	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}

	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxModuleID, err := e.encodeRef(p, res.ModuleID, "ModuleID", node, tt)
	if err != nil {
		return
	}

	auxOptions, err := e.encodeModuleFieldOptionsC(ctx, p, tt, node, res, res.Options)
	if err != nil {
		return
	}

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"config", res.Config,
		"createdAt", auxCreatedAt,
		"defaultValue", res.DefaultValue,
		"deletedAt", auxDeletedAt,
		"expressions", res.Expressions,
		"id", res.ID,
		"kind", res.Kind,
		"label", res.Label,
		"moduleID", auxModuleID,
		"multi", res.Multi,
		"name", res.Name,
		"options", auxOptions,
		"place", res.Place,
		"required", res.Required,
		"updatedAt", auxUpdatedAt,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource namespace
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeNamespaces(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeNamespace(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encodeNamespace focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeNamespace(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Namespace)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"createdAt", auxCreatedAt,
		"deletedAt", auxDeletedAt,
		"enabled", res.Enabled,
		"id", res.ID,
		"meta", res.Meta,
		"name", res.Name,
		"slug", res.Slug,
		"updatedAt", auxUpdatedAt,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	aux, err = e.encodeCharts(ctx, p, tt.ChildrenForResourceType(node, types.ChartResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"chart", aux,
	)
	if err != nil {
		return
	}

	aux, err = e.encodeModules(ctx, p, tt.ChildrenForResourceType(node, types.ModuleResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"module", aux,
	)
	if err != nil {
		return
	}

	aux, err = e.encodePages(ctx, p, tt.ChildrenForResourceType(node, types.PageResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"page", aux,
	)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource page
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodePages(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodePage(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encodePage focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodePage(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Page)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxBlocks, err := e.encodePageBlocksC(ctx, p, tt, node, res, res.Blocks)
	if err != nil {
		return
	}

	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxModuleID, err := e.encodeRef(p, res.ModuleID, "ModuleID", node, tt)
	if err != nil {
		return
	}
	auxNamespaceID, err := e.encodeRef(p, res.NamespaceID, "NamespaceID", node, tt)
	if err != nil {
		return
	}
	auxSelfID, err := e.encodeRef(p, res.SelfID, "SelfID", node, tt)
	if err != nil {
		return
	}

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"blocks", auxBlocks,
		"children", res.Children,
		"config", res.Config,
		"createdAt", auxCreatedAt,
		"deletedAt", auxDeletedAt,
		"description", res.Description,
		"handle", res.Handle,
		"id", res.ID,
		"meta", res.Meta,
		"moduleID", auxModuleID,
		"namespaceID", auxNamespaceID,
		"selfID", auxSelfID,
		"title", res.Title,
		"updatedAt", auxUpdatedAt,
		"visible", res.Visible,
		"weight", res.Weight,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	aux, err = e.encodePageLayouts(ctx, p, tt.ChildrenForResourceType(node, types.PageLayoutResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"pageLayout", aux,
	)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource pageLayout
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodePageLayouts(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodePageLayout(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encodePageLayout focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodePageLayout(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.PageLayout)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes

	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxNamespaceID, err := e.encodeRef(p, res.NamespaceID, "NamespaceID", node, tt)
	if err != nil {
		return
	}
	auxOwnedBy, err := e.encodeRef(p, res.OwnedBy, "OwnedBy", node, tt)
	if err != nil {
		return
	}
	auxPageID, err := e.encodeRef(p, res.PageID, "PageID", node, tt)
	if err != nil {
		return
	}
	auxParentID, err := e.encodeRef(p, res.ParentID, "ParentID", node, tt)
	if err != nil {
		return
	}
	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"blocks", res.Blocks,
		"config", res.Config,
		"createdAt", auxCreatedAt,
		"deletedAt", auxDeletedAt,
		"handle", res.Handle,
		"id", res.ID,
		"meta", res.Meta,
		"namespaceID", auxNamespaceID,
		"ownedBy", auxOwnedBy,
		"pageID", auxPageID,
		"parentID", auxParentID,
		"updatedAt", auxUpdatedAt,
		"weight", res.Weight,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Encoding utils
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeTimestamp(p envoyx.EncodeParams, t time.Time) (any, error) {
	if t.IsZero() {
		return nil, nil
	}

	tz := p.Encoder.PreferredTimezone
	if tz != "" {
		tzL, err := time.LoadLocation(tz)
		if err != nil {
			return nil, err
		}
		t = t.In(tzL)
	}

	ly := p.Encoder.PreferredTimeLayout
	if ly == "" {
		ly = time.RFC3339
	}

	return t.Format(ly), nil
}

func (e YamlEncoder) encodeTimestampNil(p envoyx.EncodeParams, t *time.Time) (any, error) {
	if t == nil {
		return nil, nil
	}

	return e.encodeTimestamp(p, *t)
}

func (e YamlEncoder) encodeRef(p envoyx.EncodeParams, id uint64, field string, node *envoyx.Node, tt envoyx.Traverser) (any, error) {
	parent := tt.ParentForRef(node, node.References[field])

	// @todo should we panic instead?
	//       for now gracefully fallback to the ID
	if parent == nil {
		return id, nil
	}

	return parent.Identifiers.FriendlyIdentifier(), nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility functions
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) getWriter(p envoyx.EncodeParams) (out io.Writer, err error) {
	aux, ok := p.Params[paramsKeyWriter]
	if ok {
		out, ok = aux.(io.Writer)
		if ok {
			return
		}
	}

	// @todo consider adding support for managing files from a location
	err = errors.Errorf("YAML encoder expects a writer conforming to io.Writer interface")
	return
}

func safeParentIdentifier(tt envoyx.Traverser, n *envoyx.Node, ref envoyx.Ref) (out string) {
	aux := tt.ParentForRef(n, ref)
	if aux == nil {
		return ref.Identifiers.FriendlyIdentifier()
	}

	return aux.Identifiers.FriendlyIdentifier()
}
