package agentic

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in auth_client_tools.go, in the same order.
type authClientHandler struct {
	reg toolRegistrar
}

func AuthClientHandler(reg toolRegistrar) *authClientHandler {
	h := &authClientHandler{reg: reg}
	h.register()
	return h
}

// authClientItem is the compact projection returned by list mode.
//
// It carries no secret field, and must never gain one. The service blanks
// AuthClient.Secret on the way out of both LookupByID and Search, so a handler
// that passes the service type through is safe by construction; a hand-built
// projection is the one place that guarantee can be lost, which is why this
// struct enumerates its fields rather than embedding the type.
type authClientItem struct {
	AuthClientID string     `json:"authClientID"`
	Handle       string     `json:"handle"`
	Name         string     `json:"name"`
	Enabled      bool       `json:"enabled"`
	IsDefault    bool       `json:"isDefault"`
	ValidFrom    *time.Time `json:"validFrom"`
	ExpiresAt    *time.Time `json:"expiresAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
}

func (h *authClientHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	includeDeleted := toolkit.Bool(args, "includeDeleted")

	// Reference param is never Required, so an empty ref means "list".
	ref, err := toolkit.Ref(args, "authClient")
	if err != nil {
		return nil, err
	}
	if ref != "" {
		client, err := h.resolve(ctx, ref, includeDeleted)
		if err != nil {
			return nil, err
		}
		// Returning the service type unchanged is safe: LookupByID blanks
		// Secret before it returns, and Search walks its set doing the same.
		// Nothing here reloads the client from the store, which is the only way
		// the credential could come back.
		return toolkit.JSONResult(client)
	}

	f := sysTypes.AuthClientFilter{Handle: toolkit.Str(args, "handle")}
	if includeDeleted {
		f.Deleted = filter.StateInclusive
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultAuthClient.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("auth client list", err)
	}

	items := make([]authClientItem, 0, len(set))
	for _, client := range set {
		items = append(items, authClientProjection(client))
	}

	// Pass the cursor itself, not its String(): String is a human-readable debug
	// rendering ("<id: 123, [FWD]>") while parseCursor expects base64 of the
	// cursor JSON, which is what MarshalJSON emits. json.Marshal renders a nil
	// pointer as null, so the last page is safe.
	return toolkit.JSONResult(map[string]any{
		"authClients":    items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *authClientHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "authClient")
	if err != nil {
		return nil, err
	}

	client, err := h.resolve(ctx, ref, false)
	if err != nil {
		return nil, err
	}

	// Absent leaves the field alone; present-and-empty clears it. Update copies
	// only the fields below off the object it is handed — Secret is not among
	// them — so passing back a client whose secret was blanked on the way out of
	// the service cannot wipe the stored credential.
	if v, ok := args["handle"]; ok {
		client.Handle, _ = v.(string)
	}
	if v, ok := args["validGrant"]; ok {
		client.ValidGrant, _ = v.(string)
	}
	if v, ok := args["redirectURI"]; ok {
		client.RedirectURI, _ = v.(string)
	}
	if v, ok := args["scope"]; ok {
		client.Scope, _ = v.(string)
	}
	if _, ok := args["enabled"]; ok {
		client.Enabled = toolkit.Bool(args, "enabled")
	}
	if _, ok := args["trusted"]; ok {
		client.Trusted = toolkit.Bool(args, "trusted")
	}

	if v, ok := args["name"]; ok {
		if client.Meta == nil {
			client.Meta = &sysTypes.AuthClientMeta{}
		}
		client.Meta.Name, _ = v.(string)
	}
	if v, ok := args["description"]; ok {
		if client.Meta == nil {
			client.Meta = &sysTypes.AuthClientMeta{}
		}
		client.Meta.Description, _ = v.(string)
	}

	if v, ok := args["validFrom"]; ok {
		if s, _ := v.(string); s == "" {
			client.ValidFrom = nil
		} else if client.ValidFrom, err = optTime(args, "validFrom"); err != nil {
			return nil, err
		}
	}
	if v, ok := args["expiresAt"]; ok {
		if s, _ := v.(string); s == "" {
			client.ExpiresAt = nil
		} else if client.ExpiresAt, err = optTime(args, "expiresAt"); err != nil {
			return nil, err
		}
	}

	res, err := sysService.DefaultAuthClient.Update(ctx, client)
	if err != nil {
		return nil, toolkit.Errf("auth client update", err)
	}

	// Update is the one auth client method that hands back an unblanked row: it
	// loads the client straight from the store, copies the changed fields onto
	// it and returns it, without the blanking LookupByID and Search perform.
	// Marshalling that would put a working credential into the caller's context,
	// which is exactly what this family exists to prevent. The store write has
	// already happened, so clearing the field here changes only what is
	// returned.
	if res != nil {
		res.Secret = ""
	}
	return toolkit.JSONResult(res)
}

func (h *authClientHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := h.resolveArg(ctx, req, false)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultAuthClient.DeleteByID(ctx, client.ID); err != nil {
		return nil, toolkit.Errf("auth client deletion", err)
	}
	return toolkit.TextResult("auth client %d deleted", client.ID), nil
}

// undelete resolves with deleted clients included, unlike every other tool
// here: AuthClientFilter defaults Deleted to StateExcluded, so a handle search
// would report "not found" for precisely the clients this tool exists to
// restore.
func (h *authClientHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := h.resolveArg(ctx, req, true)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultAuthClient.UndeleteByID(ctx, client.ID); err != nil {
		return nil, toolkit.Errf("auth client undelete", err)
	}
	return toolkit.TextResult("auth client %d restored", client.ID), nil
}

// resolveArg unwraps the required reference the single-target tools share.
func (h *authClientHandler) resolveArg(ctx context.Context, req mcp.CallToolRequest, includeDeleted bool) (*sysTypes.AuthClient, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "authClient")
	if err != nil {
		return nil, err
	}

	return h.resolve(ctx, ref, includeDeleted)
}

// resolve turns a reference — a numeric ID or a handle — into an auth client.
//
// The service has neither FindByAny nor FindByHandle, only FindByID and
// LookupByID, so a handle is resolved through Search, whose filter carries an
// exact-match Handle field. Both paths go through the service, so both are
// access-checked and both return a client with its secret blanked.
func (h *authClientHandler) resolve(ctx context.Context, ref string, includeDeleted bool) (*sysTypes.AuthClient, error) {
	if ref == "" {
		return nil, fmt.Errorf("auth client identifier required")
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		client, err := sysService.DefaultAuthClient.LookupByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("auth client lookup", err)
		}
		return client, nil
	}

	f := sysTypes.AuthClientFilter{Handle: ref}
	if includeDeleted {
		f.Deleted = filter.StateInclusive
	}

	set, _, err := sysService.DefaultAuthClient.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("auth client lookup", err)
	}

	switch len(set) {
	case 0:
		return nil, fmt.Errorf("auth client %q not found", ref)
	case 1:
		return set[0], nil
	default:
		ids := make([]string, 0, len(set))
		for _, client := range set {
			ids = append(ids, strconv.FormatUint(client.ID, 10))
		}
		return nil, fmt.Errorf(
			"auth client %q is ambiguous: %d clients share that handle (IDs %s) — refer to one by ID",
			ref, len(set), strings.Join(ids, ", "),
		)
	}
}

// authClientProjection trims a client for list mode. It reads only from the
// value the service returned; see authClientItem for why it never grows a
// secret field.
func authClientProjection(client *sysTypes.AuthClient) authClientItem {
	item := authClientItem{
		AuthClientID: strconv.FormatUint(client.ID, 10),
		Handle:       client.Handle,
		Enabled:      client.Enabled,
		IsDefault:    sysService.DefaultAuthClient.IsDefaultClient(client),
		ValidFrom:    client.ValidFrom,
		ExpiresAt:    client.ExpiresAt,
		DeletedAt:    client.DeletedAt,
	}

	if client.Meta != nil {
		item.Name = client.Meta.Name
	}

	return item
}
