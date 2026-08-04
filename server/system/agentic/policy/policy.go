package policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/crusttech/human/server/system/types"
)

// toolAliases mirrors the alias map in the mcp package.
// Keep both in sync when renaming tools.
//
// The duplication is deliberate — policy must not import mcp, which would close
// the cycle mcp → runtime → policy. ToolAliases below exposes it so a test can
// assert the two agree.
var toolAliases = map[string]string{
	"compose_namespace_list": "compose_namespace_lookup",
	"compose_module_list":    "compose_module_lookup",
}

// ToolAliases returns a copy of this package's alias map, for the structural
// test that keeps it in step with the mcp package's.
func ToolAliases() map[string]string {
	out := make(map[string]string, len(toolAliases))
	for k, v := range toolAliases {
		out[k] = v
	}
	return out
}

type (
	OwnershipFallback func(ctx context.Context, args ValueGetter) bool

	ValueGetter interface {
		Get(key string) (any, bool)
		Keys() []string
	}

	ValueSetter interface {
		Set(key string, val any)
	}

	MapValues map[string]any

	Decision struct {
		Allowed       bool
		Reason        string
		SanitizedArgs map[string]any
	}
)

func (m MapValues) Get(key string) (any, bool) { v, ok := m[key]; return v, ok }
func (m MapValues) Set(key string, val any)    { m[key] = val }
func (m MapValues) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Evaluate(ctx context.Context, agent *types.Agent, tool string, args ValueGetter, ownsTarget OwnershipFallback) Decision {
	if tool == "automation_workflow_exec" {
		ref, _ := args.Get("workflow")
		if findWorkflow(agent, fmt.Sprintf("%v", ref)) == nil {
			return Decision{Allowed: false, Reason: fmt.Sprintf("agent is not allowed to execute workflow %q", ref)}
		}
		return allowedDecision(agent, nil, args)
	}
	if tool == "automation_taq_lookup" && len(agent.Access.TAQs) > 0 {
		return allowedDecision(agent, nil, args)
	}
	if tool == "automation_workflow_lookup" && len(agent.Access.Workflows) > 0 {
		return allowedDecision(agent, nil, args)
	}

	if tool == "compose_page_block_schema" && agentHasPageTool(agent) {
		return allowedDecision(agent, nil, args)
	}

	// Per-TAQ tools are minted at runtime as "automation_<taqID>"
	// (runtime/executor.go mints them; the ID is always numeric). Matching on a
	// numeric suffix is what distinguishes them from the static automation_*
	// families.
	//
	// This used to be a negative match — any automation_* name outside a
	// hardcoded three-item exception list was read as a TAQ ID — which denied
	// automation_taq_exec, automation_taq_executions and
	// automation_taq_execution_trace outright, with the nonsense reason
	// `agent is not allowed to execute automation "taq_exec"`. The executor
	// already discriminates correctly by prefix; this brings policy in line.
	if taqIDStr, ok := dynamicTAQRef(tool); ok {
		if findTAQ(agent, taqIDStr) == nil {
			return Decision{Allowed: false, Reason: fmt.Sprintf("agent is not allowed to execute automation %q", taqIDStr)}
		}
		return allowedDecision(agent, nil, args)
	}

	var entry *types.AgentAccessTool
	for i := range agent.Access.Tools {
		stored := agent.Access.Tools[i].Name
		if alias, ok := toolAliases[stored]; ok {
			stored = alias
		}
		if stored == tool {
			entry = &agent.Access.Tools[i]
			break
		}
	}

	if entry == nil {
		return Decision{
			Allowed: false,
			Reason:  fmt.Sprintf("tool %q is not in the agent's allow-list", tool),
		}
	}

	// Namespace create has no source namespace to validate against — the namespace
	// doesn't exist yet. RBAC handles the real permission check inside the service.
	if tool == "compose_namespace_create" {
		return allowedDecision(agent, entry, args)
	}

	// An unmapped tool gets no resource-level narrowing. That is deliberately
	// NOT denied here: Evaluate has already required the tool to be in the
	// agent's Access.Tools allow-list above, and denying at runtime would mean
	// every newly added tool silently breaks in-process agents until someone
	// remembers to edit this file — a cross-package coupling a tool author has
	// no reason to discover.
	//
	// The gap is caught at CI time instead: IsClassified below backs a test that
	// asserts every registered tool is either mapped by buildResource or listed
	// in resourceScopeExempt. The problem was that the default was *silent*, not
	// that it was permissive.
	resource, _ := buildResource(tool, args)

	if d := checkAllow(entry.Allow, resource); !d.Allowed {
		if ownsTarget != nil && ownsTarget(ctx, args) {
			return allowedDecision(agent, entry, args)
		}
		return d
	}

	return allowedDecision(agent, entry, args)
}

