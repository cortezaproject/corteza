package envoy

import (
	"fmt"

	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

func (d *auxYamlDoc) unmarshalTriggersExtendedNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	return d.unmarshalTriggerNode(dctx, n, meta...)
}

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) { return }

// unmarshalWorkflowStepsNode extracts resource-ref arguments from YAML-decoded
// steps and registers them as path-keyed envoy refs so the YAML encoder and
// store encoder can rewrite them to friendly identifiers / resolved IDs.
//
// Called by the generated decoder after n.Decode(&r) has populated r.Steps.
func (d *auxYamlDoc) unmarshalWorkflowStepsNode(r *types.Workflow, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = make(map[string]envoyx.Ref)
	for i, s := range r.Steps {
		if s == nil {
			continue
		}
		paramKinds := types.WorkflowStepParamKinds(s)
		if len(paramKinds) == 0 {
			continue
		}
		for _, a := range s.Arguments {
			if a == nil || a.Value == nil {
				continue
			}
			kind, ok := paramKinds[a.Target]
			if !ok {
				continue
			}
			val := cast.ToString(a.Value)
			if val == "" || val == "0" {
				continue
			}
			refs[fmt.Sprintf("Steps.%d.Arguments.%s", i, a.Target)] = envoyx.Ref{
				ResourceType: kind,
				Identifiers:  envoyx.MakeIdentifiers(val),
			}
		}
	}
	return
}
