package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/crusttech/human/server/pkg/weburl"
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
		return toolkit.JSONResultWith(app, applicationLinks(app))
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

	// A custom application is served by the app view at its own ID, so the URL
	// cannot be known until the record exists. Nothing below this layer fills it
	// in — the service and REST both take the unify block as given.
	if customApplicationNeedsURL(app) {
		app.Unify.Url = customApplicationPath(app.ID)
		if app, err = sysService.DefaultApplication.Update(ctx, app); err != nil {
			return nil, toolkit.Errf("custom application url", err)
		}
	}

	return toolkit.JSONResultWith(app, applicationLinks(app))
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

	if customApplicationNeedsURL(app) {
		app.Unify.Url = customApplicationPath(app.ID)
	}

	app, err = sysService.DefaultApplication.Update(ctx, app)
	if err != nil {
		return nil, toolkit.Errf("application update", err)
	}
	return toolkit.JSONResultWith(app, applicationLinks(app))
}

// applicationSourceResult reports a stored document back: enough to verify what
// arrived, plus the link that opens it.
type applicationSourceResult struct {
	ApplicationID string `json:"applicationID"`
	Size          int    `json:"size"`
	Hash          string `json:"hash"`
	URL           string `json:"url"`
}

// applicationSourcePayload is the document and the meta the app view reads
// instead of it. The same shape REST serves on /application/{id}/source.
type applicationSourcePayload struct {
	ApplicationID string                          `json:"applicationID"`
	Source        string                          `json:"source"`
	SourceMeta    *sysTypes.ApplicationSourceMeta `json:"sourceMeta,omitempty"`
}

func (h *applicationHandler) sourceGet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// 'access' rather than the read that resolving already required: the source
	// is the application as a user experiences it, so whoever may open it may
	// read it. Same check REST makes.
	if !sysService.DefaultAccessControl.CanAccessApplication(ctx, app) {
		return nil, toolkit.Errf("application source read", sysService.ApplicationErrNotAllowedToRead())
	}

	return toolkit.JSONResult(applicationSourcePayload{
		ApplicationID: strconv.FormatUint(app.ID, 10),
		Source:        app.Source,
		SourceMeta:    app.SourceMeta,
	})
}

func (h *applicationHandler) sourceSet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	source, whole := toolkit.OptStr(args, "source")
	oldString, patching := toolkit.OptStr(args, "old_string")
	newString, replacement := toolkit.OptStr(args, "new_string")

	switch {
	case whole && (patching || replacement):
		return nil, fmt.Errorf("send either 'source' or an 'old_string'/'new_string' pair, not both")
	case patching != replacement:
		return nil, fmt.Errorf("a patch needs both 'old_string' and 'new_string'; pass an empty 'new_string' to delete the matched text")
	case patching:
		if source, err = applyApplicationSourcePatch(app.Source, oldString, newString); err != nil {
			return nil, err
		}
	case !whole:
		return nil, fmt.Errorf("nothing to store: pass 'source' with the whole document, or 'old_string' and 'new_string' to patch the stored one")
	}

	if err = checkApplicationSource(source); err != nil {
		return nil, err
	}

	// The declaration is what the bridge enforces at runtime, so a patch that
	// says nothing about it carries the stored one forward rather than emptying
	// the allowlist under an app that still works.
	meta := &sysTypes.ApplicationSourceMeta{}
	if app.SourceMeta != nil {
		meta.Namespace = app.SourceMeta.Namespace
		meta.Modules = app.SourceMeta.Modules
	}
	if v, ok := toolkit.OptStr(args, "namespace"); ok {
		meta.Namespace = v
	}
	if raw, ok := args["modules"]; ok && raw != nil {
		if meta.Modules, err = applicationStringList(raw); err != nil {
			return nil, fmt.Errorf(`invalid modules: must be a JSON array of module handles, e.g. ["Lead","Deal"]: %w`, err)
		}
	}

	// SetSource reloads the application itself and writes the hash and size onto
	// the meta it was handed, which is where the result below reads them from.
	if err = sysService.DefaultApplication.SetSource(ctx, app, source, meta); err != nil {
		return nil, toolkit.Errf("application source set", err)
	}

	return toolkit.JSONResult(applicationSourceResult{
		ApplicationID: strconv.FormatUint(app.ID, 10),
		Size:          meta.Size,
		Hash:          meta.Hash,
		URL:           weburl.CustomApplication(app.ID),
	})
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