func allowedDecision(agent *types.Agent, entry *types.AgentAccessTool, args ValueGetter) Decision {
	keys := args.Keys()
	sanitized := make(map[string]any, len(keys))
	for _, k := range keys {
		v, _ := args.Get(k)
		sanitized[k] = v
	}

	for k, v := range agent.Access.Context.Defaults {
		if _, exists := sanitized[k]; !exists {
			sanitized[k] = v
		}
	}

	if entry != nil {
		for k, v := range entry.Context.Defaults {
			if _, exists := sanitized[k]; !exists {
				sanitized[k] = v
			}
		}
		for k, v := range entry.Context.Overrides {
			sanitized[k] = v
		}
	}

	return Decision{
		Allowed:       true,
		Reason:        "tool is in agent's allow-list",
		SanitizedArgs: sanitized,
	}
}

func findTAQ(agent *types.Agent, ref string) *types.AgentAccessTAQ {
	for i := range agent.Access.TAQs {
		t := &agent.Access.TAQs[i]
		if fmt.Sprintf("%d", t.ID) == ref {
			return t
		}
	}
	return nil
}

func findWorkflow(agent *types.Agent, ref string) *types.AgentAccessWorkflow {
	for i := range agent.Access.Workflows {
		w := &agent.Access.Workflows[i]
		if fmt.Sprintf("%d", w.ID) == ref {
			return w
		}
	}
	return nil
}

func agentHasPageTool(agent *types.Agent) bool {
	for _, t := range agent.Access.Tools {
		if strings.HasPrefix(t.Name, "compose_page_") {
			return true
		}
	}
	return false
}

// dynamicTAQRef reports whether tool is a runtime-minted per-TAQ tool
// ("automation_<numeric id>") and returns the id portion.
func dynamicTAQRef(tool string) (string, bool) {
	rest, ok := strings.CutPrefix(tool, "automation_")
	if !ok || rest == "" {
		return "", false
	}
	for _, r := range rest {
		if r < '0' || r > '9' {
			return "", false
		}
	}
	return rest, true
}

