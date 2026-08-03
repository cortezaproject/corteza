package agentic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in user_tools.go, in the same order.
type userHandler struct {
	reg toolRegistrar
}

func UserHandler(reg toolRegistrar) *userHandler {
	h := &userHandler{reg: reg}
	h.register()
	return h
}

// Sentinels the user service substitutes for a masked email or name
// (maskPrivateDataEmail / maskPrivateDataName in system/service/user.go). They
// are unexported there, so they are restated here rather than reached for.
//
// They matter because every read path — FindByAny included — masks in place
// when privacy masking is enabled and the caller may not unmask. update() below
// loads the user and writes it back, so without this check a masked value would
// be persisted over the real one.
const (
	maskedUserEmail = "####.#######@######.###"
	maskedUserName  = "##### ##########"
)

// userItem is the compact projection returned by list mode. Single-item lookups
// return the service type unchanged; lists must not, because meta, settings and
// labels per row are noise when the caller is picking one user out of many.
type userItem struct {
	UserID      string     `json:"userID"`
	Email       string     `json:"email"`
	Handle      string     `json:"handle"`
	Name        string     `json:"name"`
	Kind        string     `json:"kind"`
	SuspendedAt *time.Time `json:"suspendedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

func (h *userHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Reference param is never Required, so an empty ref means "list".
	ref, err := toolkit.Ref(args, "user")
	if err != nil {
		return nil, err
	} else if ref != "" {
		u, err := sysService.DefaultUser.FindByAny(ctx, ref)
		if err != nil {
			return nil, toolkit.Errf("user lookup", err)
		}
		return toolkit.JSONResult(u)
	}

	kind := sysTypes.UserKind(toolkit.Str(args, "kind"))
	if kind != sysTypes.NormalUser && kind != sysTypes.SystemUser {
		return nil, fmt.Errorf("invalid kind: expected \"\" for ordinary users or %q for system users, got %q", sysTypes.SystemUser, kind)
	}

	f := sysTypes.UserFilter{
		Query:    toolkit.Str(args, "query"),
		Email:    toolkit.Str(args, "email"),
		Handle:   toolkit.Str(args, "handle"),
		Username: toolkit.Str(args, "username"),
		Kind:     kind,
		AllKinds: toolkit.Bool(args, "allKinds"),
	}

	// Both default to StateExcluded, so deleted and suspended users are hidden
	// unless the caller opts in.
	if toolkit.Bool(args, "includeDeleted") {
		f.Deleted = filter.StateInclusive
	}
	if toolkit.Bool(args, "includeSuspended") {
		f.Suspended = filter.StateInclusive
	}

	if ref, err := toolkit.Ref(args, "role"); err != nil {
		return nil, err
	} else if ref != "" {
		r, err := sysService.DefaultRole.FindByAny(ctx, ref)
		if err != nil {
			return nil, toolkit.Errf("role lookup", err)
		}
		f.RoleID = []string{strconv.FormatUint(r.ID, 10)}
	}

	if ref, err := toolkit.Ref(args, "userGroup"); err != nil {
		return nil, err
	} else if ref != "" {
		g, err := sysService.DefaultUserGroup.FindByAny(ctx, ref)
		if err != nil {
			return nil, toolkit.Errf("user group lookup", err)
		}
		f.UserGroupID = g.ID
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultUser.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("user list", err)
	}

	items := make([]userItem, 0, len(set))
	for _, u := range set {
		items = append(items, userItem{
			UserID:      strconv.FormatUint(u.ID, 10),
			Email:       u.Email,
			Handle:      u.Handle,
			Name:        u.Name,
			Kind:        string(u.Kind),
			SuspendedAt: u.SuspendedAt,
			DeletedAt:   u.DeletedAt,
		})
	}

	// Pass the cursor itself, not its String(): String is a human-readable debug
	// rendering ("<id: 123, [FWD]>") while parseCursor expects base64 of the
	// cursor JSON, which is what MarshalJSON emits. json.Marshal renders a nil
	// pointer as null, so the last page is safe.
	return toolkit.JSONResult(map[string]any{
		"users":          items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *userHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	email, err := toolkit.ReqStr(args, "email")
	if err != nil {
		return nil, err
	}

	u := &sysTypes.User{
		Email:  email,
		Name:   toolkit.Str(args, "name"),
		Handle: toolkit.Str(args, "handle"),
	}

	u, err = sysService.DefaultUser.Create(ctx, u)
	if err != nil {
		return nil, toolkit.Errf("user creation", err)
	}
	return toolkit.JSONResult(u)
}

func (h *userHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	u, err := resolveUser(ctx, args)
	if err != nil {
		return nil, err
	}

	_, emailGiven := args["email"]
	_, nameGiven := args["name"]

	// See the maskedUser* comment: refuse rather than persist a mask over a
	// value the caller was never shown and is not replacing.
	if u.Email == maskedUserEmail && !emailGiven {
		return nil, fmt.Errorf("user update refused: this user's email is masked by the privacy settings for your permissions, and updating would overwrite the real address — pass a replacement email, or ask an administrator with unmask permission")
	}
	if u.Name == maskedUserName && !nameGiven {
		return nil, fmt.Errorf("user update refused: this user's name is masked by the privacy settings for your permissions, and updating would overwrite the real name — pass a replacement name, or ask an administrator with unmask permission")
	}

	// Absent leaves the field alone; present-and-empty clears it. Email is the
	// exception: the service rejects an address that does not parse, so an empty
	// string is refused here with a message that says why rather than surfacing
	// a bare validation error.
	if emailGiven {
		email, _ := args["email"].(string)
		if email == "" {
			return nil, fmt.Errorf("email cannot be cleared: it is the identity the user signs in with, so pass a valid address or omit the field")
		}
		u.Email = email
	}
	if nameGiven {
		u.Name, _ = args["name"].(string)
	}
	if v, ok := args["handle"]; ok {
		u.Handle, _ = v.(string)
	}

	u, err = sysService.DefaultUser.Update(ctx, u)
	if err != nil {
		return nil, toolkit.Errf("user update", err)
	}
	return toolkit.JSONResult(u)
}

func (h *userHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	u, err := resolveUser(ctx, args)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultUser.DeleteByID(ctx, u.ID); err != nil {
		return nil, toolkit.Errf("user deletion", err)
	}
	return toolkit.TextResult("user %d deleted; restore it with system_user_undelete using that ID", u.ID), nil
}

func (h *userHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Deleted users are only reachable by ID: the store's handle, email and
	// username lookups all exclude deleted rows, so FindByAny would fail on
	// anything but a numeric reference.
	id, err := toolkit.ReqID(args, "userID")
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultUser.UndeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("user undeletion", err)
	}
	return toolkit.TextResult("user %d restored", id), nil
}

func (h *userHandler) suspend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	u, err := resolveUser(ctx, args)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultUser.Suspend(ctx, u.ID); err != nil {
		return nil, toolkit.Errf("user suspension", err)
	}
	return toolkit.TextResult("user %d suspended; access tokens revoked", u.ID), nil
}

func (h *userHandler) unsuspend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	u, err := resolveUser(ctx, args)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultUser.Unsuspend(ctx, u.ID); err != nil {
		return nil, toolkit.Errf("user unsuspension", err)
	}
	return toolkit.TextResult("user %d unsuspended", u.ID), nil
}

func (h *userHandler) setEmailConfirmed(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// False is a meaningful value here, so an absent argument must not be read
	// as one — the caller has to say which way round they mean it.
	if _, ok := args["confirmed"]; !ok {
		return nil, fmt.Errorf("confirmed is required: pass true to mark the address confirmed or false to mark it unconfirmed")
	}
	confirmed := toolkit.Bool(args, "confirmed")

	u, err := resolveUser(ctx, args)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultUser.ToggleEmailConfirmation(ctx, u.ID, confirmed); err != nil {
		return nil, toolkit.Errf("user email confirmation change", err)
	}
	if confirmed {
		return toolkit.TextResult("user %d email marked confirmed", u.ID), nil
	}
	return toolkit.TextResult("user %d email marked unconfirmed", u.ID), nil
}

// resolveUser turns the required 'user' reference — ID, handle or email — into
// a user. Every op but undelete shares it; undelete cannot, because the store
// resolves neither handle nor email to a deleted account.
func resolveUser(ctx context.Context, args map[string]any) (*sysTypes.User, error) {
	ref, err := toolkit.ReqRef(args, "user")
	if err != nil {
		return nil, err
	}

	u, err := sysService.DefaultUser.FindByAny(ctx, ref)
	if err != nil {
		return nil, toolkit.Errf("user lookup", err)
	}
	return u, nil
}