// customApplicationNeedsURL reports a custom application whose unify block has
// no URL of its own. An explicit URL is left alone — the kind decides how the
// application is rendered, not where the selector points.
func customApplicationNeedsURL(app *sysTypes.Application) bool {
	return app != nil &&
		app.ID != 0 &&
		app.Unify != nil &&
		app.Unify.Kind == sysService.ApplicationKindCustom &&
		app.Unify.Url == ""
}

// customApplicationPath is the unify URL of a custom application: the app
// view's route, relative, the way the selector stores every local URL.
func customApplicationPath(id uint64) string {
	return "app/" + strconv.FormatUint(id, 10)
}

// applyApplicationSourcePatch replaces the one occurrence of oldString in the
// stored document.
//
// Anything other than exactly one match is refused with the count. A patch that
// silently hit the second of three copies edits a place the caller never read,
// and the document is one file with no history to recover from.
func applyApplicationSourcePatch(stored, oldString, newString string) (string, error) {
	if oldString == "" {
		return "", fmt.Errorf("old_string is empty; pass the exact text to replace, or send the whole document as 'source'")
	}

	if stored == "" {
		return "", fmt.Errorf("this application has no source stored yet; send the whole document as 'source'")
	}

	switch n := strings.Count(stored, oldString); n {
	case 1:
		return strings.Replace(stored, oldString, newString, 1), nil
	case 0:
		return "", fmt.Errorf("old_string does not appear in the stored source; read it with system_application_source_get and copy the text exactly, whitespace included")
	default:
		return "", fmt.Errorf("old_string matches %d times in the stored source; extend it with the surrounding lines until exactly one place matches", n)
	}
}

// checkApplicationSource refuses a document the app sandbox cannot run.
//
// The sandbox serves the HTML with a CSP of default-src 'none', connect-src
// 'none' and script-src limited to inline script and cdnjs: an ES module never
// loads, and every network call fails with nothing to catch it on. Both faults
// surface as a blank frame in front of a user, minutes after this call
// returned success, so they are refused here where the wording can say what to
// write instead.
func checkApplicationSource(source string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("the source is empty; a custom application is one whole HTML document")
	}

	if len(source) > sysService.ApplicationSourceMaxSize {
		return fmt.Errorf(
			"the source is %d bytes and the limit is %d; keep a custom application to a single document with inline script, well under that",
			len(source), sysService.ApplicationSourceMaxSize,
		)
	}

	const plainHTML = "; write one plain HTML document with inline <script>, loading libraries from https://cdnjs.cloudflare.com"

	if strings.Contains(source, `<script type="module"`) || strings.Contains(source, "<script type='module'") {
		return fmt.Errorf(`the source has a <script type="module">, which the sandbox's script-src never loads` + plainHTML)
	}

	if strings.Contains(source, "from 'react'") || strings.Contains(source, `from "react"`) {
		return fmt.Errorf("the source imports React, which the sandbox cannot load and which needs a build step" + plainHTML)
	}

	for i, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		for _, keyword := range []string{"import ", "export "} {
			if strings.HasPrefix(trimmed, keyword) {
				return fmt.Errorf("line %d is an ES module statement (%q), which the sandbox never evaluates%s", i+1, strings.TrimSpace(keyword), plainHTML)
			}
		}
	}

	for _, api := range []string{"fetch(", "XMLHttpRequest", "WebSocket"} {
		if usesIdentifier(source, api) {
			return fmt.Errorf(
				"the source uses %s, and the sandbox is served with connect-src 'none' — the call is blocked at runtime and returns nothing to render; read data through the bridge (human.records.list, human.records.read, human.records.report) instead",
				strings.TrimSuffix(api, "("),
			)
		}
	}

	for _, store := range []string{"localStorage", "sessionStorage", "indexedDB", "document.cookie"} {
		if usesIdentifier(source, store) {
			return fmt.Errorf(
				"the source uses %s; the sandbox's origin is opaque, so every storage API throws SecurityError there — keep state in page variables, it lasts as long as the page is open",
				store,
			)
		}
	}

	for _, dialog := range []string{"alert", "confirm", "prompt"} {
		if callsGlobal(source, dialog) {
			return fmt.Errorf(
				"the source calls %s(), which the sandbox suppresses without a word (confirm() always answers false, so what it guards never runs); show an in-page dialog instead",
				dialog,
			)
		}
	}

	if strings.Contains(source, "window.open(") {
		return fmt.Errorf("the source calls window.open(), and the sandbox allows no new windows; show the content in the page, for example in a detail pane")
	}

	if newTab.MatchString(source) {
		return fmt.Errorf(`the source has target="_blank", and the sandbox allows no new windows, so the link does nothing; show the content in the page instead`)
	}

	if downloadLink.MatchString(source) || downloadProperty.MatchString(source) {
		return fmt.Errorf("the source offers a file through a download link, which the sandbox blocks; downloads are not available to a custom app yet, so leave the export out")
	}

	for _, m := range linkHref.FindAllStringSubmatch(source, -1) {
		href := strings.TrimSpace(m[1] + m[2] + m[3])
		if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(strings.ToLower(href), "javascript:") {
			continue
		}
		return fmt.Errorf(
			"the source links to %s; in Human a link that leaves the page does nothing, though it works in a preview — the app is one document, so link only to #anchors within it, and show an email address or phone number as text",
			href,
		)
	}

	for _, m := range externalScript.FindAllStringSubmatch(source, -1) {
		if !strings.HasPrefix(strings.ToLower(m[1]), "https://cdnjs.cloudflare.com/") {
			return fmt.Errorf("the source loads a script from %s; the sandbox's script-src admits https://cdnjs.cloudflare.com only — load the library from there, pinned to an exact version", m[1])
		}
	}

	if m := externalLink.FindStringSubmatch(source); m != nil {
		return fmt.Errorf("the source links %s; the sandbox loads no external stylesheet or font — inline the CSS in <style> and use a system font stack", m[1])
	}

	if m := externalCSSURL.FindStringSubmatch(source); m != nil {
		return fmt.Errorf("the source's CSS reaches %s; the sandbox loads nothing from outside — inline it, or use a data: URL", m[1])
	}

	if m := externalImage.FindStringSubmatch(source); m != nil {
		return fmt.Errorf("the source shows the image %s; the sandbox's img-src admits data: URLs only — embed it as data: or drop it", m[1])
	}

	return nil
}

