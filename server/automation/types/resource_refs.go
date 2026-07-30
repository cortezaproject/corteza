package types

import (
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/spf13/cast"
)

// resourceRefsExt extends the generated Trigger.ResourceRefs() with the
// Constraints module-scope refs that require an isModuleConstraint check.
func (t Trigger) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	for _, c := range t.Constraints {
		if c == nil || !isModuleConstraint(t.ResourceType, c.Name) {
			continue
		}

		for _, v := range c.Values {
			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				v,
				resourceref.ReasonTriggerModule,
			))
		}
	}

	return out
}

// resourceRefsExt extends the generated Workflow.ResourceRefs() with
// step-argument refs.
func (w Workflow) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	for _, s := range w.Steps {
		if s == nil {
			continue
		}

		out = append(out, stepResourceRefs(string(s.Kind), s.Ref, s.Arguments)...)
	}

	return out
}

// resourceRefsExt extends the generated NgAutomation.ResourceRefs() with
// trigger-constraint + step refs.
func (a NgAutomation) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	for _, t := range a.Triggers {
		if t == nil {
			continue
		}

		for _, c := range t.Constraints {
			if !isModuleConstraint(t.ResourceType, c.Name) {
				continue
			}

			for _, v := range c.Values {
				out = resourceref.Append(out, resourceref.MakeIdent(
					resourceref.KindComposeModule,
					v.Value,
					resourceref.ReasonTriggerModule,
				))
			}
		}
	}

	for _, s := range a.Steps {
		if s == nil {
			continue
		}

		out = append(out, stepResourceRefs(s.Kind, s.Ref, s.Arguments)...)
	}

	return out
}

// NgAutomationStepParamKinds returns the map of argument target → resource kind
// for one ng-automation step — the same mapping WorkflowStepParamKinds returns
// for a workflow step, against the ng step shape.
//
// Exported for the project revision copy (system/service/project_revision_clone.go),
// which has to rewrite a copied automation's resource-typed arguments onto the
// draft's own modules, namespace and agents.
func NgAutomationStepParamKinds(s *NgAutomationStep) map[string]string {
	if s == nil {
		return nil
	}
	if s.Kind == string(WorkflowStepKindExecWorkflow) {
		return map[string]string{"workflow": resourceref.KindAutomationWorkflow}
	}
	return stepArgumentKinds(s.Ref)
}

// WorkflowStepParamKinds returns the map of argument target → resource kind
// for the given step. Exported for use by the envoyx encode/decode path.
func WorkflowStepParamKinds(s *WorkflowStep) map[string]string {
	if s == nil {
		return nil
	}
	if s.Kind == WorkflowStepKindExecWorkflow {
		return map[string]string{"workflow": resourceref.KindAutomationWorkflow}
	}
	return stepArgumentKinds(s.Ref)
}

// stepResourceRefs extracts references from a step's arguments
//
// Constant argument values produce resolvable refs; computed arguments
// (expressions, scope variables) targeting resource-typed parameters produce
// unresolved refs which consumers surface as warnings instead of edges.
func stepResourceRefs(kind, ref string, args []*Expr) (out []resourceref.Ref) {
	// connection-generated functions encode the configured connection in the
	// ref itself: conn_{connectionID}_{operation} (see configured connection
	// function registration in system/service/configured_connection.go)
	if id := connFunctionRef(ref); id > 0 {
		out = append(out, resourceref.Make(
			resourceref.KindConfiguredConnection,
			id,
			resourceref.ReasonStepConnection,
		))
	}

	paramKinds := stepArgumentKinds(ref)

	// subworkflow exec carries the target in the "workflow" argument
	// (see workflowConverter.convExecWorkflowStep)
	if kind == string(WorkflowStepKindExecWorkflow) {
		paramKinds = map[string]string{"workflow": resourceref.KindAutomationWorkflow}
	}

	if len(paramKinds) == 0 {
		return
	}

	for _, a := range args {
		if a == nil {
			continue
		}

		k, is := paramKinds[a.Target]
		if !is {
			continue
		}

		if a.Value != nil {
			if r := resourceref.MakeIdent(k, cast.ToString(a.Value), resourceref.ReasonStepArgument); !r.IsEmpty() {
				out = append(out, r)
				continue
			}
		}

		if a.Expr != "" || a.Source != "" {
			out = append(out, resourceref.MakeDynamic(k, resourceref.ReasonStepArgument))
		}
	}

	return
}

// stepArgumentKinds maps function refs to their resource-referencing parameters
//
// Static mirror of the function registry metadata (generated *_handler.gen.go
// files declare param types like "ComposeModule"); extend when new
// resource-bound functions land.
func stepArgumentKinds(ref string) map[string]string {
	switch {
	case strings.HasPrefix(ref, "composeRecords"):
		return map[string]string{
			"module":    resourceref.KindComposeModule,
			"namespace": resourceref.KindComposeNamespace,
		}

	case strings.HasPrefix(ref, "roles"):
		// rolesLookup & friends use "lookup", member ops use "role";
		// "user" params intentionally skipped (users are not graph resources)
		return map[string]string{
			"lookup": resourceref.KindRole,
			"role":   resourceref.KindRole,
		}

	case strings.HasPrefix(ref, "templates"):
		return map[string]string{
			"lookup":   resourceref.KindTemplate,
			"template": resourceref.KindTemplate,
		}
	}

	switch ref {
	case "composeModulesLookup":
		return map[string]string{
			"module":    resourceref.KindComposeModule,
			"namespace": resourceref.KindComposeNamespace,
		}

	case "composeNamespacesLookup":
		return map[string]string{
			"namespace": resourceref.KindComposeNamespace,
		}

	case "agentRun", "agentPrompt":
		return map[string]string{
			"agentID": resourceref.KindAgent,
		}

	case "notificationSendRecord":
		return map[string]string{
			"module":    resourceref.KindComposeModule,
			"namespace": resourceref.KindComposeNamespace,
		}
	}

	return nil
}

// connFunctionRef extracts the configured connection ID from
// connection-generated function refs (conn_{connectionID}_{operation})
func connFunctionRef(ref string) uint64 {
	rest, found := strings.CutPrefix(ref, "conn_")
	if !found {
		return 0
	}

	idStr, _, found := strings.Cut(rest, "_")
	if !found {
		return 0
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)
	return id
}

// isModuleConstraint matches trigger constraints which scope compose events to
// modules (constraint matching: compose/service/event/module.go)
func isModuleConstraint(resourceType, name string) bool {
	if !strings.HasPrefix(resourceType, "compose:") && !strings.HasPrefix(resourceType, "corteza::compose:") {
		return false
	}

	switch name {
	case "module", "module.handle":
		return true
	}

	return false
}
