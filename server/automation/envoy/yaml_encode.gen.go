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

	"github.com/cortezaproject/corteza/server/automation/types"
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
	case types.WorkflowResourceType:
		aux, err = e.encodeWorkflows(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "workflow", aux)
		if err != nil {
			return
		}

	case types.TriggerResourceType:
		aux, err = e.encodeTriggers(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "trigger", aux)
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
} // // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource workflow
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeWorkflows(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeWorkflow(ctx, p, n, tt)
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

// encodeWorkflow focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeWorkflow(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Workflow)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxCreatedBy, err := e.encodeRef(p, res.CreatedBy, "CreatedBy", node, tt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}
	auxDeletedBy, err := e.encodeRef(p, res.DeletedBy, "DeletedBy", node, tt)
	if err != nil {
		return
	}

	auxOwnedBy, err := e.encodeRef(p, res.OwnedBy, "OwnedBy", node, tt)
	if err != nil {
		return
	}

	auxRunAs, err := e.encodeRef(p, res.RunAs, "RunAs", node, tt)
	if err != nil {
		return
	}

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}
	auxUpdatedBy, err := e.encodeRef(p, res.UpdatedBy, "UpdatedBy", node, tt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"createdAt", auxCreatedAt,
		"createdBy", auxCreatedBy,
		"deletedAt", auxDeletedAt,
		"deletedBy", auxDeletedBy,
		"enabled", res.Enabled,
		"handle", res.Handle,
		"id", res.ID,
		"keepSessions", res.KeepSessions,
		"meta", res.Meta,
		"ownedBy", auxOwnedBy,
		"paths", res.Paths,
		"runAs", auxRunAs,
		"scope", res.Scope,
		"steps", res.Steps,
		"trace", res.Trace,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	aux, err = e.encodeTriggers(ctx, p, tt.ChildrenForResourceType(node, types.TriggerResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"trigger", aux,
	)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource trigger
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeTriggers(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeTrigger(ctx, p, n, tt)
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

// encodeTrigger focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeTrigger(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Trigger)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes

	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxCreatedBy, err := e.encodeRef(p, res.CreatedBy, "CreatedBy", node, tt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}
	auxDeletedBy, err := e.encodeRef(p, res.DeletedBy, "DeletedBy", node, tt)
	if err != nil {
		return
	}

	auxOwnedBy, err := e.encodeRef(p, res.OwnedBy, "OwnedBy", node, tt)
	if err != nil {
		return
	}

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}
	auxUpdatedBy, err := e.encodeRef(p, res.UpdatedBy, "UpdatedBy", node, tt)
	if err != nil {
		return
	}
	auxWorkflowID, err := e.encodeRef(p, res.WorkflowID, "WorkflowID", node, tt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"constraints", res.Constraints,
		"createdAt", auxCreatedAt,
		"createdBy", auxCreatedBy,
		"deletedAt", auxDeletedAt,
		"deletedBy", auxDeletedBy,
		"enabled", res.Enabled,
		"eventType", res.EventType,
		"id", res.ID,
		"input", res.Input,
		"meta", res.Meta,
		"ownedBy", auxOwnedBy,
		"resourceType", res.ResourceType,
		"stepID", res.StepID,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
		"workflowID", auxWorkflowID,
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
