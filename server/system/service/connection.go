package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
	}
}

func (svc *connection) WithConfiguredConnection(cc *configuredConnection) *connection {
	svc.configuredConnection = cc
	return svc
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

		svc.preprocess(new)

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

		svc.deriveParams(new)

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

		svc.preprocess(upd)

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

		svc.deriveParams(res)

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
	res, err = svc.configuredConnection.Update(ctx, upd)
	if err != nil {
		return nil, err
	}

	return
}

// -- preprocessing -----------------------------------------------------------

func (svc *connection) preprocess(c *types.Connection) {
	svc.normaliseOperationTypes(c)
}

func (svc *connection) normaliseOperationTypes(c *types.Connection) {
	for i := range c.Operations {
		for j := range c.Operations[i].Input {
			if strings.EqualFold(c.Operations[i].Input[j].Type, "Object") {
				c.Operations[i].Input[j].Type = "Vars"
			}
		}

		for j := range c.Operations[i].Output {
			if strings.EqualFold(c.Operations[i].Output[j].Type, "Object") {
				c.Operations[i].Output[j].Type = "Vars"
			}
		}
	}
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

	for _, op := range c.Operations {
		if len(op.Steps) > 1 {
			return fmt.Errorf("multi step operations are currently not supoorted")
		}
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
				Label:    labelFromName(name),
				Scope:    st.scope,
				Type:     "string",
				Required: true,
			}
			if m, ok := meta[name]; ok {
				if m.Type != "" {
					p.Type = m.Type
				}
				if m.Label != "" {
					p.Label = m.Label
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

// labelFromName converts any common identifier style into a human-friendly
// title-cased string. Supported input styles:
//
//	camelCase        → "Camel Case"
//	PascalCase       → "Pascal Case"
//	snake_case       → "Snake Case"
//	kebab-case       → "Kebab Case"
//	SCREAMING_SNAKE  → "Screaming Snake"
//	APIKey / apiID   → "API Key" / "Api ID"  (consecutive-uppercase runs kept as words)
func labelFromName(name string) string {
	runes := []rune(name)
	n := len(runes)
	if n == 0 {
		return ""
	}

	// Split into words first
	var words []string
	start := 0

	upper := func(r rune) bool { return r >= 'A' && r <= 'Z' }
	lower := func(r rune) bool { return r >= 'a' && r <= 'z' }
	sep   := func(r rune) bool { return r == '_' || r == '-' }

	flush := func(end int) {
		if end > start {
			words = append(words, string(runes[start:end]))
		}
	}

	for i := 1; i < n; i++ {
		prev, cur := runes[i-1], runes[i]
		switch {
		case sep(cur):
			flush(i)
			start = i + 1
		case sep(prev):
			start = i
		case upper(cur) && lower(prev):
			// fooBar → foo | Bar
			flush(i)
			start = i
		case upper(cur) && i+1 < n && lower(runes[i+1]) && upper(prev):
			// APIKey → API | Key
			flush(i)
			start = i
		}
	}
	flush(n)

	// Title-case each word; preserve all-uppercase words (acronyms) as-is
	var b strings.Builder
	for i, w := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		rr := []rune(w)
		allUpper := true
		for _, r := range rr {
			if r >= 'a' && r <= 'z' {
				allUpper = false
				break
			}
		}
		if allUpper && len(rr) > 1 {
			// acronym — keep as-is (e.g. API, ID)
			b.WriteString(w)
		} else {
			// title-case: uppercase first, lowercase rest
			for j, r := range rr {
				if j == 0 && r >= 'a' && r <= 'z' {
					r -= 32
				} else if j > 0 && r >= 'A' && r <= 'Z' {
					r += 32
				}
				b.WriteRune(r)
			}
		}
	}
	return b.String()
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