// resourceScopeExempt lists tools that deliberately carry no compose/automation
// resource dimension, so checkAllow has nothing to narrow on. Being on this list
// is not a free pass: Evaluate has already required the tool to be in the
// agent's Access.Tools allow-list before reaching here.
//
// Membership is asserted at CI time via IsClassified, so a new tool arriving
// with no resource scoping fails a test rather than shipping unnoticed.
//
// TODO: compose_page_* and compose_chart_* are namespace-scoped in reality and
// should graduate to real buildResource cases. Left exempt here because adding
// mappings would change authorization for existing agent configs, which needs
// its own migration rather than riding along with a bug fix.
var resourceScopeExempt = map[string]bool{
	"compose_page_lookup":            true,
	"compose_page_create":            true,
	"compose_page_update":            true,
	"compose_page_delete":            true,
	"compose_page_undelete":          true,
	"compose_page_remove_blocks":     true,
	"compose_page_reorder":           true,
	"compose_page_block_schema":      true,
	"compose_chart_lookup":           true,
	"compose_chart_create":           true,
	"compose_chart_update":           true,
	"compose_chart_delete":           true,
	"compose_chart_undelete":         true,
	"automation_taq_lookup":          true,
	"automation_taq_executions":      true,
	"automation_taq_execution_trace": true,
	"automation_workflow_lookup":     true,

	// Undelete restores a definition rather than running one, so it is
	// classified with its lookup sibling and not with exec. Mapping it in
	// buildResource instead would route it through checkAllow, which denies a
	// mapped resource whenever the agent has no allow entries — a restriction
	// exec carries deliberately and a configuring op has no reason to inherit.
	"automation_taq_undelete":      true,
	"automation_workflow_undelete": true,
	"discovery_search":             true,

	// A trigger carries no compose resource dimension: it is scoped by the
	// workflow it fires, and canManageTrigger checks CanManageTriggersOnWorkflow
	// against that workflow on every write. Event types are a compile-time
	// catalogue, identical for every caller.
	//
	// Listed by name rather than by prefix because the prefix list must not gain
	// an automation_ entry — that space is shared with the runtime-minted
	// per-TAQ tools.
	"automation_trigger_lookup":    true,
	"automation_trigger_create":    true,
	"automation_trigger_update":    true,
	"automation_trigger_delete":    true,
	"automation_trigger_undelete":  true,
	"automation_event_type_lookup": true,

	// TAQ authoring. Classified with lookup rather than with exec for the same
	// reason undelete is: it changes what an automation is, not which one runs.
	// A TAQ carries no compose namespace or module dimension for checkAllow to
	// narrow, and the service gates both writes itself —
	// CanCreateNgAutomation on create, CanUpdateNgAutomation inside
	// handleUpdate.
	"automation_taq_create": true,
	"automation_taq_update": true,

	// TAQ and workflow delete. Their exec siblings are mapped in buildResource
	// because running one acts on a specific automation; deleting is a
	// definition-level change gated by CanDeleteNgAutomation /
	// CanDeleteWorkflow in the service.
	"automation_taq_delete":      true,
	"automation_workflow_delete": true,

	// Workflow create and update, for the same reason: authoring a definition is
	// not running one. CanCreateWorkflow is a component-level check with no
	// workflow to narrow on at all, and CanUpdateWorkflow is checked in the
	// service against the loaded workflow. Mapping either in buildResource would
	// route it through checkAllow, which denies whenever the agent has no allow
	// entries — a restriction exec carries deliberately and authoring has no
	// reason to inherit.
	"automation_workflow_create": true,
	"automation_workflow_update": true,

	// Reminders carry no compose resource dimension — they are scoped to their
	// assignee inside the service (onLookup refuses a reminder assigned to
	// someone else, onSearch filters by the same predicate), which is an
	// authorization model checkAllow has nothing to narrow.
	"system_reminder_lookup":    true,
	"system_reminder_create":    true,
	"system_reminder_update":    true,
	"system_reminder_delete":    true,
	"system_reminder_dismiss":   true,
	"system_reminder_undismiss": true,
	"system_reminder_snooze":    true,
}

// resourceScopeExemptPrefixes covers whole tool families that carry no
// compose/automation resource dimension.
//
// The per-name map below stays the default, because naming a tool is what
// forces someone to think about its scoping. A prefix is only right where the
// whole resource is out of the compose/automation dimension by construction —
// a user or a role is not scoped to a namespace, and never will be — and it
// keeps a 40-tool batch from burying the one-off entries that do carry meaning.
//
// Do not add a compose_ or automation_ prefix here. Those resources have a real
// dimension and belong in buildResource.
var resourceScopeExemptPrefixes = []string{
	"system_user_",
	"system_user_group_",
	"system_role_",
	"system_auth_client_",
	"system_application_",
	"system_reminder_",
}

