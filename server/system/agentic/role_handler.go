package agentic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in role_tools.go, in the same order.
type roleHandler struct {
	reg toolRegistrar
}

func RoleHandler(reg toolRegistrar) *roleHandler {
	h := &roleHandler{reg: reg}
	h.register()
	return h
}

// roleItem is the compact projection returned by list mode.
//
// isSystem is in it for a reason beyond disambiguation: a system role refuses
// delete, archive and undelete and cannot be renamed, so without it in the list
// a model learns which roles it cannot touch only by trying and failing, once
// per role. isClosed is the same argument for membership changes.
type roleItem struct {
	RoleID     string     `json:"roleID"`
	Name       string     `json:"name"`
	Handle     string     `json:"handle"`
	IsSystem   bool       `json:"isSystem"`
	IsClosed   bool       `json:"isClosed"`
	ArchivedAt *time.Time `json:"archivedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}

// roleDetail is a single role as the service returns it, plus the two
// predicates that decide whether a write against it can succeed. They are
// computed from instance configuration rather than stored on the row, so they
// cannot be read off the role itself — the REST payload adds them for the same
// reason.
type roleDetail struct {
	*sysTypes.Role
	IsSystem bool `json:"isSystem"`
	IsClosed bool `json:"isClosed"`
}

// roleMemberItem projects a role membership. Membership is stored as an RBAC
// resource string identifying either a user or a user group; splitting it into
// a kind and an ID is what a caller can actually act on.
type roleMemberItem struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

func (h *roleHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Reference param is never Required, so an empty ref means "list".
	ref, err := toolkit.Ref(args, "role")
	if err != nil {
		return nil, err
	}
	if ref != "" {
		r, err := resolveRole(ctx, ref)
		if err != nil {
			return nil, toolkit.Errf("role lookup", err)
		}
		return toolkit.JSONResult(newRoleDetail(r))
	}

	f := sysTypes.RoleFilter{
		Query:  toolkit.Str(args, "query"),
		Name:   toolkit.Str(args, "name"),
		Handle: toolkit.Str(args, "handle"),
	}

	if toolkit.Bool(args, "includeDeleted") {
		f.Deleted = filter.StateInclusive
	}
	if toolkit.Bool(args, "includeArchived") {
		f.Archived = filter.StateInclusive
	}

	member, err := toolkit.Ref(args, "member")
	if err != nil {
		return nil, err
	}
	group, err := toolkit.Ref(args, "userGroup")
	if err != nil {
		return nil, err
	}

	// RoleFilter.Resource carries the membership narrowing, and it holds one
	// resource, so the two filters cannot be combined. Saying so beats returning
	// whichever one happened to be applied last.
	switch {
	case member != "" && group != "":
		return nil, fmt.Errorf("member and userGroup cannot be combined: the role filter narrows by one membership at a time")

	case member != "":
		u, err := sysService.DefaultUser.FindByAny(ctx, member)
		if err != nil {
			return nil, toolkit.Errf("user lookup", err)
		}
		// Note this is RoleFilter.Resource and not RoleFilter.MemberID: MemberID
		// and UserGroupID are declared on the filter but no store adapter reads
		// them, so setting them would silently filter nothing. Resource is the
		// field the store honours, and is how the service itself asks for a
		// user's roles (see system/service/user.go).
		f.Resource = sysTypes.UserRbacResource(u.ID)

	case group != "":
		g, err := sysService.DefaultUserGroup.FindByAny(ctx, group)
		if err != nil {
			return nil, toolkit.Errf("user group lookup", err)
		}
		f.Resource = sysTypes.UserGroupRbacResource(g.ID)
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultRole.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("role list", err)
	}

	items := make([]roleItem, 0, len(set))
	for _, r := range set {
		items = append(items, roleItem{
			RoleID:     fmt.Sprintf("%d", r.ID),
			Name:       r.Name,
			Handle:     r.Handle,
			IsSystem:   sysService.DefaultRole.IsSystem(r),
			IsClosed:   sysService.DefaultRole.IsClosed(r),
			ArchivedAt: r.ArchivedAt,
			DeletedAt:  r.DeletedAt,
		})
	}

	// The cursor itself, never its String(): String is a debug rendering that
	// parseCursor cannot read back, so returning it would make paging one-way.
	// A nil pointer marshals to null, so the last page needs no special case.
	return toolkit.JSONResult(map[string]any{
		"roles":          items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *roleHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	// Meta is always allocated: the meta column is NOT NULL, so a role created
	// without a description failed at the store with a constraint violation
	// rather than anything a caller could act on.
	r := &sysTypes.Role{
		Name:   name,
		Handle: toolkit.Str(args, "handle"),
		Meta:   &sysTypes.RoleMeta{Description: toolkit.Str(args, "description")},
	}

	r, err = sysService.DefaultRole.Create(ctx, r)
	if err != nil {
		return nil, toolkit.Errf("role creation", err)
	}
	return toolkit.JSONResult(newRoleDetail(r))
}

func (h *roleHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "role")
	if err != nil {
		return nil, err
	}

	// The service copies name, handle and meta from the update wholesale, so the
	// current role is the starting point — building a bare types.Role here would
	// blank every field the caller did not mention.
	current, err := resolveRole(ctx, ref)
	if err != nil {
		return nil, toolkit.Errf("role lookup", err)
	}

	upd := current.Clone()

	// Absent leaves the field alone; present-and-empty clears it.
	if v, ok := args["name"]; ok {
		upd.Name, _ = v.(string)
	}
	if v, ok := args["handle"]; ok {
		upd.Handle, _ = v.(string)
	}
	if v, ok := args["description"]; ok {
		s, _ := v.(string)
		// Meta is replaced wholesale by the service too, so mutate the clone
		// rather than building a fresh RoleMeta: a role's context expression
		// lives in the same struct and rebuilding would drop it.
		if upd.Meta == nil {
			upd.Meta = &sysTypes.RoleMeta{}
		}
		upd.Meta.Description = s
	}

	// UpdatedAt stays nil: the service reads it as an optimistic-locking stamp
	// and treats nil as "no stale-data check". A tool caller has nowhere to hold
	// a version between calls, so sending one it did not read would be theatre.
	upd.UpdatedAt = nil

	res, err := sysService.DefaultRole.Update(ctx, upd)
	if err != nil {
		return nil, toolkit.Errf("role update", err)
	}
	return toolkit.JSONResult(newRoleDetail(res))
}

func (h *roleHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := roleRefID(ctx, req, "role")
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultRole.DeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("role deletion", err)
	}
	return toolkit.TextResult("role %d deleted; its members no longer hold what it granted", id), nil
}

func (h *roleHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// An ID, not a ref: the store's handle and name lookups exclude deleted
	// roles, so a deleted role is unreachable by any other form of reference.
	id, err := toolkit.ReqID(args, "roleID")
	if err != nil {
		return nil, err
	}

	if err := sysService.DefaultRole.UndeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("role undeletion", err)
	}
	return toolkit.TextResult("role %d restored; its members hold what it grants again", id), nil
}

func (h *roleHandler) archive(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := roleRefID(ctx, req, "role")
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultRole.Archive(ctx, id); err != nil {
		return nil, toolkit.Errf("role archiving", err)
	}
	return toolkit.TextResult("role %d archived; it stops applying until unarchived", id), nil
}

func (h *roleHandler) unarchive(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := roleRefID(ctx, req, "role")
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultRole.Unarchive(ctx, id); err != nil {
		return nil, toolkit.Errf("role unarchiving", err)
	}
	return toolkit.TextResult("role %d unarchived; it applies to its members again", id), nil
}

func (h *roleHandler) memberList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := roleRefID(ctx, req, "role")
	if err != nil {
		return nil, err
	}

	mm, err := sysService.DefaultRole.MemberList(ctx, id)
	if err != nil {
		return nil, toolkit.Errf("role member list", err)
	}

	items := make([]roleMemberItem, 0, len(mm))
	for _, m := range mm {
		kind := "user"
		if strings.HasPrefix(m.Resource, sysTypes.UserGroupResourceType+"/") {
			kind = "userGroup"
		}
		items = append(items, roleMemberItem{
			Kind: kind,
			ID:   fmt.Sprintf("%d", rbac.ResourceID(m.Resource)),
		})
	}

	return toolkit.JSONResult(map[string]any{
		"roleID":  fmt.Sprintf("%d", id),
		"members": items,
	})
}

func (h *roleHandler) memberAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	roleID, userID, groupID, err := roleMemberArgs(ctx, req)
	if err != nil {
		return nil, err
	}

	if userID > 0 {
		if err := sysService.DefaultRole.MemberAdd(ctx, roleID, userID); err != nil {
			return nil, toolkit.Errf("role member add", err)
		}
		return toolkit.TextResult("user %d added to role %d and now holds everything it grants", userID, roleID), nil
	}

	if err := sysService.DefaultRole.MemberAddGroup(ctx, roleID, groupID); err != nil {
		return nil, toolkit.Errf("role member group add", err)
	}
	return toolkit.TextResult("user group %d added to role %d; every member of the group now holds what it grants", groupID, roleID), nil
}

func (h *roleHandler) memberRemove(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	roleID, userID, groupID, err := roleMemberArgs(ctx, req)
	if err != nil {
		return nil, err
	}

	if userID > 0 {
		if err := sysService.DefaultRole.MemberRemove(ctx, roleID, userID); err != nil {
			return nil, toolkit.Errf("role member remove", err)
		}
		return toolkit.TextResult("user %d removed from role %d", userID, roleID), nil
	}

	if err := sysService.DefaultRole.MemberRemoveGroup(ctx, roleID, groupID); err != nil {
		return nil, toolkit.Errf("role member group remove", err)
	}
	return toolkit.TextResult("user group %d removed from role %d", groupID, roleID), nil
}

func (h *roleHandler) cloneRules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	from, err := roleRefID(ctx, req, "role")
	if err != nil {
		return nil, err
	}

	to, err := roleRefID(ctx, req, "targetRole")
	if err != nil {
		return nil, err
	}

	if from == to {
		return nil, fmt.Errorf("role and targetRole resolve to the same role %d: cloning a role's rules onto itself rewrites them for no effect", from)
	}

	if err := sysService.DefaultRole.CloneRules(ctx, from, to); err != nil {
		return nil, toolkit.Errf("role rule cloning", err)
	}
	return toolkit.TextResult(
		"permission rules copied from role %d onto role %d; role %d's previous rules were discarded and its members now hold role %d's access",
		from, to, to, from,
	), nil
}

// newRoleDetail attaches the two configuration-derived predicates to a role.
func newRoleDetail(r *sysTypes.Role) roleDetail {
	return roleDetail{
		Role:     r,
		IsSystem: sysService.DefaultRole.IsSystem(r),
		IsClosed: sysService.DefaultRole.IsClosed(r),
	}
}

// resolveRole turns a role reference — ID, handle or name — into a role the
// caller is allowed to read.
//
// FindByAny only access-checks the ID path: it routes an ID to FindByID, which
// refuses a role the caller cannot read, while a handle or a name goes to
// FindByHandle / FindByName, which read the store and check nothing. Re-fetching
// by ID puts every reference form through the same check, so a caller cannot
// read by handle what it could not read by ID.
func resolveRole(ctx context.Context, ref string) (*sysTypes.Role, error) {
	r, err := sysService.DefaultRole.FindByAny(ctx, ref)
	if err != nil {
		return nil, err
	}
	return sysService.DefaultRole.FindByID(ctx, r.ID)
}

// roleRefID resolves a required role reference argument to an ID.
//
// Unlike resolveRole this does not re-fetch through the read check: the write
// and membership operations it feeds each apply their own guard — update and
// archive check update permission, delete checks delete permission, the member
// operations check member management — and requiring read permission on top
// would invent a restriction the service does not have.
func roleRefID(ctx context.Context, req mcp.CallToolRequest, key string) (uint64, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return 0, err
	}

	ref, err := toolkit.ReqRef(args, key)
	if err != nil {
		return 0, err
	}

	r, err := sysService.DefaultRole.FindByAny(ctx, ref)
	if err != nil {
		return 0, toolkit.Errf("role lookup", err)
	}
	return r.ID, nil
}

// roleMemberArgs unwraps the arguments the two membership operations share: a
// role, and exactly one of a user or a user group. Exactly one, because the
// service has a separate call for each and silently preferring one over the
// other would make the tool do something the caller did not ask for.
func roleMemberArgs(ctx context.Context, req mcp.CallToolRequest) (roleID, userID, groupID uint64, err error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return 0, 0, 0, err
	}

	if roleID, err = roleRefID(ctx, req, "role"); err != nil {
		return 0, 0, 0, err
	}

	userRef, err := toolkit.Ref(args, "user")
	if err != nil {
		return 0, 0, 0, err
	}
	groupRef, err := toolkit.Ref(args, "userGroup")
	if err != nil {
		return 0, 0, 0, err
	}

	switch {
	case userRef != "" && groupRef != "":
		return 0, 0, 0, fmt.Errorf("pass either user or userGroup, not both: adding or removing a user and a group are separate changes")

	case userRef == "" && groupRef == "":
		return 0, 0, 0, fmt.Errorf("one of user or userGroup is required")

	case userRef != "":
		u, err := sysService.DefaultUser.FindByAny(ctx, userRef)
		if err != nil {
			return 0, 0, 0, toolkit.Errf("user lookup", err)
		}
		userID = u.ID

	default:
		g, err := sysService.DefaultUserGroup.FindByAny(ctx, groupRef)
		if err != nil {
			return 0, 0, 0, toolkit.Errf("user group lookup", err)
		}
		groupID = g.ID
	}

	return roleID, userID, groupID, nil
}
