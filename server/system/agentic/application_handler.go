package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in application_tools.go, in the same order.
type applicationHandler struct {
	reg toolRegistrar
}

func ApplicationHandler(reg toolRegistrar) *applicationHandler {
	h := &applicationHandler{reg: reg}
	h.register()
	return h
}

// applicationItem is the compact projection returned by list mode. Single-item
// lookups return the service type unchanged; lists must not, because the Unify
// block — URL, config string, icon and logo — is bulky and tells the caller
// nothing useful about which application to pick.
type applicationItem struct {
	ApplicationID string     `json:"applicationID"`
	Name          string     `json:"name"`
	Enabled       bool       `json:"enabled"`
	Weight        int        `json:"weight"`
	Flags         []string   `json:"flags,omitempty"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

func (h *applicationHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	deleted := filter.StateExcluded
	if toolkit.Bool(args, "includeDeleted") {
		deleted = filter.StateInclusive
	}

	// Reference param is never Required, so an empty ref means "list".
	ref, err := toolkit.Ref(args, "application")
	if err != nil {
		return nil, err
	}
	if ref != "" {
		app, err := resolveApplication(ctx, ref, deleted)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(app)
	}

	f := sysTypes.ApplicationFilter{
		Query:   toolkit.Str(args, "query"),
		Deleted: deleted,
	}

	if raw, ok := args["flags"]; ok && raw != nil {
		if f.Flags, err = applicationStringList(raw); err != nil {
			return nil, fmt.Errorf("invalid flags: must be a JSON array of strings: %w", err)
		}
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultApplication.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("application list", err)
	}

	items := make([]applicationItem, 0, len(set))
	for _, app := range set {
		items = append(items, applicationItem{
			ApplicationID: strconv.FormatUint(app.ID, 10),
			Name:          app.Name,
			Enabled:       app.Enabled,
			Weight:        app.Weight,
			Flags:         app.Flags,
			DeletedAt:     app.DeletedAt,
		})
	}

	// Pass the cursor itself, not its String(): String is a human-readable debug
	// rendering ("<id: 123, [FWD]>") while parseCursor expects base64 of the
	// cursor JSON, which is what MarshalJSON emits. json.Marshal renders a nil
	// pointer as null, so the last page is safe.
	return toolkit.JSONResult(map[string]any{
		"applications":   items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *applicationHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	app := &sysTypes.Application{
		Name:    name,
		Enabled: toolkit.Bool(args, "enabled"),
	}

	// The service fills in an empty Unify block when none is given, so merging
	// onto a fresh struct is the same shape create and update both use.
	if raw, ok := args["unify"]; ok && raw != nil {
		app.Unify = &sysTypes.ApplicationUnify{}
		if err = mergeApplicationUnify(app.Unify, raw); err != nil {
			return nil, err
		}
	}

	app, err = sysService.DefaultApplication.Create(ctx, app)
	if err != nil {
		return nil, toolkit.Errf("application creation", err)
	}
	return toolkit.JSONResult(app)
}

func (h *applicationHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "application")
	if err != nil {
		return nil, err
	}

	app, err := resolveApplication(ctx, ref, filter.StateExcluded)
	if err != nil {
		return nil, err
	}

	// Absent leaves the field alone. Name is the one field that cannot be
	// cleared: the service does not validate it and a nameless application is
	// unfindable in every listing that shows one.
	if v, ok := args["name"]; ok {
		s, _ := v.(string)
		if s == "" {
			return nil, fmt.Errorf("name cannot be set to an empty string; omit it to leave it unchanged")
		}
		app.Name = s
	}
	if _, ok := args["enabled"]; ok {
		app.Enabled = toolkit.Bool(args, "enabled")
	}

	// Merged key by key rather than replaced wholesale: every field of the block
	// is a scalar, so an absent key means "unchanged" and a present empty one
	// means "clear", with no collection that would need a remove companion.
	if raw, ok := args["unify"]; ok && raw != nil {
		if app.Unify == nil {
			app.Unify = &sysTypes.ApplicationUnify{}
		}
		if err = mergeApplicationUnify(app.Unify, raw); err != nil {
			return nil, err
		}
	}

	app, err = sysService.DefaultApplication.Update(ctx, app)
	if err != nil {
		return nil, toolkit.Errf("application update", err)
	}
	return toolkit.JSONResult(app)
}

func (h *applicationHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "application")
	if err != nil {
		return nil, err
	}

	app, err := resolveApplication(ctx, ref, filter.StateExcluded)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultApplication.DeleteByID(ctx, app.ID); err != nil {
		return nil, toolkit.Errf("application deletion", err)
	}
	return toolkit.TextResult("application %d deleted", app.ID), nil
}

func (h *applicationHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// ID only, deliberately: ApplicationFilter.Deleted defaults to excluding
	// deleted records, so resolving a name would never find the application this
	// tool exists to restore.
	id, err := toolkit.ReqID(args, "applicationID")
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultApplication.UndeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("application undeletion", err)
	}
	return toolkit.TextResult("application %d restored", id), nil
}

func (h *applicationHandler) reorder(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	raw, ok := args["order"]
	if !ok || raw == nil {
		return nil, fmt.Errorf("order is required")
	}

	// applicationStringList refuses a JSON array of numbers, which is the point:
	// an ID that arrived as a number has already lost precision.
	idStrs, err := applicationStringList(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid order: must be a JSON array of ID strings: %w", err)
	}
	if len(idStrs) == 0 {
		return nil, fmt.Errorf("order must list at least one application")
	}

	order := make([]uint64, len(idStrs))
	for i, s := range idStrs {
		if order[i], err = strconv.ParseUint(s, 10, 64); err != nil {
			return nil, fmt.Errorf("invalid application ID %q: %w", s, err)
		}
	}

	if err = sysService.DefaultApplication.Reorder(ctx, order); err != nil {
		return nil, toolkit.Errf("application reorder", err)
	}
	return toolkit.TextResult("%d applications reordered", len(order)), nil
}

func (h *applicationHandler) flag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app, ownedBy, name, err := applicationFlagArgs(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultApplication.Flag(ctx, app, ownedBy, name); err != nil {
		return nil, toolkit.Errf("application flagging", err)
	}
	return toolkit.TextResult("application %d flagged %q", app.ID, name), nil
}

func (h *applicationHandler) unflag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app, ownedBy, name, err := applicationFlagArgs(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultApplication.Unflag(ctx, app, ownedBy, name); err != nil {
		return nil, toolkit.Errf("application unflagging", err)
	}
	return toolkit.TextResult("flag %q removed from application %d", name, app.ID), nil
}

// resolveApplication turns a name-or-ID reference into an application.
//
// The application service has neither FindByAny nor FindByHandle, so this is
// written out here. A name is matched exactly through ApplicationFilter.Name;
// ApplicationFilter.Query is deliberately not used as a fallback, because it is
// an ILIKE substring match and would resolve "Reports" to "Sales Reports"
// without the caller ever learning that a different application was touched.
func resolveApplication(ctx context.Context, ref string, deleted filter.State) (*sysTypes.Application, error) {
	if ref == "" {
		return nil, fmt.Errorf("application is required")
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		app, err := sysService.DefaultApplication.FindByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("application lookup", err)
		}
		return app, nil
	}

	f := sysTypes.ApplicationFilter{Name: ref, Deleted: deleted}
	// Small enough to stay cheap, large enough that "several" is reported as
	// several rather than as one arbitrary match.
	f.Limit = 10

	set, _, err := sysService.DefaultApplication.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("application lookup", err)
	}

	switch len(set) {
	case 0:
		return nil, fmt.Errorf("no application named %q; names are matched exactly, use system_application_lookup with 'query' to search by fragment", ref)
	case 1:
		return set[0], nil
	default:
		ids := make([]string, 0, len(set))
		for _, app := range set {
			ids = append(ids, strconv.FormatUint(app.ID, 10))
		}
		return nil, fmt.Errorf("%d applications are named %q (IDs %v); pass one of those IDs instead of the name", len(set), ref, ids)
	}
}

// applicationFlagArgs unwraps the three arguments the flag operations share.
//
// The owner is derived from an explicit 'mode' and never inferred from an ID.
// The service branches on ownedBy — zero is a global flag gated by the
// global-flag permission, non-zero must equal the caller's own identity and is
// gated by the self-flag permission — so a tool that let an ID decide would let
// a caller write shared state while believing it was setting a private one.
func applicationFlagArgs(ctx context.Context, req mcp.CallToolRequest) (*sysTypes.Application, uint64, string, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, 0, "", err
	}

	ref, err := toolkit.ReqRef(args, "application")
	if err != nil {
		return nil, 0, "", err
	}

	name, err := toolkit.ReqStr(args, "flag")
	if err != nil {
		return nil, 0, "", err
	}

	mode, err := toolkit.ReqStr(args, "mode")
	if err != nil {
		return nil, 0, "", err
	}

	var ownedBy uint64
	switch mode {
	case "global":
		// Zero owner is what the service reads as "global".
		ownedBy = 0
	case "own":
		ownedBy = a.GetIdentityFromContext(ctx).Identity()
		if ownedBy == 0 {
			// Falling through with zero would silently become a global flag.
			return nil, 0, "", fmt.Errorf(`mode "own" needs an authenticated identity, and this request has none`)
		}
	default:
		return nil, 0, "", fmt.Errorf(`invalid mode %q: expected "own" or "global"`, mode)
	}

	app, err := resolveApplication(ctx, ref, filter.StateExcluded)
	if err != nil {
		return nil, 0, "", err
	}
	return app, ownedBy, name, nil
}

// mergeApplicationUnify applies a JSON object onto an existing unify block.
//
// Unmarshalling onto the loaded struct is what gives the tool its documented
// semantics for free: a key that is absent leaves the field alone, a key that is
// present overwrites it, and a present empty value clears it.
func mergeApplicationUnify(into *sysTypes.ApplicationUnify, raw any) error {
	var data []byte

	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil
		}
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return fmt.Errorf("cannot encode unify: %w", err)
		}
	}

	if err := json.Unmarshal(data, into); err != nil {
		return fmt.Errorf("invalid unify: must be a JSON object, e.g. {\"listed\":true,\"url\":\"/compose\"}: %w", err)
	}
	return nil
}

// applicationStringList reads a JSON array of strings, accepting it either as a
// string holding JSON or as an already-decoded array. A JSON array of numbers is
// refused rather than coerced, because an ID that arrived as a number has
// already lost precision.
func applicationStringList(raw any) ([]string, error) {
	var data []byte

	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil, nil
		}
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return nil, err
		}
	}

	var out []string
	return out, json.Unmarshal(data, &out)
}
