package envoy

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/envoyx"
	"gopkg.in/yaml.v3"
)

func (e YamlEncoder) encode(ctx context.Context, base *yaml.Node, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	return
}

// encodeWorkflowStepsC rewrites resource-ref argument values to friendly
// identifiers before YAML serialisation. The original steps are copied so the
// in-memory workflow resource is not mutated.
func (e YamlEncoder) encodeWorkflowStepsC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, wf *types.Workflow, steps types.WorkflowStepSet) (_ any, err error) {
	out := make(types.WorkflowStepSet, len(steps))
	for i, s := range steps {
		if s == nil {
			out[i] = s
			continue
		}
		sCopy := *s
		argsCopy := make([]*types.Expr, len(s.Arguments))
		for j, a := range s.Arguments {
			if a == nil {
				argsCopy[j] = nil
				continue
			}
			aCopy := *a
			argsCopy[j] = &aCopy
		}
		sCopy.Arguments = argsCopy

		for j, a := range sCopy.Arguments {
			if a == nil {
				continue
			}
			ref, ok := n.References[fmt.Sprintf("Steps.%d.Arguments.%s", i, a.Target)]
			if !ok {
				continue
			}
			sCopy.Arguments[j].Value = safeParentIdentifier(tt, n, ref)
		}
		out[i] = &sCopy
	}
	return out, nil
}
