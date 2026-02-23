package service

import (
	"context"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	configuredConnection struct {
		actionlog     actionlog.Recorder
		store         store.Storer
		ac            configuredConnectionAccessController
		connectionSvc *connection
	}

	configuredConnectionAccessController interface {
		CanSearchConfiguredConnections(ctx context.Context) bool

		CanCreateConfiguredConnection(context.Context) bool
		CanReadConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
		CanDeleteConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
	}
)

func ConfiguredConnectionSvc() *configuredConnection {
	return &configuredConnection{
		actionlog:     DefaultActionlog,
		store:         DefaultStore,
		ac:            DefaultAccessControl,
		connectionSvc: DefaultConnection,
	}
}

func (svc *configuredConnection) FindByID(ctx context.Context, ID uint64) (res *types.ConfiguredConnection, err error) {
	res, err = store.LookupConfiguredConnectionByID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadConfiguredConnection(ctx, res) {
		return nil, errors.Unauthorized("connection read denied")
	}

	return
}

func (svc *configuredConnection) Install(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	if !svc.ac.CanCreateConfiguredConnection(ctx) {
		return nil, errors.Unauthorized("connection install denied")
	}

	// 1. Fetch and snapshot the connection definition
	var conn *types.Connection
	conn, err = svc.connectionSvc.FindByID(ctx, new.ConnectionID)
	if err != nil {
		return nil, err
	}

	// Resolve templates and dispatch sub-system provisioning
	resolved := svc.resolveTemplates(conn, new.Config.Params)
	if err = svc.dispatch(ctx, resolved, new); err != nil {
		return nil, err
	}

	// 7. Store Connection record
	new.ID = nextID()
	new.Connection = *conn
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	if new.Status == "" {
		new.Status = "active"
	}

	// Attach ownership labels
	if new.Labels == nil {
		new.Labels = make(map[string]labelTypes.LabelValue)
	}
	new.Labels["corteza/connection-id"] = labelTypes.LabelValue{Val: strconv.FormatUint(conn.ID, 10)}
	new.Labels["corteza/connection-revision"] = labelTypes.LabelValue{Val: strconv.Itoa(conn.Revision)}
	new.Labels["corteza/configured-connection-id"] = labelTypes.LabelValue{Val: strconv.FormatUint(new.ID, 10)}

	if err = store.CreateConfiguredConnection(ctx, svc.store, new); err != nil {
		return nil, err
	}

	return new, nil
}

// dispatch orchestrates sub-system provisioning from a resolved connection definition.
// TODO: implement DAL, Compose, Automation, and IG provisioning.
func (svc *configuredConnection) dispatch(ctx context.Context, resolved *types.Connection, conn *types.ConfiguredConnection) error {
	// dalConn, err := svc.provisionDAL(ctx, resolved.Service, conn.Config.CredentialID)
	// modules, err := svc.provisionModules(ctx, resolved.Resources, conn.Config.NamespaceID, dalConn.ID)
	// err = svc.provisionAutomation(ctx, resolved.Operations, modules, dalConn.ID)
	// err = svc.provisionWebhooks(ctx, resolved.Resources, modules)
	return nil
}

// resolveTemplates substitutes all {{placeholder}} variables in a connection's
// Template fields with values from the connection's scoped params.
// Returns a deep copy of the connection with all templates resolved.
func (svc *configuredConnection) resolveTemplates(conn *types.Connection, params []types.ConfiguredConnectionParam) *types.Connection {
	// Build a quick lookup: "scope|name" → value
	lookup := make(map[string]string, len(params))
	for _, p := range params {
		key := joinScope(p.Scope) + "|" + p.Name
		lookup[key] = p.Value
	}

	resolve := func(tpl *types.ConnectionTemplate, scope []string) {
		scopeKey := joinScope(scope)
		tpl.Value = placeholderRe.ReplaceAllStringFunc(tpl.Value, func(match string) string {
			sub := placeholderRe.FindStringSubmatch(match)
			if len(sub) < 2 {
				return match
			}
			if v, ok := lookup[scopeKey+"|"+sub[1]]; ok {
				return v
			}
			return match
		})
	}

	// Deep-copy connection so we don't mutate the original
	out := *conn

	// Service-level
	resolve(&out.Service.BaseURL, []string{"service", "baseURL"})

	headers := make(map[string]types.ConnectionTemplate, len(out.Service.Headers))
	for k, h := range out.Service.Headers {
		resolve(&h, []string{"service", "headers"})
		headers[k] = h
	}
	out.Service.Headers = headers

	authParams := make(map[string]types.ConnectionTemplate, len(out.Service.Auth.Params))
	for k, p := range out.Service.Auth.Params {
		resolve(&p, []string{"service", "auth"})
		authParams[k] = p
	}
	out.Service.Auth.Params = authParams

	// Resources
	resources := make(types.ConnectionResources, len(out.Resources))
	copy(resources, out.Resources)
	for i := range resources {
		resolve(&resources[i].Endpoint, []string{"resources", resources[i].Handle, "endpoint"})
	}
	out.Resources = resources

	// Standard operations
	resolveHTTPAction := func(op *types.ConnectionHTTPAction, scope []string) {
		if op == nil {
			return
		}
		resolve(&op.Path, scope)
		resolve(&op.BodyTemplate, scope)
		for k, h := range op.Headers {
			resolve(&h, scope)
			op.Headers[k] = h
		}
		for k, q := range op.QueryParams {
			resolve(&q, scope)
			op.QueryParams[k] = q
		}
	}

	resolveHTTPAction(out.StandardOperations.List, []string{"standardOperations", "list"})
	resolveHTTPAction(out.StandardOperations.Read, []string{"standardOperations", "read"})
	resolveHTTPAction(out.StandardOperations.Create, []string{"standardOperations", "create"})
	resolveHTTPAction(out.StandardOperations.Update, []string{"standardOperations", "update"})
	resolveHTTPAction(out.StandardOperations.Delete, []string{"standardOperations", "delete"})

	// Custom operations
	ops := make(types.ConnectionOperations, len(out.Operations))
	copy(ops, out.Operations)
	for i := range ops {
		resolveHTTPAction(&ops[i].HTTP, []string{"operations", ops[i].Handle})
	}
	out.Operations = ops

	return &out
}

func (svc *configuredConnection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	return ConfiguredConnectionErrDeletionNotSupported()
}

func (svc *configuredConnection) Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (set types.ConfiguredConnectionSet, f types.ConfiguredConnectionFilter, err error) {
	if !svc.ac.CanSearchConfiguredConnections(ctx) {
		return nil, f, errors.Unauthorized("connection search denied")
	}

	set, f, err = store.SearchConfiguredConnections(ctx, svc.store, filter)
	return
}

func loadConfiguredConnection(ctx context.Context, s store.ConfiguredConnections, ID uint64) (res *types.ConfiguredConnection, err error) {
	if ID == 0 {
		return nil, errors.NotFound("connection not found")
	}

	if res, err = store.LookupConfiguredConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		err = errors.NotFound("connection not found")
	}

	return
}
