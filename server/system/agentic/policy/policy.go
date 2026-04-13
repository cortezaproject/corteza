package policy

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/system/types"
)

// toolAliases mirrors the alias map in the mcp package.
// Keep both in sync when renaming tools.
var toolAliases = map[string]string{
	"compose_namespace_list": "compose_namespace_lookup",
	"compose_module_list":    "compose_module_lookup",
}

type (
	ValueGetter interface {
		Get(key string) (any, bool)
		Keys() []string
	}
	
	ValueSetter interface {
		Set(key string, val any)
	}

	MapValues map[string]any
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

type Decision struct {
	Allowed       bool
	Reason        string
	SanitizedArgs map[string]any
}

func Evaluate(agent *types.Agent, tool string, args ValueGetter) Decision {
	// TAQ and workflow tools are auto-injected — validate against the agent's allowlist
	if tool == "automation_taq_exec" {
		ref, _ := args.Get("taq")
		if findTAQ(agent, fmt.Sprintf("%v", ref)) == nil {
			return Decision{Allowed: false, Reason: fmt.Sprintf("agent is not allowed to execute TAQ %q", ref)}
		}
		return allowedDecision(agent, nil, args)
	}
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

	if d := checkAllow(entry.Allow, buildResource(tool, args)); !d.Allowed {
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

// buildResource constructs a Corteza resource identifier from the tool name and args.
// Missing or zero-value segments are replaced with "*".
func buildResource(tool string, args ValueGetter) string {
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
		return fmt.Sprintf("corteza::compose:record/%s/%s/%s", seg("namespaceID"), seg("moduleID"), seg("recordID"))
	case strings.HasPrefix(tool, "compose_module_"):
		return fmt.Sprintf("corteza::compose:module/%s/%s", seg("namespaceID"), seg("moduleID"))
	case strings.HasPrefix(tool, "compose_namespace_"):
		return fmt.Sprintf("corteza::compose:namespace/%s", seg("namespaceID"))
	case tool == "automation_taq_exec":
		return fmt.Sprintf("corteza::automation:ng-automation/%s", seg("taq"))
	case tool == "automation_workflow_exec":
		return fmt.Sprintf("corteza::automation:workflow/%s", seg("workflow"))
	default:
		return ""
	}
}

// checkAllow returns denied if allow is non-empty and no entry covers the resource.
// Each allow entry covers a namespace; if ModuleIDs is empty it covers all modules in that namespace.
func checkAllow(allow []types.AgentAccessAllow, resource string) Decision {
	if len(allow) == 0 || resource == "" {
		return Decision{Allowed: true}
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
		// Wildcard namespace means a listing/discovery operation (no specific namespace
		// was provided). Allow it — the tool is already in the agent's allow-list.
		return Decision{Allowed: true}
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
			return Decision{Allowed: true}
		}
		for _, mid := range a.ModuleIDs {
			if fmt.Sprintf("%d", mid) == modStr {
				return Decision{Allowed: true}
			}
		}
	}

	return Decision{Allowed: false, Reason: "resource not in allow-list"}
}

func matchResource(pattern, actual []string) bool {
	for i, p := range pattern {
		if i >= len(actual) {
			return false
		}
		if p != "*" && p != actual[i] {
			return false
		}
	}
	return true
}
