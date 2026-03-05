package policy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

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
	var entry *types.AgentAccessTool
	for i := range agent.Access.Tools {
		if agent.Access.Tools[i].Name == tool {
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

	for k, v := range entry.Context.Defaults {
		if _, exists := sanitized[k]; !exists {
			sanitized[k] = v
		}
	}

	for k, v := range entry.Context.Overrides {
		sanitized[k] = v
	}

	return Decision{
		Allowed:       true,
		Reason:        "tool is in agent's allow-list",
		SanitizedArgs: sanitized,
	}
}

func FilterResponse(ctx context.Context, agent *types.Agent, resource string, data ValueGetter) map[string]any {
	var entry *types.AgentAccessAllow
	for i := range agent.Access.Allow {
		if agent.Access.Allow[i].Resource == resource {
			entry = &agent.Access.Allow[i]
			break
		}
	}

	if entry == nil {
		return map[string]any{}
	}

	if entry.Filter != "" {
		match, err := evalFilter(ctx, entry.Filter, data)
		if err != nil || !match {
			return map[string]any{}
		}
	}

	if len(entry.Properties) == 0 {
		result := make(map[string]any)
		for _, k := range data.Keys() {
			v, _ := data.Get(k)
			result[k] = v
		}
		return result
	}

	allowed := make(map[string]bool, len(entry.Properties))
	for _, p := range entry.Properties {
		if p.Access == "allow" || p.Access == "" {
			allowed[p.Name] = true
		}
	}

	result := make(map[string]any, len(allowed))
	for _, k := range data.Keys() {
		if allowed[k] {
			v, _ := data.Get(k)
			result[k] = v
		}
	}
	return result
}

func evalFilter(ctx context.Context, filter string, data ValueGetter) (bool, error) {
	parser := expr.NewParser()
	evaluable, err := parser.Parse(filter)
	if err != nil {
		return false, fmt.Errorf("invalid filter expression %q: %w", filter, err)
	}

	dataMap := make(map[string]any)
	for _, k := range data.Keys() {
		v, _ := data.Get(k)
		dataMap[k] = v
	}

	vars, err := expr.NewVars(dataMap)
	if err != nil {
		return false, fmt.Errorf("failed to build filter vars: %w", err)
	}

	return evaluable.Test(ctx, vars)
}
