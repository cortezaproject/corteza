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
		resourceref.Make(resourceref.KindAutomationWorkflow, 600, resourceref.ReasonTriggerWorkflow),
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonTriggerModule),
		resourceref.MakeIdent(resourceref.KindComposeModule, "lead-module", resourceref.ReasonTriggerModule),
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
		resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonStepArgument),
		resourceref.MakeDynamic(resourceref.KindComposeNamespace, resourceref.ReasonStepArgument),
		resourceref.Make(resourceref.KindAutomationWorkflow, 601, resourceref.ReasonStepArgument),
		resourceref.Make(resourceref.KindAgent, 6001, resourceref.ReasonStepArgument),
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
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonStepArgument),
		resourceref.MakeDynamic(resourceref.KindComposeNamespace, resourceref.ReasonStepArgument),
		resourceref.Make(resourceref.KindConfiguredConnection, 4001, resourceref.ReasonStepConnection),
		resourceref.Make(resourceref.KindAgent, 6001, resourceref.ReasonStepArgument),
		resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonStepArgument),
		resourceref.Make(resourceref.KindRole, 8001, resourceref.ReasonStepArgument),
		resourceref.MakeIdent(resourceref.KindTemplate, "invoice-tpl", resourceref.ReasonStepArgument),
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
				resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonStepArgument),
				resourceref.Make(resourceref.KindComposeNamespace, 55, resourceref.ReasonStepArgument),
			},
		},
		{
			"numeric JSON constant",
			step("function", "composeRecordsCreate", &Expr{Target: "module", Value: float64(1001)}),
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonStepArgument),
			},
		},
		{
			"constant handle kept as label",
			step("function", "composeRecordsLookup", &Expr{Target: "module", Value: "lead"}),
			[]resourceref.Ref{
				resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonStepArgument),
			},
		},
		{
			"zero value, no expression: nothing",
			step("function", "composeRecordsLookup", &Expr{Target: "module", Value: "0"}),
			nil,
		},
		{
			"expression arg: unresolved warning ref",
			step("function", "composeRecordsSearch", &Expr{Target: "module", Expr: "scope.mod"}),
			[]resourceref.Ref{
				resourceref.MakeDynamic(resourceref.KindComposeModule, resourceref.ReasonStepArgument),
			},
		},
		{
			"scope variable arg: unresolved warning ref",
			step("function", "composeRecordsSearch", &Expr{Target: "module", Source: "mod"}),
			[]resourceref.Ref{
				resourceref.MakeDynamic(resourceref.KindComposeModule, resourceref.ReasonStepArgument),
			},
		},
		{
			"constant value wins over expression",
			step("function", "composeRecordsSearch", &Expr{Target: "module", Value: "1001", Expr: "scope.mod"}),
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonStepArgument),
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
				resourceref.Make(resourceref.KindConfiguredConnection, 4001, resourceref.ReasonStepConnection),
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
				resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonStepArgument),
			},
		},
		{
			"roles: lookup and role extracted, user skipped",
			step("function", "rolesAddMember",
				&Expr{Target: "role", Value: "8001"},
				&Expr{Target: "user", Value: "42"}),
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindRole, 8001, resourceref.ReasonStepArgument),
			},
		},
		{
			"templates lookup by handle",
			step("function", "templatesRender", &Expr{Target: "lookup", Value: "invoice-tpl"}),
			[]resourceref.Ref{
				resourceref.MakeIdent(resourceref.KindTemplate, "invoice-tpl", resourceref.ReasonStepArgument),
			},
		},
		{
			"agentPrompt by ID",
			step("function", "agentPrompt", &Expr{Target: "agentID", Value: "6001"}),
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindAgent, 6001, resourceref.ReasonStepArgument),
			},
		},
		{
			"notificationSendRecord module",
			step("function", "notificationSendRecord", &Expr{Target: "module", Value: "lead"}),
			[]resourceref.Ref{
				resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonStepArgument),
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
				resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonTriggerModule),
				resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonTriggerModule),
			},
		},
		{
			"module.handle constraint accepted",
			trigger("compose:record", constraint("module.handle", "lead")),
			[]resourceref.Ref{
				resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonTriggerModule),
			},
		},
		{
			"long-form resource type accepted",
			trigger("corteza::compose:record", constraint("module", "1001")),
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonTriggerModule),
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
		resourceref.MakeIdent(resourceref.KindComposeModule, "lead-module", resourceref.ReasonTriggerModule),
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonTriggerModule),
	}, a.ResourceRefs())
}
