package types

import (
	"testing"

	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/stretchr/testify/require"
)

func TestTriggerResourceRefs(t *testing.T) {
	tr := Trigger{
		WorkflowID:   600,
		ResourceType: "compose:record",
		Constraints: TriggerConstraintSet{
			{Name: "module", Values: []string{"1001", "lead-module"}},
			{Name: "namespace", Values: []string{"crm"}},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindAutomationWorkflow, ID: 600, Reason: resourceref.ReasonTriggerWorkflow, Path: "WorkflowID"},
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonTriggerModule, Path: "Constraints.0.Values.0"},
		{Kind: resourceref.KindComposeModule, Ident: "lead-module", Reason: resourceref.ReasonTriggerModule, Path: "Constraints.0.Values.1"},
	}, tr.ResourceRefs())
}

func TestTriggerResourceRefsNonComposeResource(t *testing.T) {
	tr := Trigger{
		ResourceType: "system:user",
		Constraints: TriggerConstraintSet{
			{Name: "module", Values: []string{"1001"}},
		},
	}

	require.Empty(t, tr.ResourceRefs())
}

func TestWorkflowResourceRefsSteps(t *testing.T) {
	w := Workflow{
		Steps: WorkflowStepSet{
			// constant module handle + dynamic namespace expression
			{Kind: WorkflowStepKindFunction, Ref: "composeRecordsSearch", Arguments: []*Expr{
				{Target: "module", Value: "lead"},
				{Target: "namespace", Expr: "scope.ns"},
				{Target: "query", Value: "x"},
			}},
			// subworkflow exec by ID
			{Kind: WorkflowStepKindExecWorkflow, Arguments: []*Expr{
				{Target: "workflow", Value: "601"},
			}},
			// unknown function: nothing extracted, no warnings
			{Kind: WorkflowStepKindFunction, Ref: "logInfo", Arguments: []*Expr{
				{Target: "message", Expr: "scope.msg"},
			}},
			// agent run with constant ID
			{Kind: WorkflowStepKindFunction, Ref: "agentRun", Arguments: []*Expr{
				{Target: "agentID", Value: "6001"},
			}},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindComposeModule, Ident: "lead", Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
		{Kind: resourceref.KindComposeNamespace, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.1", Dynamic: true},
		{Kind: resourceref.KindAutomationWorkflow, ID: 601, Reason: resourceref.ReasonStepArgument, Path: "Steps.1.Arguments.0"},
		{Kind: resourceref.KindAgent, ID: 6001, Reason: resourceref.ReasonStepArgument, Path: "Steps.3.Arguments.0"},
	}, w.ResourceRefs())
}

func TestNgAutomationResourceRefsSteps(t *testing.T) {
	a := NgAutomation{
		Steps: NgAutomationStepSet{
			{Kind: "function", Ref: "composeRecordsCreate", Arguments: []*Expr{
				{Target: "module", Value: "1001"},
				{Target: "namespace", Source: "ns"},
			}},
			// configured-connection function; connection ID lives in the ref
			{Kind: "function", Ref: "conn_4001_sheetAppend", Arguments: []*Expr{
				{Target: "range", Value: "A1"},
			}},
			{Kind: "function", Ref: "agentPrompt", Arguments: []*Expr{
				{Target: "agentID", Value: "6001"},
			}},
			{Kind: "function", Ref: "notificationSendRecord", Arguments: []*Expr{
				{Target: "module", Value: "lead"},
			}},
			{Kind: "function", Ref: "rolesAddMember", Arguments: []*Expr{
				{Target: "role", Value: "8001"},
				{Target: "user", Value: "42"},
			}},
			{Kind: "function", Ref: "templatesRender", Arguments: []*Expr{
				{Target: "lookup", Value: "invoice-tpl"},
			}},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
		{Kind: resourceref.KindComposeNamespace, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.1", Dynamic: true},
		{Kind: resourceref.KindConfiguredConnection, ID: 4001, Reason: resourceref.ReasonStepConnection, Path: "Steps.1.Ref"},
		{Kind: resourceref.KindAgent, ID: 6001, Reason: resourceref.ReasonStepArgument, Path: "Steps.2.Arguments.0"},
		{Kind: resourceref.KindComposeModule, Ident: "lead", Reason: resourceref.ReasonStepArgument, Path: "Steps.3.Arguments.0"},
		{Kind: resourceref.KindRole, ID: 8001, Reason: resourceref.ReasonStepArgument, Path: "Steps.4.Arguments.0"},
		{Kind: resourceref.KindTemplate, Ident: "invoice-tpl", Reason: resourceref.ReasonStepArgument, Path: "Steps.5.Arguments.0"},
	}, a.ResourceRefs())
}

func TestNgAutomationStepResolution(t *testing.T) {
	step := func(kind, ref string, args ...*Expr) *NgAutomationStep {
		return &NgAutomationStep{Kind: kind, Ref: ref, Arguments: args}
	}

	cc := []struct {
		name string
		step *NgAutomationStep
		out  []resourceref.Ref
	}{
		{
			"composeRecords prefix covers every record function",
			step("iterator", "composeRecordsEach",
				&Expr{Target: "module", Value: "1001"},
				&Expr{Target: "namespace", Value: "55"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
				{Kind: resourceref.KindComposeNamespace, ID: 55, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.1"},
			},
		},
		{
			"numeric JSON constant",
			step("function", "composeRecordsCreate", &Expr{Target: "module", Value: float64(1001)}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
		{
			"constant handle kept as ident",
			step("function", "composeRecordsLookup", &Expr{Target: "module", Value: "lead"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, Ident: "lead", Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
		{
			"zero value, no expression: nothing",
			step("function", "composeRecordsLookup", &Expr{Target: "module", Value: "0"}),
			nil,
		},
		{
			"expression arg: dynamic warning ref",
			step("function", "composeRecordsSearch", &Expr{Target: "module", Expr: "scope.mod"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0", Dynamic: true},
			},
		},
		{
			"scope variable arg: dynamic warning ref",
			step("function", "composeRecordsSearch", &Expr{Target: "module", Source: "mod"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0", Dynamic: true},
			},
		},
		{
			"constant value wins over expression",
			step("function", "composeRecordsSearch", &Expr{Target: "module", Value: "1001", Expr: "scope.mod"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
		{
			"non-resource args ignored even when dynamic",
			step("function", "composeRecordsSearch",
				&Expr{Target: "query", Expr: "scope.q"},
				&Expr{Target: "limit", Value: "10"}),
			nil,
		},
		{
			"unknown function: nothing extracted",
			step("function", "httpRequestSend", &Expr{Target: "module", Value: "1001"}),
			nil,
		},
		{
			"connection function: ID from ref itself",
			step("function", "conn_4001_sheetAppend", &Expr{Target: "range", Value: "A1"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindConfiguredConnection, ID: 4001, Reason: resourceref.ReasonStepConnection, Path: "Steps.0.Ref"},
			},
		},
		{
			"malformed connection refs: nothing",
			step("function", "conn_abc_op"),
			nil,
		},
		{
			"connection ref without operation: nothing",
			step("function", "conn_4001"),
			nil,
		},
		{
			"nil argument tolerated",
			step("function", "composeRecordsLookup", nil, &Expr{Target: "module", Value: "1001"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.1"},
			},
		},
		{
			"roles: lookup and role extracted, user skipped",
			step("function", "rolesAddMember",
				&Expr{Target: "role", Value: "8001"},
				&Expr{Target: "user", Value: "42"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindRole, ID: 8001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
		{
			"templates lookup by handle",
			step("function", "templatesRender", &Expr{Target: "lookup", Value: "invoice-tpl"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindTemplate, Ident: "invoice-tpl", Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
		{
			"agentPrompt by ID",
			step("function", "agentPrompt", &Expr{Target: "agentID", Value: "6001"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindAgent, ID: 6001, Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
		{
			"notificationSendRecord module",
			step("function", "notificationSendRecord", &Expr{Target: "module", Value: "lead"}),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, Ident: "lead", Reason: resourceref.ReasonStepArgument, Path: "Steps.0.Arguments.0"},
			},
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			a := NgAutomation{Steps: NgAutomationStepSet{c.step}}

			var want []resourceref.Ref
			want = append(want, c.out...)

			require.Equal(t, want, a.ResourceRefs())
		})
	}
}

func TestNgAutomationTriggerResolution(t *testing.T) {
	trigger := func(rt string, cc ...NgTriggerConstraint) *NgAutomationTrigger {
		return &NgAutomationTrigger{ResourceType: rt, Constraints: cc}
	}
	constraint := func(name string, vv ...string) NgTriggerConstraint {
		c := NgTriggerConstraint{Name: name}
		for _, v := range vv {
			c.Values = append(c.Values, NgTriggerConstraintValue{Value: v})
		}
		return c
	}

	cc := []struct {
		name    string
		trigger *NgAutomationTrigger
		out     []resourceref.Ref
	}{
		{
			"module constraint by ID and handle",
			trigger("compose:record", constraint("module", "1001", "lead")),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonTriggerModule, Path: "Triggers.0.Constraints.0.Values.0"},
				{Kind: resourceref.KindComposeModule, Ident: "lead", Reason: resourceref.ReasonTriggerModule, Path: "Triggers.0.Constraints.0.Values.1"},
			},
		},
		{
			"module.handle constraint accepted",
			trigger("compose:record", constraint("module.handle", "lead")),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, Ident: "lead", Reason: resourceref.ReasonTriggerModule, Path: "Triggers.0.Constraints.0.Values.0"},
			},
		},
		{
			"long-form resource type accepted",
			trigger("corteza::compose:record", constraint("module", "1001")),
			[]resourceref.Ref{
				{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonTriggerModule, Path: "Triggers.0.Constraints.0.Values.0"},
			},
		},
		{
			"non-compose resource type: nothing",
			trigger("system:user", constraint("module", "1001")),
			nil,
		},
		{
			"non-module constraint: nothing",
			trigger("compose:record", constraint("namespace", "crm"), constraint("record.values.status", "open")),
			nil,
		},
		{
			"empty constraint values: nothing",
			trigger("compose:record", constraint("module")),
			nil,
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			a := NgAutomation{Triggers: NgAutomationTriggerSet{c.trigger}}

			var want []resourceref.Ref
			want = append(want, c.out...)

			require.Equal(t, want, a.ResourceRefs())
		})
	}
}

func TestNgAutomationResourceRefs(t *testing.T) {
	a := NgAutomation{
		Triggers: NgAutomationTriggerSet{
			{
				ResourceType: "compose:record",
				Constraints: []NgTriggerConstraint{
					{Name: "module.handle", Values: []NgTriggerConstraintValue{
						{Type: "String", Value: "lead-module"},
						{Type: "ID", Value: "1001"},
					}},
				},
			},
			nil,
			{
				ResourceType: "system:user",
				Constraints: []NgTriggerConstraint{
					{Name: "module", Values: []NgTriggerConstraintValue{{Value: "1002"}}},
				},
			},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindComposeModule, Ident: "lead-module", Reason: resourceref.ReasonTriggerModule, Path: "Triggers.0.Constraints.0.Values.0"},
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonTriggerModule, Path: "Triggers.0.Constraints.0.Values.1"},
	}, a.ResourceRefs())
}
