package agentic

import (
	"context"
	"fmt"
	"sort"
	"strings"

	autService "github.com/crusttech/human/server/automation/service"
	autTypes "github.com/crusttech/human/server/automation/types"
	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/crusttech/human/server/pkg/rbac"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in permission_tools.go.
type permissionHandler struct {
	reg toolRegistrar
}

func PermissionHandler(reg toolRegistrar) *permissionHandler {
	h := &permissionHandler{reg: reg}
	h.register()
	return h
}

type (
	// permissionAccessControl is the slice of a component's generated access
	// controller these tools use. Each of the three implements it identically —
	// the file is generated per component from the same template.
	permissionAccessControl interface {
		CanGrant(ctx context.Context) bool
		List() []map[string]string
		FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error)
		Grant(ctx context.Context, rr ...*rbac.Rule) error
	}

	// component pairs a component's RBAC prefix with the controller that owns
	// its rules. Rules are stored in one global table but written through the
	// component that declares the resource, because only that component's
	// generated validator knows the resource shapes.
	component struct {
		name   string
		prefix string
		ac     permissionAccessControl
	}

	permissionRuleArg struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Access    string `json:"access"`
	}

	permissionRuleItem struct {
		RoleID    string `json:"roleID"`
		Role      string `json:"role"`
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Access    string `json:"access"`
	}

	permissionTypeItem struct {
		Component    string   `json:"component"`
		ResourceType string   `json:"resourceType"`
		Resource     string   `json:"resource"`
		Operations   []string `json:"operations"`
	}
)

// components returns the three components whose permissions these tools reach,
// resolved at call time because compose and automation initialise after the MCP
// registry is wired.
//
// Federation is deliberately absent: it is off unless FEDERATION_ENABLED is
// set, its services are initialised conditionally and later than this, and its
// REST routes do not load at all in a default instance. A federation resource is
// refused by name rather than silently doing nothing.
func components() []component {
	out := make([]component, 0, 3)

	if ac := sysService.DefaultAccessControl; ac != nil {
		out = append(out, component{"system", rbac.ResourceComponent(sysTypes.ComponentRbacResource()), ac})
	}
	if ac := cmpService.DefaultAccessControl; ac != nil {
		out = append(out, component{"compose", rbac.ResourceComponent(cmpTypes.ComponentRbacResource()), ac})
	}
	if ac := autService.DefaultAccessControl; ac != nil {
		out = append(out, component{"automation", rbac.ResourceComponent(autTypes.ComponentRbacResource()), ac})
	}

	return out
}

// componentFor finds the component that owns a resource string.
func componentFor(resource string) (component, bool) {
	return componentIn(components(), resource)
}

func componentIn(comps []component, resource string) (component, bool) {
	prefix := rbac.ResourceComponent(resource)
	for _, c := range comps {
		if c.prefix == prefix {
			return c, true
		}
	}
	return component{}, false
}

// schema returns the grantable catalogue: resource types and their operations.
func (h *permissionHandler) schema(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	var (
		wantComponent = strings.ToLower(strings.TrimSpace(toolkit.Str(args, "component")))
		wantType      = strings.ToLower(strings.TrimSpace(toolkit.Str(args, "resourceType")))
		all           []permissionTypeItem
	)

	for _, c := range components() {
		if wantComponent != "" && wantComponent != c.name && wantComponent != c.prefix {
			continue
		}
		all = append(all, permissionTypes(c)...)
	}

	out := map[string]any{}
	items := all

	if wantType != "" {
		items = make([]permissionTypeItem, 0, len(all))
		for _, it := range all {
			if matchesPermissionType(it.ResourceType, wantType) {
				items = append(items, it)
			}
		}
		// A resource type that matches nothing is the likely failure — the
		// strings are not guessable — so answer with the ones that exist.
		if len(items) == 0 {
			known := make([]string, 0, len(all))
			for _, it := range all {
				known = append(known, it.ResourceType)
			}
			out["knownResourceTypes"] = known
		}
	}

	out["resourceTypes"] = items
	out["count"] = len(items)

	return toolkit.JSONResult(out)
}

// permissionTypes collapses a component's flat operation list — one row per
// resource-type-and-operation pair — into one entry per resource type.
func permissionTypes(c component) []permissionTypeItem {
	var (
		order  []string
		byType = map[string]*permissionTypeItem{}
	)

	for _, row := range c.ac.List() {
		t := row["type"]
		if t == "" {
			continue
		}

		item, ok := byType[t]
		if !ok {
			item = &permissionTypeItem{Component: c.name, ResourceType: t, Resource: row["any"]}
			byType[t] = item
			order = append(order, t)
		}
		item.Operations = append(item.Operations, row["op"])
	}

	out := make([]permissionTypeItem, 0, len(order))
	for _, t := range order {
		out = append(out, *byType[t])
	}
	return out
}

// matchesPermissionType matches the full resource type or its bare tail, so a
// caller who knows only "module" finds "corteza::compose:module".
func matchesPermissionType(resourceType, query string) bool {
	lower := strings.ToLower(resourceType)
	if lower == query {
		return true
	}
	if i := strings.LastIndex(lower, ":"); i >= 0 {
		return lower[i+1:] == query
	}
	return false
}

