package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in user_group_tools.go, in the same order.
type userGroupHandler struct {
	reg toolRegistrar
}

func UserGroupHandler(reg toolRegistrar) *userGroupHandler {
	h := &userGroupHandler{reg: reg}
	h.register()
	return h
}

// userGroupItem is the compact projection returned by list mode. Single-item
// lookups return the service type unchanged.
//
// Membership is absent on purpose: it is the heavy field on a group and the
// reason system_user_group_member_list exists as its own tool. There is no name
// field to project — types.UserGroup has only a handle, so Meta.Short stands in
// as the human label.
type userGroupItem struct {
	UserGroupID string     `json:"userGroupID"`
	Handle      string     `json:"handle"`
	Short       string     `json:"short"`
	IsRoot      bool       `json:"isRoot"`
	ParentIDs   []string   `json:"parentIDs"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

// userGroupMember is the projection for member_list. The REST controller returns
// bare IDs; handle and name are added because a model picking a person out of a
// list needs something readable. Email and the rest of the user record are left
// out — this tool's authorization is "can read the group", which is a lower bar
// than reading a user's contact detail.
type userGroupMember struct {
	UserID string `json:"userID"`
	Handle string `json:"handle"`
	Name   string `json:"name"`
}

func (h *userGroupHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Archived and deleted are independent states and both default to excluded,
	// in single-fetch mode as well as in list mode. Keeping the two modes on the
	// same defaults is the point: a by-handle fetch that could see a deleted
	// group the list cannot would make "not found" mean two different things.
	deleted := userGroupState(toolkit.Bool(args, "includeDeleted"))
	archived := userGroupState(toolkit.Bool(args, "includeArchived"))

	// Reference param is never Required, so an empty ref means "list".
	ref, err := toolkit.Ref(args, "userGroup")
	if err != nil {
		return nil, err
	} else if ref != "" {
		g, err := h.resolve(ctx, ref, deleted, archived)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(g)
	}

	f := sysTypes.UserGroupFilter{
		Query:    toolkit.Str(args, "query"),
		Deleted:  deleted,
		Archived: archived,
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultUserGroup.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("user group list", err)
	}

	items := make([]userGroupItem, 0, len(set))
	for _, g := range set {
		items = append(items, newUserGroupItem(g))
	}

	// The cursor value itself, never its String(): String is a debug rendering
	// that parseCursor cannot read back, and json.Marshal renders a nil pointer
	// as null, so the last page needs no special case.
	return toolkit.JSONResult(map[string]any{
		"userGroups":     items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *userGroupHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	hdl, err := toolkit.ReqStr(args, "handle")
	if err != nil {
		return nil, err
	}

	g := &sysTypes.UserGroup{
		Handle: hdl,
		Meta: &sysTypes.UserGroupMeta{
			Short:       toolkit.Str(args, "short"),
			Description: toolkit.Str(args, "description"),
		},
		Config: &sysTypes.UserGroupConfig{},
	}

	if g.Config.Paths, err = h.parsePaths(ctx, toolkit.Str(args, "parents")); err != nil {
		return nil, err
	}

	g, err = sysService.DefaultUserGroup.Create(ctx, g)
	if err != nil {
		return nil, toolkit.Errf("user group creation", err)
	}
	return toolkit.JSONResult(g)
}

func (h *userGroupHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "userGroup")
	if err != nil {
		return nil, err
	}

	// Write ops resolve with both states inclusive: the caller named a specific
	// group, so reporting an archived one as missing would be wrong. Whether the
	// operation is allowed on it is the service's call, not the resolver's.
	g, err := h.resolve(ctx, ref, filter.StateInclusive, filter.StateInclusive)
	if err != nil {
		return nil, err
	}

	// The service replaces handle, meta and config wholesale from what it is
	// given, so the update starts as a copy of the stored group and only the
	// arguments present are overwritten. Absent leaves a field alone;
	// present-and-empty clears it.
	upd := g.Clone()

	if v, ok := args["handle"]; ok {
		upd.Handle, _ = v.(string)
	}

	if upd.Meta == nil {
		upd.Meta = &sysTypes.UserGroupMeta{}
	}
	if v, ok := args["short"]; ok {
		upd.Meta.Short, _ = v.(string)
	}
	if v, ok := args["description"]; ok {
		upd.Meta.Description, _ = v.(string)
	}

	// Collections replace, they do not merge.
	if upd.Config == nil {
		upd.Config = &sysTypes.UserGroupConfig{}
	}
	if _, ok := args["parents"]; ok {
		if upd.Config.Paths, err = h.parsePaths(ctx, toolkit.Str(args, "parents")); err != nil {
			return nil, err
		}
	}

	upd, err = sysService.DefaultUserGroup.Update(ctx, upd)
	if err != nil {
		return nil, toolkit.Errf("user group update", err)
	}
	return toolkit.JSONResult(upd)
}

func (h *userGroupHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	g, err := h.refArg(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultUserGroup.DeleteByID(ctx, g.ID); err != nil {
		return nil, toolkit.Errf("user group deletion", err)
	}
	return toolkit.TextResult("user group %s (%d) deleted", g.Handle, g.ID), nil
}

func (h *userGroupHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	g, err := h.refArg(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultUserGroup.UndeleteByID(ctx, g.ID); err != nil {
		return nil, toolkit.Errf("user group undeletion", err)
	}
	return toolkit.TextResult("user group %s (%d) restored", g.Handle, g.ID), nil
}

func (h *userGroupHandler) memberList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	g, err := h.refArg(ctx, req)
	if err != nil {
		return nil, err
	}

	mm, err := sysService.DefaultUserGroup.MemberList(ctx, g.ID)
	if err != nil {
		return nil, toolkit.Errf("user group member list", err)
	}

	// MemberList takes no filter and returns the whole set, so there is no cursor
	// to hand back. The projection plus toolkit.JSONResult's ceiling are what
	// keep a very large group from flooding the caller.
	members := make([]userGroupMember, 0, len(mm))
	for _, m := range mm {
		members = append(members, userGroupMember{
			UserID: strconv.FormatUint(m.ID, 10),
			Handle: m.Handle,
			Name:   m.Name,
		})
	}

	return toolkit.JSONResult(map[string]any{
		"userGroupID": strconv.FormatUint(g.ID, 10),
		"handle":      g.Handle,
		"members":     members,
	})
}

func (h *userGroupHandler) memberAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	groupRef, err := toolkit.ReqRef(args, "userGroup")
	if err != nil {
		return nil, err
	}

	userRef, err := toolkit.ReqRef(args, "user")
	if err != nil {
		return nil, err
	}

	g, err := h.resolve(ctx, groupRef, filter.StateInclusive, filter.StateInclusive)
	if err != nil {
		return nil, err
	}

	u, err := sysService.DefaultUser.FindByAny(ctx, userRef)
	if err != nil {
		return nil, toolkit.Errf("user lookup", err)
	}

	if err := sysService.DefaultUserGroup.MemberAdd(ctx, g.ID, u.ID); err != nil {
		return nil, toolkit.Errf("user group member add", err)
	}

	// A user carries a single UserGroupID, so this was a move rather than an
	// addition. Saying so keeps the caller from assuming the previous membership
	// survived.
	return toolkit.TextResult(
		"user %s (%d) is now a member of user group %s (%d), and of no other group",
		u.Handle, u.ID, g.Handle, g.ID,
	), nil
}

// refArg unwraps and resolves the single required group reference the
// single-target ops share.
func (h *userGroupHandler) refArg(ctx context.Context, req mcp.CallToolRequest) (*sysTypes.UserGroup, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "userGroup")
	if err != nil {
		return nil, err
	}

	return h.resolve(ctx, ref, filter.StateInclusive, filter.StateInclusive)
}

// resolve turns a user group reference — a numeric ID or a handle — into the
// group itself.
//
// It goes through Search rather than the service's FindByID / FindByHandle /
// FindByAny deliberately. Those three run the store lookup, load labels and stop
// there: unlike role.onLookup (which calls CanReadRole) or compose
// namespace.lookup (CanReadNamespace), they never consult the access controller,
// so a tool built on them would hand any authenticated caller any group by ID or
// handle. Search is the authorized read path — it gates on CanSearchUserGroups
// and applies CanReadUserGroup per row through filter.Check — and the store
// filter supports both an ID set and an exact handle match, so it resolves
// exactly what FindByAny would. See CONVENTIONS.md §8.6; the missing check in
// FindByID is a service gap, and widening the tool surface is not this file's
// job to fix.
func (h *userGroupHandler) resolve(ctx context.Context, ref string, deleted, archived filter.State) (*sysTypes.UserGroup, error) {
	if ref == "" {
		return nil, fmt.Errorf("userGroup is required")
	}

	f := sysTypes.UserGroupFilter{Deleted: deleted, Archived: archived}
	if _, err := strconv.ParseUint(ref, 10, 64); err == nil {
		f.UserGroupID = []string{ref}
	} else {
		f.Handle = ref
	}

	set, _, err := sysService.DefaultUserGroup.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("user group lookup", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("user group %q not found", ref)
	}

	return set[0], nil
}

// parsePaths turns the `parents` argument into the group's config paths.
//
// UserGroupConfig.Paths[].SelfID is the *parent* group's ID, not the group's
// own — the permission graph is built as
// AddNode(child, handle, GroupNodePath{SelfID: parent}). The param is named
// `parents` because "selfID pointing at someone else" is a trap worth not
// passing on to the caller.
//
// Each element is either a bare reference string or an object
// {"parent": "<handle or ID>", "name": "<label>"}. The bare form covers the
// common single-parent case; the object form exists because the service rejects
// a group whose parent links do not have distinct names, so more than one parent
// forces the caller to name them.
func (h *userGroupHandler) parsePaths(ctx context.Context, raw string) ([]sysTypes.UserGroupPath, error) {
	if raw == "" {
		return nil, nil
	}

	var elems []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &elems); err != nil {
		return nil, fmt.Errorf(
			"invalid parents: expected a JSON array such as [\"engineering\"] or "+
				"[{\"parent\":\"engineering\",\"name\":\"reporting-line\"}]: %w", err,
		)
	}

	out := make([]sysTypes.UserGroupPath, 0, len(elems))
	for _, e := range elems {
		var (
			ref  string
			name string
		)

		if err := json.Unmarshal(e, &ref); err != nil {
			var obj struct {
				Parent string `json:"parent"`
				Name   string `json:"name"`
			}
			if err := json.Unmarshal(e, &obj); err != nil {
				return nil, fmt.Errorf(
					"invalid parents entry %s: expected a group reference string or "+
						"{\"parent\":\"<handle or ID>\",\"name\":\"<label>\"}", string(e),
				)
			}
			ref, name = obj.Parent, obj.Name
		}

		if ref == "" {
			return nil, fmt.Errorf("invalid parents entry %s: the parent reference is empty", string(e))
		}

		// Config.Paths carries IDs only, so a parent named by handle has to be
		// resolved here. A deleted group is not a usable parent, so it does not
		// resolve.
		parent, err := h.resolve(ctx, ref, filter.StateExcluded, filter.StateInclusive)
		if err != nil {
			return nil, err
		}

		out = append(out, sysTypes.UserGroupPath{SelfID: parent.ID, Name: name})
	}

	return out, nil
}

// newUserGroupItem builds the list projection for one group.
func newUserGroupItem(g *sysTypes.UserGroup) userGroupItem {
	item := userGroupItem{
		UserGroupID: strconv.FormatUint(g.ID, 10),
		Handle:      g.Handle,
		IsRoot:      g.IsRoot,
		ParentIDs:   []string{},
		ArchivedAt:  g.ArchivedAt,
		DeletedAt:   g.DeletedAt,
	}

	if g.Meta != nil {
		item.Short = g.Meta.Short
	}

	if g.Config != nil {
		for _, p := range g.Config.Paths {
			item.ParentIDs = append(item.ParentIDs, strconv.FormatUint(p.SelfID, 10))
		}
	}

	return item
}

// userGroupState maps an include-X boolean onto the tri-state filter the store
// uses. Excluded is the default for both deleted and archived, which is why a
// caller has to opt in to see either.
func userGroupState(include bool) filter.State {
	if include {
		return filter.StateInclusive
	}
	return filter.StateExcluded
}
