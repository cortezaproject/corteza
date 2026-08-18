package runtime

import (
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/expr"
	"go.uber.org/zap"
)

func mkResolveRuntime(t *testing.T) *runtime {
	t.Helper()

	global, err := expr.NewVars(map[string]any{"greeting": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return &runtime{
		globalState: global,
		entryPoint:  "trigger_1",
		scheduler:   &scheduler{},
		log:         zap.NewNop(),
	}
}

func argScoped(scope string) types.StepArg {
	return types.StepArg{Expr: &types.Expr{ArgumentName: "a", Scope: scope, Source: "greeting"}}
}

// TL;DR: the reserved scope names are served from the global state, never looked up among step outputs.
// Example: identities the server injects (invoker, runner) live in the global scope, so a step
// referencing them names scope "global" — which no step ever produces an output under.
func TestResolveInputs_ReservedScopesComeFromGlobalState(t *testing.T) {
	for _, scope := range []string{"", "global", "trigger_1"} {
		r := mkResolveRuntime(t)
		step := &types.Step{Arguments: []types.StepArg{argScoped(scope)}}

		out, err := r.resolveInputs(step)
		if err != nil {
			t.Fatalf("scope %q: unexpected error: %v", scope, err)
		}
		if out[scope] != r.globalState {
			t.Errorf("scope %q: expected the global state, got %v", scope, out[scope])
		}
	}
}

// TL;DR: an unknown scope is still an error — the reserved-name shortcut must not swallow a real miss.
// Example: a step pulling from a sibling that was deleted should fail loudly, not read the global scope.
func TestResolveInputs_UnknownScopeStillFails(t *testing.T) {
	r := mkResolveRuntime(t)
	step := &types.Step{Arguments: []types.StepArg{argScoped("step_9")}}

	if _, err := r.resolveInputs(step); err == nil {
		t.Fatal("expected an error for a scope no step produces")
	}
}