// lookup reads rules, and — when a resource is named — what the caller may do
// with it.
func (h *permissionHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	roles, err := permissionRoles(ctx, strings.TrimSpace(toolkit.Str(args, "role")))
	if err != nil {
		return nil, err
	}

	var (
		resource = strings.TrimSpace(toolkit.Str(args, "resource"))
		limit    = int(toolkit.Page(args).Limit)

		items   []permissionRuleItem
		skipped []string
		read    int
	)

	// Each component gates rule reading behind its own grant permission, and
	// FindRulesByRoleID answers with the role's rules across every component —
	// so the result is filtered back to the component that authorised it.
	// Anything else would let grant on one component read another's rules.
	for _, c := range components() {
		if !c.ac.CanGrant(ctx) {
			skipped = append(skipped, c.name)
			continue
		}
		read++

		for _, role := range roles {
			rules, err := c.ac.FindRulesByRoleID(ctx, role.ID)
			if err != nil {
				return nil, toolkit.Errf("permission rule lookup", err)
			}

			for _, rule := range rules {
				if rbac.ResourceComponent(rule.Resource) != c.prefix {
					continue
				}
				if resource != "" && !strings.HasPrefix(rule.Resource, resource) {
					continue
				}

				items = append(items, permissionRuleItem{
					RoleID:    fmt.Sprint(rule.RoleID),
					Role:      role.Name,
					Resource:  rule.Resource,
					Operation: rule.Operation,
					Access:    rule.Access.String(),
				})
			}
		}
	}

	if read == 0 {
		return nil, fmt.Errorf("permission rule lookup failed: reading rules needs the grant permission, and you hold it on none of %s", strings.Join(componentNames(), ", "))
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Resource != items[j].Resource {
			return items[i].Resource < items[j].Resource
		}
		if items[i].RoleID != items[j].RoleID {
			return items[i].RoleID < items[j].RoleID
		}
		return items[i].Operation < items[j].Operation
	})

	out := map[string]any{"count": len(items)}
	if len(items) > limit {
		out["truncated"] = true
		out["total"] = len(items)
		out["count"] = limit
		items = items[:limit]
	}
	out["rules"] = items

	if len(skipped) > 0 {
		out["componentsSkipped"] = skipped
	}

	if resource != "" {
		out["yourAccess"] = callerOperations(ctx, resource)
	}

	return toolkit.JSONResult(out)
}

// permissionRoles resolves the role argument, or lists every role the caller can
// see when it is absent.
func permissionRoles(ctx context.Context, ref string) ([]*sysTypes.Role, error) {
	if ref != "" {
		r, err := resolveRole(ctx, ref)
		if err != nil {
			return nil, toolkit.Errf("role lookup", err)
		}
		return []*sysTypes.Role{r}, nil
	}

	rr, _, err := sysService.DefaultRole.Find(ctx, sysTypes.RoleFilter{Paging: filter.Paging{Limit: 200}})
	if err != nil {
		return nil, toolkit.Errf("role lookup", err)
	}
	return rr, nil
}

func componentNames() []string { return namesOf(components()) }

func namesOf(comps []component) []string {
	out := make([]string, 0, len(comps))
	for _, c := range comps {
		out = append(out, c.name)
	}
	return out
}

// callerOperations is the ceiling on what the caller may grant on one resource:
// the operations they themselves hold on it.
func callerOperations(ctx context.Context, resource string) []string {
	c, ok := componentFor(resource)
	if !ok {
		return nil
	}

	resourceType := rbac.ResourceType(resource)
	out := make([]string, 0, 8)

	for _, it := range permissionTypes(c) {
		if it.ResourceType != resourceType {
			continue
		}
		for _, op := range it.Operations {
			if callerHolds(ctx, resource, op) {
				out = append(out, op)
			}
		}
	}

	return out
}

// callerHolds answers the one question every write here turns on: does the
// calling user already have this operation on this resource?
//
// It runs the same RBAC evaluation the server runs when the caller performs the
// operation themselves, which is what makes the boundary structural rather than
// a rule someone has to remember. Two entry points because one of them refuses
// half the question:
//
//   - Can is the full evaluation, org tree included, and is what a concrete
//     resource deserves. It answers false for ANY resource carrying a wildcard —
//     service.checkValidity rejects those outright — so it cannot be the only
//     path.
//   - Trace evaluates a wildcard resource honestly: rules are matched with
//     path.Match against the requested pattern, so a rule at the same or a
//     broader scope matches and a narrower one does not. That is exactly
//     "do you hold this at this breadth". It skips the user-group branch, so it
//     can only under-report, which is the safe direction for a ceiling.
//
// Both err towards refusal. A caller wrongly refused loses a grant they could
// have made by hand; a caller wrongly allowed escalates.
func callerHolds(ctx context.Context, resource, operation string) bool {
	svc := rbac.Global()
	if svc == nil {
		return false
	}

	var (
		ses = rbac.ContextToSession(ctx)
		res = rbac.NewResource(resource)
	)

	if svc.Can(ses, operation, res) {
		return true
	}

	return svc.Trace(ses, operation, res).Access == rbac.Allow
}

