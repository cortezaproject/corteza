package service

import (
	"context"
	"regexp"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/label"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	connection struct {
		actionlog            actionlog.Recorder
		store                store.Storer
		ac                   connectionAccessController
		configuredConnection *configuredConnection
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

	ConfiguredConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
		Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (types.ConfiguredConnectionSet, types.ConfiguredConnectionFilter, error)
		DeleteByID(ctx context.Context, ID uint64) error
		Enable(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
	}

	dispatchRsp struct {
		dalConnection *types.DalConnection
	}
)

func Connection() *connection {
	return &connection{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,

		configuredConnection: DefaultConfiguredConnection,
	}
}

func (svc *connection) FindByID(ctx context.Context, ID uint64) (res *types.Connection, err error) {
	var (
		aProps = &connectionActionProps{connection: &types.Connection{ID: ID}}
	)

	err = func() error {
		res, err = loadConnection(ctx, svc.store, ID)
		if err != nil {
			return err
		}

		aProps.setConnection(res)

		if !svc.ac.CanReadConnection(ctx, res) {
			return ConnectionErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		svc.deriveParams(res)
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConnectionActionLookup, err)
}

func (svc *connection) Create(ctx context.Context, new *types.Connection) (res *types.Connection, err error) {
	var (
		aProps = &connectionActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateConnection(ctx) {
			return ConnectionErrNotAllowedToCreate()
		}

		if err = svc.validateConnection(new); err != nil {
			return err
		}

		if err = svc.uniqueHandleCheck(ctx, new.Handle, 0); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
		new.Revision = 1

		if new.Status == "" {
			new.Status = "draft"
		}

		if err = store.CreateConnection(ctx, svc.store, new); err != nil {
			return err
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return err
		}

		return nil
	}()

	return new, svc.recordAction(ctx, aProps, ConnectionActionCreate, err)
}

func (svc *connection) Update(ctx context.Context, upd *types.Connection) (res *types.Connection, err error) {
	var (
		aProps = &connectionActionProps{update: upd}
	)

	err = func() (err error) {
		if res, err = loadConnection(ctx, svc.store, upd.ID); err != nil {
			return err
		}

		aProps.setConnection(res)

		if !svc.ac.CanUpdateConnection(ctx, res) {
			return ConnectionErrNotAllowedToUpdate()
		}

		if err = svc.validateConnection(upd); err != nil {
			return err
		}

		if upd.Handle != res.Handle {
			if err = svc.uniqueHandleCheck(ctx, upd.Handle, upd.ID); err != nil {
				return err
			}
		}

		res.Handle = upd.Handle
		res.Meta = upd.Meta
		res.Service = upd.Service
		res.Resources = upd.Resources
		res.Operations = upd.Operations
		res.Revision++

		n := now()
		res.UpdatedAt = n
		res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConnection(ctx, svc.store, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return err
			}
			res.Labels = upd.Labels
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConnectionActionUpdate, err)
}

func (svc *connection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		res    *types.Connection
		aProps = &connectionActionProps{connection: &types.Connection{ID: ID}}
	)

	err = func() (err error) {
		if res, err = loadConnection(ctx, svc.store, ID); err != nil {
			return err
		}

		aProps.setConnection(res)

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
	}()

	return svc.recordAction(ctx, aProps, ConnectionActionDelete, err)
}

func (svc *connection) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		res    *types.Connection
		aProps = &connectionActionProps{connection: &types.Connection{ID: ID}}
	)

	err = func() (err error) {
		if res, err = loadConnection(ctx, svc.store, ID); err != nil {
			return err
		}

		aProps.setConnection(res)

		if !svc.ac.CanDeleteConnection(ctx, res) {
			return ConnectionErrNotAllowedToUndelete()
		}

		res.DeletedAt = nil
		res.DeletedBy = 0

		return store.UpdateConnection(ctx, svc.store, res)
	}()

	return svc.recordAction(ctx, aProps, ConnectionActionUndelete, err)
}

func (svc *connection) Search(ctx context.Context, filter types.ConnectionFilter) (set types.ConnectionSet, f types.ConnectionFilter, err error) {
	var (
		aProps = &connectionActionProps{filter: &filter}
	)

	filter.Check = func(res *types.Connection) (bool, error) {
		if !svc.ac.CanReadConnection(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchConnections(ctx) {
			return ConnectionErrNotAllowedToSearch()
		}

		set, f, err = store.SearchConnections(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledConnections(set)...); err != nil {
			return err
		}

		for _, c := range set {
			svc.deriveParams(c)
		}
		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ConnectionActionSearch, err)
}

// Configure creates a new ConfiguredConnection in draft status.
// No provisioning is done yet; config can be freely edited.
func (svc *connection) Configure(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	return svc.configuredConnection.Create(ctx, new)
}

func (svc *connection) UpdateConfiguration(ctx context.Context, upd *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	upd, err = svc.configuredConnection.Update(ctx, upd)
	if err != nil {
		return nil, err
	}

	return
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

	// Standard operations
	collectStdOp := func(name string, op *types.ConnectionHTTPAction) {
		if op == nil {
			return
		}
		for _, t := range collectHTTPActionTemplates(op) {
			tt = append(tt, scopedTemplate{[]string{"standardOperations", name}, t})
		}
	}

	// Resource-level templates
	for _, r := range c.Resources {
		tt = append(tt, scopedTemplate{[]string{"resources", r.Handle, "endpoint"}, r.Endpoint})

		collectStdOp("list", r.Operations.List)
		collectStdOp("read", r.Operations.Read)
		collectStdOp("create", r.Operations.Create)
		collectStdOp("update", r.Operations.Update)
		collectStdOp("delete", r.Operations.Delete)
	}

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
