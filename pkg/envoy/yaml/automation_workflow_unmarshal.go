package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

func (wset *automationWorkflowSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &automationWorkflow{}
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed automation workflow definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if err = decodeRef(k, "automation workflow handle", &wrap.res.Handle); err != nil {
			return y7s.NodeErr(n, "Automation workflow reference must be a valid handle")
		}

		if wrap.res.Meta.Name == "" {
			// if name is not set, use handle
			wrap.res.Meta.Name = wrap.res.Handle
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *automationWorkflow) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.rbac = make(rbacRuleSet, 0, 10)
		wrap.res = &types.Workflow{}
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "handle":
			return y7s.DecodeScalar(v, "workflow handle", &wrap.res.Handle)
		case "meta":
			return v.Decode(&wrap.res.Meta)
		case "enabled":
			return y7s.DecodeScalar(v, "workflow enabled", &wrap.res.Enabled)
		case "trace":
			return y7s.DecodeScalar(v, "workflow trace", &wrap.res.Trace)
		case "keepSessions":
			return y7s.DecodeScalar(v, "workflow keepSessions", &wrap.res.KeepSessions)
		case "scope":
			err = v.Decode(&wrap.res.Scope)
			if err != nil {
				return err
			}

		case "triggers":
			wrap.triggers = make(automationTriggerSet, 0, 100)

			err = v.Decode(&wrap.triggers)
			if err != nil {
				return err
			}

		case "steps":
			wrap.steps = make(automationWorkflowStepSet, 0, 100)

			err = v.Decode(&wrap.steps)
			if err != nil {
				return err
			}

		case "paths":
			wrap.paths = make(automationWorkflowPathSet, 0, 100)

			err = v.Decode(&wrap.paths)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (wrap *automationTrigger) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.Trigger{}
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "enabled":
			return y7s.DecodeScalar(v, "trigger enabled", &wrap.res.Enabled)
		case "stepID":
			return y7s.DecodeScalar(v, "trigger step", &wrap.res.StepID)
		case "resourceType":
			return y7s.DecodeScalar(v, "trigger resourceType", &wrap.res.ResourceType)
		case "eventType":
			return y7s.DecodeScalar(v, "trigger eventType", &wrap.res.EventType)
		case "constraints":
			return v.Decode(&wrap.res.Constraints)
		case "input":
			err = v.Decode(&wrap.res.Input)
			if err != nil {
				return err
			}

		case "meta":
			return v.Decode(&wrap.res.Meta)
		}

		return nil
	})
}

func (wrap *automationWorkflowStep) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.WorkflowStep{}
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "stepid",
			"id":
			return y7s.DecodeScalar(v, "trigger step", &wrap.res.ID)
		case "kind":
			return y7s.DecodeScalar(v, "step kind", &wrap.res.Kind)
		case "ref":
			return y7s.DecodeScalar(v, "step ref", &wrap.res.Ref)
		case "arguments":
			wrap.res.Arguments, err = unmarshalExprSet(v)
			return err
		case "results":
			wrap.res.Results, err = unmarshalExprSet(v)
			return err
		case "meta":
			return v.Decode(&wrap.res.Meta)
		}

		return nil
	})
}

func (wrap *automationWorkflowPath) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.WorkflowPath{}
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "expr":
			return y7s.DecodeScalar(v, "path expr", &wrap.res.Expr)
		case "parentid",
			"parent":
			return y7s.DecodeScalar(v, "parent ref", &wrap.res.ParentID)
		case "childid",
			"child":
			return y7s.DecodeScalar(v, "child ref", &wrap.res.ChildID)
		case "meta":
			return v.Decode(&wrap.res.Meta)
		}

		return nil
	})
}

func unmarshalExprSet(n *yaml.Node) ([]*types.Expr, error) {
	ee := make([]*types.Expr, 0, 10)

	err := y7s.EachSeq(n, func(v *yaml.Node) (err error) {
		wrap, err := unmarshalExpr(v)
		ee = append(ee, wrap)
		return err
	})

	return ee, err
}

func unmarshalExpr(n *yaml.Node) (*types.Expr, error) {
	wrap := &types.Expr{}

	err := y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "target":
			return y7s.DecodeScalar(v, "expression target", &wrap.Target)
		case "source":
			return y7s.DecodeScalar(v, "expression source", &wrap.Source)
		case "expr":
			return y7s.DecodeScalar(v, "expression expr", &wrap.Expr)
		case "value":
			return y7s.DecodeScalar(v, "expression value", &wrap.Value)
		case "type":
			return y7s.DecodeScalar(v, "expression type", &wrap.Type)
		case "tests":
			tt := make(types.TestSet, 0, 2)
			err = v.Decode(&tt)
			if err != nil {
				return err
			}
			wrap.Tests = tt
			return nil
		}

		return nil
	})

	return wrap, err
}

func (wset automationWorkflowSet) MarshalEnvoy() ([]resource.Interface, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]resource.Interface, 0, len(wset)*10)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap automationWorkflow) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewAutomationWorkflow(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetUserstamps(wrap.us)
	rs.SetConfig(wrap.envoyConfig)

	for _, t := range wrap.triggers {
		trs := rs.AddAutomationTrigger(t.res)
		trs.SetTimestamps(t.ts)
		trs.SetUserstamps(t.us)
	}

	for _, s := range wrap.steps {
		rs.AddAutomationWorkflowStep(s.res)
	}

	for _, p := range wrap.paths {
		rs.AddAutomationWorkflowPath(p.res)
	}

	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