// grant writes allow/deny rules.
func (h *permissionHandler) grant(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.write(ctx, req, false)
}

// revoke sets rules back to inherit, which deletes them.
func (h *permissionHandler) revoke(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.write(ctx, req, true)
}

// write is grant and revoke: they differ only in the access they write and in
// whether the caller may choose it.
//
// Nothing is applied until every rule has passed. A permission change that
// half-lands leaves a role in a state nobody asked for and nothing records, and
// the failures this refuses on are all things the caller can fix and retry.
func (h *permissionHandler) write(ctx context.Context, req mcp.CallToolRequest, revoking bool) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	subject := "permission grant"
	if revoking {
		subject = "permission revoke"
	}

	ref, err := toolkit.ReqRef(args, "role")
	if err != nil {
		return nil, err
	}

	role, err := sysService.DefaultRole.FindByAny(ctx, ref)
	if err != nil {
		return nil, toolkit.Errf("role lookup", err)
	}

	var in []permissionRuleArg
	present, err := toolkit.JSONArg(args, "rules", "rules", &in)
	if err != nil {
		return nil, err
	}
	if !present || len(in) == 0 {
		return nil, fmt.Errorf("%s failed: 'rules' is required and must hold at least one rule", subject)
	}

	byComponent, acs, err := buildRules(role.ID, in, revoking, subject, components(), func(resource, operation string) bool {
		return callerHolds(ctx, resource, operation)
	})
	if err != nil {
		return nil, err
	}

	written := 0
	for name, rules := range byComponent {
		if err = acs[name].Grant(ctx, rules...); err != nil {
			return nil, toolkit.Errf(subject, err)
		}
		written += len(rules)
	}

	verb := "granted"
	if revoking {
		verb = "revoked"
	}

	return toolkit.JSONResult(map[string]any{
		"roleID": fmt.Sprint(role.ID),
		"role":   role.Name,
		verb:     written,
		"note":   "already signed-in users keep the access their session was built with until they sign in again",
	})
}

// buildRules validates every rule and groups them by the component that owns
// them. It is where the grant boundary is enforced, and it is separated from the
// handler so that boundary can be tested against a predicate rather than against
// a live RBAC service.
//
// It returns on the first failure, having written nothing: a permission change
// that half-lands leaves a role in a state nobody asked for, and every failure
// here is one the caller can fix and retry.
func buildRules(roleID uint64, in []permissionRuleArg, revoking bool, subject string, comps []component, holds func(resource, operation string) bool) (map[string][]*rbac.Rule, map[string]permissionAccessControl, error) {
	byComponent := map[string][]*rbac.Rule{}
	acs := map[string]permissionAccessControl{}

	for i, r := range in {
		resource := strings.TrimSpace(r.Resource)
		operation := strings.TrimSpace(r.Operation)

		if resource == "" || operation == "" {
			return nil, nil, fmt.Errorf("%s failed: rule %d needs both 'resource' and 'operation'; system_permission_schema lists them", subject, i)
		}

		access, err := ruleAccess(r.Access, revoking)
		if err != nil {
			return nil, nil, fmt.Errorf("%s failed: rule %d: %w", subject, i, err)
		}

		c, ok := componentIn(comps, resource)
		if !ok {
			return nil, nil, fmt.Errorf("%s failed: rule %d names resource %q, which belongs to no component these tools reach (%s). Federation resources are not covered. system_permission_schema lists every resource string",
				subject, i, resource, strings.Join(namesOf(comps), ", "))
		}

		// The boundary. Checked before the component's own grant permission so
		// that the refusal names the operation rather than the component: a
		// caller who lacks the operation is not going to be helped by being
		// told to ask for grant.
		if !holds(resource, operation) {
			return nil, nil, fmt.Errorf("%s refused: you do not hold %q on %s yourself, and this tool can only pass on permissions you already have. Nothing was changed. system_permission_lookup with that resource returns 'yourAccess' — the operations you can give away",
				subject, operation, resource)
		}

		byComponent[c.name] = append(byComponent[c.name], &rbac.Rule{
			RoleID:    roleID,
			Resource:  resource,
			Operation: operation,
			Access:    access,
		})
		acs[c.name] = c.ac
	}

	return byComponent, acs, nil
}

// ruleAccess reads the access a rule asks for.
//
// Revoke writes Inherit and ignores what was asked. Grant refuses it: "inherit"
// on a grant means "remove this rule", which is a different operation with a
// different risk level, and letting the word through would make a tool
// annotated as a write quietly take access away.
func ruleAccess(raw string, revoking bool) (rbac.Access, error) {
	if revoking {
		return rbac.Inherit, nil
	}

	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "allow":
		return rbac.Allow, nil
	case "deny":
		return rbac.Deny, nil
	case "inherit":
		return rbac.Inherit, fmt.Errorf("access \"inherit\" removes a rule rather than granting one; use system_permission_revoke")
	default:
		return rbac.Inherit, fmt.Errorf("access %q is not one of \"allow\" or \"deny\"", raw)
	}
}
