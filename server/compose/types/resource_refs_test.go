package types

import (
	"testing"

	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/stretchr/testify/require"
)

func TestModuleResourceRefs(t *testing.T) {
	m := Module{
		Config: ModuleConfig{
			DAL: ModuleConfigDAL{ConnectionID: 42},
		},
		Fields: ModuleFieldSet{
			&ModuleField{Name: "lead", Kind: "Record", Options: ModuleFieldOptions{"moduleID": "1001"}},
			&ModuleField{Name: "title", Kind: "String", Options: ModuleFieldOptions{}},
		},
	}

	require.Equal(t, []resourceref.Ref{
		resourceref.Make(resourceref.KindDalConnection, 42, resourceref.ReasonModuleConnection),
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonModuleFieldRef),
	}, m.ResourceRefs())
}

func TestModuleFieldResourceRefs(t *testing.T) {
	cc := []struct {
		name string
		opt  ModuleFieldOptions
		out  []resourceref.Ref
	}{
		{
			"ID ref",
			ModuleFieldOptions{"moduleID": "1001"},
			[]resourceref.Ref{resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonModuleFieldRef)},
		},
		{
			"handle ref",
			ModuleFieldOptions{"moduleID": "lead-module"},
			[]resourceref.Ref{resourceref.MakeIdent(resourceref.KindComposeModule, "lead-module", resourceref.ReasonModuleFieldRef)},
		},
		{"empty", ModuleFieldOptions{}, nil},
		{"zero", ModuleFieldOptions{"moduleID": "0"}, nil},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, ModuleField{Options: c.opt}.ResourceRefs())
		})
	}
}

func TestChartResourceRefs(t *testing.T) {
	c := Chart{
		Config: ChartConfig{
			Reports: []*ChartConfigReport{
				{ModuleID: 1001},
				nil,
				{ModuleID: 0},
				{ModuleID: 1002},
			},
		},
	}

	require.Equal(t, []resourceref.Ref{
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonChartModule),
		resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonChartModule),
	}, c.ResourceRefs())
}

func TestPageResourceRefs(t *testing.T) {
	p := Page{
		ModuleID: 1001,
		Blocks: PageBlocks{
			{Kind: "RecordList", Options: map[string]interface{}{"moduleID": "1002"}},
			{Kind: "Chart", Options: map[string]interface{}{"chartID": "3001"}},
			{Kind: "Calendar", Options: map[string]interface{}{
				"feeds": []interface{}{
					map[string]interface{}{"options": map[string]interface{}{"moduleID": "1003"}},
				},
			}},
			{Kind: "Metric", Options: map[string]interface{}{
				"metrics": []interface{}{
					map[string]interface{}{"moduleID": "1004"},
				},
			}},
			{Kind: "Progress", Options: map[string]interface{}{
				"value": map[string]interface{}{"moduleID": "1005"},
			}},
			{Kind: "Automation", Options: map[string]interface{}{
				"buttons": []interface{}{
					map[string]interface{}{"workflowID": "5001"},
				},
			}},
			{Kind: "Content", Options: map[string]interface{}{}},
		},
	}

	require.Equal(t, []resourceref.Ref{
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonPageModule),
		resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonPageModule),
		resourceref.Make(resourceref.KindComposeChart, 3001, resourceref.ReasonPageChart),
		resourceref.Make(resourceref.KindComposeModule, 1003, resourceref.ReasonPageModule),
		resourceref.Make(resourceref.KindComposeModule, 1004, resourceref.ReasonPageModule),
		resourceref.Make(resourceref.KindComposeModule, 1005, resourceref.ReasonPageModule),
		resourceref.Make(resourceref.KindAutomationWorkflow, 5001, resourceref.ReasonPageWorkflow),
	}, p.ResourceRefs())
}

func TestPageBlockResourceRefs(t *testing.T) {
	cc := []struct {
		name string
		b    PageBlock
		out  []resourceref.Ref
	}{
		{
			"RecordList selection buttons",
			PageBlock{Kind: "RecordList", Options: map[string]interface{}{
				"moduleID": "1002",
				"selectionButtons": []interface{}{
					map[string]interface{}{"workflowID": "5001"},
					map[string]interface{}{"automationID": "6001"},
				},
			}},
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonPageModule),
				resourceref.Make(resourceref.KindAutomationWorkflow, 5001, resourceref.ReasonPageWorkflow),
				resourceref.Make(resourceref.KindNgAutomation, 6001, resourceref.ReasonPageAutomation),
			},
		},
		{
			"Automation button workflow + ng-automation",
			PageBlock{Kind: "Automation", Options: map[string]interface{}{
				"buttons": []interface{}{
					map[string]interface{}{"workflowID": "5001", "automationID": "6001"},
				},
			}},
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindAutomationWorkflow, 5001, resourceref.ReasonPageWorkflow),
				resourceref.Make(resourceref.KindNgAutomation, 6001, resourceref.ReasonPageAutomation),
			},
		},
		{
			"AgentChat default + allowed",
			PageBlock{Kind: "AgentChat", Options: map[string]interface{}{
				"defaultAgentID":  "7001",
				"allowedAgentIDs": []interface{}{"7002", "7003"},
			}},
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindAgent, 7001, resourceref.ReasonPageAgent),
				resourceref.Make(resourceref.KindAgent, 7002, resourceref.ReasonPageAgent),
				resourceref.Make(resourceref.KindAgent, 7003, resourceref.ReasonPageAgent),
			},
		},
		{
			"ChatbotInbox",
			PageBlock{Kind: "ChatbotInbox", Options: map[string]interface{}{
				"chatbotIDs": []interface{}{"8001", "8002"},
			}},
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindChatbot, 8001, resourceref.ReasonPageChatbot),
				resourceref.Make(resourceref.KindChatbot, 8002, resourceref.ReasonPageChatbot),
			},
		},
		{
			"Geometry feeds",
			PageBlock{Kind: "Geometry", Options: map[string]interface{}{
				"feeds": []interface{}{
					map[string]interface{}{"options": map[string]interface{}{"moduleID": "1009"}},
				},
			}},
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindComposeModule, 1009, resourceref.ReasonPageModule),
			},
		},
		{
			"Navigation items",
			PageBlock{Kind: "Navigation", Options: map[string]interface{}{
				"navigationItems": []interface{}{
					map[string]interface{}{"options": map[string]interface{}{"item": map[string]interface{}{
						"pageID":       "2001",
						"pageLayoutID": "2002",
						"moduleID":     "1010",
					}}},
				},
			}},
			[]resourceref.Ref{
				resourceref.Make(resourceref.KindComposePage, 2001, resourceref.ReasonPageNavigation),
				resourceref.Make(resourceref.KindComposePageLayout, 2002, resourceref.ReasonPageNavigation),
				resourceref.Make(resourceref.KindComposeModule, 1010, resourceref.ReasonPageModule),
			},
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, PageBlocks{c.b}.ResourceRefs())
		})
	}
}
