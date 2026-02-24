package service

import (
	"context"
	"regexp"
	"strconv"

	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/label"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	connection struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        connectionAccessController
	}

	connectionAccessController interface {
		CanSearchConnections(ctx context.Context) bool

		CanCreateConnection(context.Context) bool
		CanReadConnection(context.Context, *types.Connection) bool
		CanUpdateConnection(context.Context, *types.Connection) bool
		CanDeleteConnection(context.Context, *types.Connection) bool
	}

	ConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Connection, error)
		Create(ctx context.Context, new *types.Connection) (*types.Connection, error)
		Update(ctx context.Context, upd *types.Connection) (*types.Connection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ConnectionFilter) (types.ConnectionSet, types.ConnectionFilter, error)
	}
)

func Connection() *connection {
	return &connection{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *connection) FindByID(ctx context.Context, ID uint64) (res *types.Connection, err error) {
	res, err = loadConnection(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadConnection(ctx, res) {
		return nil, ConnectionErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, res); err != nil {
		return nil, err
	}

	svc.deriveParams(res)
	return
}

func (svc *connection) Create(ctx context.Context, new *types.Connection) (res *types.Connection, err error) {
	if !svc.ac.CanCreateConnection(ctx) {
		return nil, ConnectionErrNotAllowedToCreate()
	}

	if err = svc.validateConnection(new); err != nil {
		return nil, err
	}

	if err = svc.uniqueHandleCheck(ctx, new.Handle, 0); err != nil {
		return nil, err
	}

	new.ID = nextID()
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	new.Revision = 1

	if new.Status == "" {
		new.Status = "draft"
	}

	if err = store.CreateConnection(ctx, svc.store, new); err != nil {
		return nil, err
	}

	if err = label.Create(ctx, svc.store, new); err != nil {
		return nil, err
	}

	return new, nil
}

func (svc *connection) Update(ctx context.Context, upd *types.Connection) (res *types.Connection, err error) {
	res, err = loadConnection(ctx, svc.store, upd.ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanUpdateConnection(ctx, res) {
		return nil, ConnectionErrNotAllowedToUpdate()
	}

	if err = svc.validateConnection(upd); err != nil {
		return nil, err
	}

	if upd.Handle != res.Handle {
		if err = svc.uniqueHandleCheck(ctx, upd.Handle, upd.ID); err != nil {
			return nil, err
		}
	}

	res.Handle = upd.Handle
	res.Meta = upd.Meta
	res.Service = upd.Service
	res.Resources = upd.Resources
	res.StandardOperations = upd.StandardOperations
	res.Operations = upd.Operations
	res.Revision++

	n := now()
	res.UpdatedAt = n
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateConnection(ctx, svc.store, res); err != nil {
		return nil, err
	}

	if label.Changed(res.Labels, upd.Labels) {
		if err = label.Update(ctx, svc.store, upd); err != nil {
			return nil, err
		}
		res.Labels = upd.Labels
	}

	return
}

func (svc *connection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	res, err := loadConnection(ctx, svc.store, ID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteConnection(ctx, res) {
		return ConnectionErrNotAllowedToDelete()
	}

	// Fail if active connections exist
	cc, _, err := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
		ConnectionID: ID,
	})
	if err != nil {
		return err
	}
	if len(cc) > 0 {
		return ConnectionErrHasActiveConnections()
	}

	n := now()
	res.DeletedAt = n
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	return store.UpdateConnection(ctx, svc.store, res)
}

func (svc *connection) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	res, err := loadConnection(ctx, svc.store, ID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteConnection(ctx, res) {
		return ConnectionErrNotAllowedToUndelete()
	}

	res.DeletedAt = nil
	res.DeletedBy = 0

	return store.UpdateConnection(ctx, svc.store, res)
}

func (svc *connection) Search(ctx context.Context, filter types.ConnectionFilter) (set types.ConnectionSet, f types.ConnectionFilter, err error) {
	if !svc.ac.CanSearchConnections(ctx) {
		return nil, f, ConnectionErrNotAllowedToSearch()
	}

	filter.Check = func(res *types.Connection) (bool, error) {
		if !svc.ac.CanReadConnection(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	set, f, err = store.SearchConnections(ctx, svc.store, filter)
	if err != nil {
		return nil, f, err
	}

	if err = label.Load(ctx, svc.store, toLabeledConnections(set)...); err != nil {
		return nil, f, err
	}

	for _, c := range set {
		svc.deriveParams(c)
	}
	return
}

func (svc *connection) Install(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	// if !svc.ac.CanCreateConfiguredConnection(ctx) {
	// 	return nil, errors.Unauthorized("connection install denied")
	// }

	// 1. Fetch and snapshot the connection definition
	var conn *types.Connection
	conn, err = svc.FindByID(ctx, new.ConnectionID)
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
func (svc *connection) dispatch(ctx context.Context, resolved *types.Connection, conn *types.ConfiguredConnection) error {
	// dalConn, err := svc.provisionDAL(ctx, resolved.Service, conn.Config.CredentialID)
	// modules, err := svc.provisionModules(ctx, resolved.Resources, conn.Config.NamespaceID, dalConn.ID)
	// err = svc.provisionAutomation(ctx, resolved.Operations, modules, dalConn.ID)
	// err = svc.provisionWebhooks(ctx, resolved.Resources, modules)
	return nil
}

// resolveTemplates substitutes all {{placeholder}} variables in a connection's
// Template fields with values from the connection's scoped params.
// Returns a deep copy of the connection with all templates resolved.
func (svc *connection) resolveTemplates(conn *types.Connection, params []types.ConfiguredConnectionParam) *types.Connection {
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

// -- validation ---------------------------------------------------------------

func (svc *connection) validateConnection(c *types.Connection) error {
	if c.Meta.Short == "" {
		return ConnectionErrMissingShortName()
	}
	if c.Service.BaseURL.Value == "" {
		return ConnectionErrMissingBaseURL()
	}
	if c.Service.Auth.Method == "" {
		return ConnectionErrMissingAuthMethod()
	}
	return nil
}

func (svc *connection) uniqueHandleCheck(ctx context.Context, handle string, excludeID uint64) error {
	if handle == "" {
		return nil
	}
	set, _, err := store.SearchConnections(ctx, svc.store, types.ConnectionFilter{Handle: handle})
	if err != nil {
		return err
	}
	for _, c := range set {
		if c.ID != excludeID {
			return ConnectionErrHandleNotUnique()
		}
	}
	return nil
}

// -- parameter derivation -----------------------------------------------------

var placeholderRe = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// deriveParams scans all Template fields in a Connection, extracts {{placeholder}}
// variables, and populates the read-only DerivedParams field.
func (svc *connection) deriveParams(c *types.Connection) {
	type scopedTemplate struct {
		scope []string
		tpl   types.ConnectionTemplate
	}

	var tt []scopedTemplate

	// Service-level templates
	tt = append(tt, scopedTemplate{[]string{"service", "baseURL"}, c.Service.BaseURL})
	for _, h := range c.Service.Headers {
		tt = append(tt, scopedTemplate{[]string{"service", "headers"}, h})
	}
	for _, p := range c.Service.Auth.Params {
		tt = append(tt, scopedTemplate{[]string{"service", "auth"}, p})
	}

	// Resource-level templates
	for _, r := range c.Resources {
		tt = append(tt, scopedTemplate{[]string{"resources", r.Handle, "endpoint"}, r.Endpoint})
	}

	// Standard operations
	collectStdOp := func(name string, op *types.ConnectionHTTPAction) {
		if op == nil {
			return
		}
		for _, t := range collectHTTPActionTemplates(op) {
			tt = append(tt, scopedTemplate{[]string{"standardOperations", name}, t})
		}
	}
	collectStdOp("list", c.StandardOperations.List)
	collectStdOp("read", c.StandardOperations.Read)
	collectStdOp("create", c.StandardOperations.Create)
	collectStdOp("update", c.StandardOperations.Update)
	collectStdOp("delete", c.StandardOperations.Delete)

	// Custom operations
	for _, op := range c.Operations {
		for _, t := range collectHTTPActionTemplates(&op.HTTP) {
			tt = append(tt, scopedTemplate{[]string{"operations", op.Handle}, t})
		}
	}

	// Extract unique params (keyed by scope+name)
	seen := make(map[string]bool)
	var params []types.ConnectionDerivedParam

	for _, st := range tt {
		// Build a lookup of inline placeholder metadata
		meta := make(map[string]types.ConnectionPlaceholder, len(st.tpl.Placeholders))
		for _, p := range st.tpl.Placeholders {
			meta[p.Name] = p
		}

		for _, match := range placeholderRe.FindAllStringSubmatch(st.tpl.Value, -1) {
			name := match[1]
			key := name + "|" + joinScope(st.scope)
			if seen[key] {
				continue
			}
			seen[key] = true

			p := types.ConnectionDerivedParam{
				Name:     name,
				Scope:    st.scope,
				Type:     "string",
				Required: true,
			}
			if m, ok := meta[name]; ok {
				if m.Type != "" {
					p.Type = m.Type
				}
				p.Description = m.Description
				p.Required = m.Required
				p.Default = m.Default
				p.Options = m.Options
			}
			params = append(params, p)
		}
	}

	c.DerivedParams = params
}

func collectHTTPActionTemplates(a *types.ConnectionHTTPAction) []types.ConnectionTemplate {
	tt := []types.ConnectionTemplate{a.Path, a.BodyTemplate}
	for _, h := range a.Headers {
		tt = append(tt, h)
	}
	for _, q := range a.QueryParams {
		tt = append(tt, q)
	}
	return tt
}

func joinScope(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += "."
		}
		out += s
	}
	return out
}

// -- helpers ------------------------------------------------------------------

func loadConnection(ctx context.Context, s store.Connections, ID uint64) (res *types.Connection, err error) {
	if ID == 0 {
		return nil, ConnectionErrInvalidID()
	}

	if res, err = store.LookupConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		err = ConnectionErrNotFound()
	}

	return
}

func toLabeledConnections(set types.ConnectionSet) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
