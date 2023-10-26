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

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
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
	case types.ApplicationResourceType:
		aux, err = e.encodeApplications(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "application", aux)
		if err != nil {
			return
		}
	case types.ApigwRouteResourceType:
		aux, err = e.encodeApigwRoutes(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "apigwRoute", aux)
		if err != nil {
			return
		}

	case types.AuthClientResourceType:
		aux, err = e.encodeAuthClients(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "authClient", aux)
		if err != nil {
			return
		}

	case types.QueueResourceType:
		aux, err = e.encodeQueues(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "queue", aux)
		if err != nil {
			return
		}

	case types.ReportResourceType:
		aux, err = e.encodeReports(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "report", aux)
		if err != nil {
			return
		}

	case types.RoleResourceType:
		aux, err = e.encodeRoles(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "role", aux)
		if err != nil {
			return
		}

	case types.TemplateResourceType:
		aux, err = e.encodeTemplates(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "template", aux)
		if err != nil {
			return
		}
	case types.UserResourceType:
		aux, err = e.encodeUsers(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "user", aux)
		if err != nil {
			return
		}
	case types.DalConnectionResourceType:
		aux, err = e.encodeDalConnections(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "dalConnection", aux)
		if err != nil {
			return
		}
	case types.DalSensitivityLevelResourceType:
		aux, err = e.encodeDalSensitivityLevels(ctx, p, nodes, tt)
		if err != nil {
			return
		}
		// Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "dalSensitivityLevel", aux)
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
// Functions for resource application
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeApplications(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeApplication(ctx, p, n, tt)
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

// encodeApplication focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeApplication(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Application)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxOwnerID, err := e.encodeRef(p, res.OwnerID, "OwnerID", node, tt)
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
		"name", res.Name,
		"ownerID", auxOwnerID,
		"unify", res.Unify,
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
// Functions for resource apigwRoute
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeApigwRoutes(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeApigwRoute(ctx, p, n, tt)
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

// encodeApigwRoute focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeApigwRoute(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.ApigwRoute)

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

	auxGroup, err := e.encodeRef(p, res.Group, "Group", node, tt)
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
		"endpoint", res.Endpoint,
		"group", auxGroup,
		"id", res.ID,
		"meta", res.Meta,
		"method", res.Method,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux

	aux, err = e.encodeApigwFilters(ctx, p, tt.ChildrenForResourceType(node, types.ApigwFilterResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"filters", aux,
	)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource apigwFilter
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeApigwFilters(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeApigwFilter(ctx, p, n, tt)
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

// encodeApigwFilter focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeApigwFilter(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.ApigwFilter)

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

	auxRoute, err := e.encodeRef(p, res.Route, "Route", node, tt)
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
		"id", res.ID,
		"kind", res.Kind,
		"params", res.Params,
		"ref", res.Ref,
		"route", auxRoute,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
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
// Functions for resource authClient
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeAuthClients(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeAuthClient(ctx, p, n, tt)
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

// encodeAuthClient focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeAuthClient(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.AuthClient)

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

	auxExpiresAt, err := e.encodeTimestampNil(p, res.ExpiresAt)
	if err != nil {
		return
	}

	auxOwnedBy, err := e.encodeRef(p, res.OwnedBy, "OwnedBy", node, tt)
	if err != nil {
		return
	}

	auxSecurity, err := e.encodeAuthClientSecurityC(ctx, p, tt, node, res, res.Security)
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
	auxValidFrom, err := e.encodeTimestampNil(p, res.ValidFrom)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"createdAt", auxCreatedAt,
		"createdBy", auxCreatedBy,
		"deletedAt", auxDeletedAt,
		"deletedBy", auxDeletedBy,
		"enabled", res.Enabled,
		"expiresAt", auxExpiresAt,
		"handle", res.Handle,
		"id", res.ID,
		"meta", res.Meta,
		"ownedBy", auxOwnedBy,
		"redirectURI", res.RedirectURI,
		"scope", res.Scope,
		"secret", res.Secret,
		"security", auxSecurity,
		"trusted", res.Trusted,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
		"validFrom", auxValidFrom,
		"validGrant", res.ValidGrant,
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
// Functions for resource queue
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeQueues(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeQueue(ctx, p, n, tt)
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

// encodeQueue focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeQueue(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Queue)

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

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}
	auxUpdatedBy, err := e.encodeRef(p, res.UpdatedBy, "UpdatedBy", node, tt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"consumer", res.Consumer,
		"createdAt", auxCreatedAt,
		"createdBy", auxCreatedBy,
		"deletedAt", auxDeletedAt,
		"deletedBy", auxDeletedBy,
		"id", res.ID,
		"meta", res.Meta,
		"queue", res.Queue,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
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
// Functions for resource report
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeReports(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeReport(ctx, p, n, tt)
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

// encodeReport focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeReport(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Report)

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

	out, err = y7s.AddMap(out,
		"blocks", res.Blocks,
		"createdAt", auxCreatedAt,
		"createdBy", auxCreatedBy,
		"deletedAt", auxDeletedAt,
		"deletedBy", auxDeletedBy,
		"handle", res.Handle,
		"id", res.ID,
		"meta", res.Meta,
		"ownedBy", auxOwnedBy,
		"scenarios", res.Scenarios,
		"sources", res.Sources,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
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
// Functions for resource role
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeRoles(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeRole(ctx, p, n, tt)
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

// encodeRole focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeRole(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Role)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxArchivedAt, err := e.encodeTimestampNil(p, res.ArchivedAt)
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

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"archivedAt", auxArchivedAt,
		"createdAt", auxCreatedAt,
		"deletedAt", auxDeletedAt,
		"handle", res.Handle,
		"id", res.ID,
		"meta", res.Meta,
		"name", res.Name,
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
// Functions for resource template
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeTemplates(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeTemplate(ctx, p, n, tt)
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

// encodeTemplate focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeTemplate(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.Template)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxLastUsedAt, err := e.encodeTimestampNil(p, res.LastUsedAt)
	if err != nil {
		return
	}

	auxOwnerID, err := e.encodeRef(p, res.OwnerID, "OwnerID", node, tt)
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
		"handle", res.Handle,
		"id", res.ID,
		"language", res.Language,
		"lastUsedAt", auxLastUsedAt,
		"meta", res.Meta,
		"ownerID", auxOwnerID,
		"partial", res.Partial,
		"template", res.Template,
		"type", res.Type,
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
// Functions for resource user
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeUsers(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeUser(ctx, p, n, tt)
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

// encodeUser focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeUser(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.User)

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	auxCreatedAt, err := e.encodeTimestamp(p, res.CreatedAt)
	if err != nil {
		return
	}
	auxDeletedAt, err := e.encodeTimestampNil(p, res.DeletedAt)
	if err != nil {
		return
	}

	auxSuspendedAt, err := e.encodeTimestampNil(p, res.SuspendedAt)
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
		"email", res.Email,
		"emailConfirmed", res.EmailConfirmed,
		"handle", res.Handle,
		"id", res.ID,
		"kind", res.Kind,
		"meta", res.Meta,
		"name", res.Name,
		"suspendedAt", auxSuspendedAt,
		"updatedAt", auxUpdatedAt,
		"username", res.Username,
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
// Functions for resource dalConnection
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeDalConnections(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeDalConnection(ctx, p, n, tt)
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

// encodeDalConnection focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeDalConnection(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.DalConnection)

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

	auxUpdatedAt, err := e.encodeTimestampNil(p, res.UpdatedAt)
	if err != nil {
		return
	}
	auxUpdatedBy, err := e.encodeRef(p, res.UpdatedBy, "UpdatedBy", node, tt)
	if err != nil {
		return
	}

	out, err = y7s.AddMap(out,
		"config", res.Config,
		"createdAt", auxCreatedAt,
		"createdBy", auxCreatedBy,
		"deletedAt", auxDeletedAt,
		"deletedBy", auxDeletedBy,
		"handle", res.Handle,
		"id", res.ID,
		"meta", res.Meta,
		"type", res.Type,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
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
// Functions for resource dalSensitivityLevel
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeDalSensitivityLevels(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encodeDalSensitivityLevel(ctx, p, n, tt)
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

// encodeDalSensitivityLevel focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encodeDalSensitivityLevel(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.DalSensitivityLevel)

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
		"handle", res.Handle,
		"id", res.ID,
		"level", res.Level,
		"meta", res.Meta,
		"updatedAt", auxUpdatedAt,
		"updatedBy", auxUpdatedBy,
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
