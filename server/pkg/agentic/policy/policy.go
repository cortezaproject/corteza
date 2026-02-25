package policy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

type Decision struct {
	Allowed       bool
	Reason        string
	SanitizedArgs map[string]any
}

func Evaluate(agent *types.Agent, tool string, args map[string]any) Decision {
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

	sanitized := make(map[string]any, len(args))
	for k, v := range args {
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

func FilterResponse(ctx context.Context, agent *types.Agent, resource string, data map[string]any) map[string]any {
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
		return data
	}

	allowed := make(map[string]bool, len(entry.Properties))
	for _, p := range entry.Properties {
		if p.Access == "allow" || p.Access == "" {
			allowed[p.Name] = true
		}
	}

	result := make(map[string]any, len(allowed))
	for k, v := range data {
		if allowed[k] {
			result[k] = v
		}
	}
	return result
}

func evalFilter(ctx context.Context, filter string, data map[string]any) (bool, error) {
	parser := expr.NewParser()
	evaluable, err := parser.Parse(filter)
	if err != nil {
		return false, fmt.Errorf("invalid filter expression %q: %w", filter, err)
	}

	vars, err := expr.NewVars(data)
	if err != nil {
		return false, fmt.Errorf("failed to build filter vars: %w", err)
	}

	return evaluable.Test(ctx, vars)
}