// IsClassified reports whether a tool has been given a resource mapping or been
// explicitly marked as carrying no resource dimension. It exists so a test in
// the mcp package can assert that every registered tool is one or the other —
// see the comment in Evaluate for why this is a CI-time check and not a runtime
// denial.
//
// Runtime-minted per-TAQ tools are classified by construction.
func IsClassified(tool string) bool {
	if _, ok := dynamicTAQRef(tool); ok {
		return true
	}
	if resourceScopeExempt[tool] {
		return true
	}
	for _, prefix := range resourceScopeExemptPrefixes {
		if strings.HasPrefix(tool, prefix) {
			return true
		}
	}
	_, mapped := buildResource(tool, MapValues{})
	return mapped
}

// IsScopeExempt reports whether a tool is exempt from resource-level narrowing,
// by name or by family prefix. Evaluate does not branch on this — an unmapped
// tool is allowed either way, see the comment there — but it makes the two
// exemption mechanisms testable as one thing.
func IsScopeExempt(tool string) bool {
	if resourceScopeExempt[tool] {
		return true
	}
	for _, prefix := range resourceScopeExemptPrefixes {
		if strings.HasPrefix(tool, prefix) {
			return true
		}
	}
	return false
}

// buildResource constructs a Human resource identifier from the tool name and
// args. Missing or zero-value segments are replaced with "*". The second return
// value reports whether the tool is known to this mapping at all.
func buildResource(tool string, args ValueGetter) (string, bool) {
	seg := func(key string) string {
		v, ok := args.Get(key)
		if !ok {
			return "*"
		}
		s := fmt.Sprintf("%v", v)
		if s == "" || s == "0" {
			return "*"
		}
		return s
	}

	switch {
	case strings.HasPrefix(tool, "compose_record_"):
		return fmt.Sprintf("corteza::compose:record/%s/%s/%s", seg("namespaceID"), seg("moduleID"), seg("recordID")), true
	case strings.HasPrefix(tool, "compose_module_"):
		return fmt.Sprintf("corteza::compose:module/%s/%s", seg("namespaceID"), seg("moduleID")), true
	case strings.HasPrefix(tool, "compose_namespace_"):
		return fmt.Sprintf("corteza::compose:namespace/%s", seg("namespaceID")), true
	case tool == "automation_taq_exec":
		return fmt.Sprintf("corteza::automation:ng-automation/%s", seg("taq")), true
	case tool == "automation_workflow_exec":
		return fmt.Sprintf("corteza::automation:workflow/%s", seg("workflow")), true
	default:
		return "", false
	}
}

// checkAllow returns denied if allow is non-empty and no entry covers the resource.
// Each allow entry covers a namespace; if ModuleIDs is empty it covers all modules in that namespace.
func checkAllow(allow []types.AgentAccessAllow, resource string) Decision {
	if resource == "" {
		return Decision{Allowed: true}
	}
	if len(allow) == 0 {
		return Decision{Allowed: false, Reason: "tool has no allow entries"}
	}

	if !strings.HasPrefix(resource, "corteza::compose:") {
		return Decision{Allowed: true}
	}

	parts := strings.Split(resource, "/")
	if len(parts) < 2 {
		return Decision{Allowed: true}
	}
	nsStr := parts[1]
	if nsStr == "*" {
		// The executor already errors out early when the LLM provides an invalid namespace,
		// so reaching "*" here means no namespace was provided at all.
		// Deny — when an explicit allow-list is present a namespace must be specified.
		return Decision{Allowed: false, Reason: "namespace is required but was not specified"}
	}
	modStr := ""
	if len(parts) > 2 {
		modStr = parts[2]
	}

	for _, a := range allow {
		if fmt.Sprintf("%d", a.NamespaceID) != nsStr {
			continue
		}
		if len(a.ModuleIDs) == 0 {
			return Decision{Allowed: true}
		}
		if modStr == "" || modStr == "*" {
			return Decision{Allowed: false, Reason: "module is required but was not specified"}
		}
		for _, mid := range a.ModuleIDs {
			if fmt.Sprintf("%d", mid) == modStr {
				return Decision{Allowed: true}
			}
		}
	}

	return Decision{Allowed: false, Reason: "resource not in allow-list"}
}