var (
	newTab           = regexp.MustCompile(`(?i)target\s*=\s*["']?_blank`)
	linkHref         = regexp.MustCompile(`(?i)<a\b[^>]*\shref\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]*))`)
	downloadLink     = regexp.MustCompile(`(?i)<a\b[^>]*\sdownload\b`)
	downloadProperty = regexp.MustCompile(`\.download\s*=[^=]`)
	externalScript   = regexp.MustCompile(`(?i)<script\b[^>]*\ssrc\s*=\s*["']?((?:https?:)?//[^"'\s>]+)`)
	externalLink     = regexp.MustCompile(`(?i)<link\b[^>]*\shref\s*=\s*["']?((?:https?:)?//[^"'\s>]+)`)
	externalCSSURL   = regexp.MustCompile(`(?i)(?:@import\s+(?:url\()?|url\()\s*["']?((?:https?:)?//[^"')\s]+)`)
	externalImage    = regexp.MustCompile(`(?i)<img\b[^>]*\ssrc\s*=\s*["']?((?:https?:)?//[^"'\s>]+)`)
)

// callsGlobal reports whether the source calls the global name — `name(` or
// `window.name(` — rather than a method of that name on some object of its own.
func callsGlobal(source, name string) bool {
	call := name + "("
	for i := 0; i < len(source); {
		at := strings.Index(source[i:], call)
		if at < 0 {
			return false
		}
		at += i
		switch {
		case at == 0:
			return true
		case strings.HasSuffix(source[:at], "window."):
			return true
		case source[at-1] != '.' && !isIdentifierByte(source[at-1]):
			return true
		}
		i = at + len(call)
	}
	return false
}

// usesIdentifier reports whether the source uses name as an identifier of its
// own rather than as the tail of a longer one — "prefetch(" is not "fetch(".
func usesIdentifier(source, name string) bool {
	for i := 0; i < len(source); {
		at := strings.Index(source[i:], name)
		if at < 0 {
			return false
		}
		at += i
		if at == 0 || !isIdentifierByte(source[at-1]) {
			return true
		}
		i = at + len(name)
	}
	return false
}

func isIdentifierByte(b byte) bool {
	return b == '_' || b == '$' ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
