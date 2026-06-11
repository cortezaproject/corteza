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
		{Kind: resourceref.KindDalConnection, ID: 42, Reason: resourceref.ReasonModuleConnection, Path: "Config.DAL.ConnectionID"},
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonModuleFieldRef, Path: "Fields.lead.Options.ModuleID"},
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
			[]resourceref.Ref{{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonModuleFieldRef, Path: "Options.ModuleID"}},
		},
		{
			"handle ref",
			ModuleFieldOptions{"moduleID": "lead-module"},
			[]resourceref.Ref{{Kind: resourceref.KindComposeModule, Ident: "lead-module", Reason: resourceref.ReasonModuleFieldRef, Path: "Options.ModuleID"}},
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
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonChartModule, Path: "Config.Reports.0.ModuleID"},
		{Kind: resourceref.KindComposeModule, ID: 1002, Reason: resourceref.ReasonChartModule, Path: "Config.Reports.3.ModuleID"},
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
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonPageModule, Path: "ModuleID"},
		{Kind: resourceref.KindComposeModule, ID: 1002, Reason: resourceref.ReasonPageModule, Path: "Blocks.0.Options.ModuleID"},
		{Kind: resourceref.KindComposeChart, ID: 3001, Reason: resourceref.ReasonPageChart, Path: "Blocks.1.Options.ChartID"},
		{Kind: resourceref.KindComposeModule, ID: 1003, Reason: resourceref.ReasonPageModule, Path: "Blocks.2.Options.feeds.0.ModuleID"},
		{Kind: resourceref.KindComposeModule, ID: 1004, Reason: resourceref.ReasonPageModule, Path: "Blocks.3.Options.metrics.0.ModuleID"},
		{Kind: resourceref.KindComposeModule, ID: 1005, Reason: resourceref.ReasonPageModule, Path: "Blocks.4.Options.value.ModuleID"},
		{Kind: resourceref.KindAutomationWorkflow, ID: 5001, Reason: resourceref.ReasonPageWorkflow, Path: "Blocks.5.Options.buttons.0.WorkflowID"},
	}, p.ResourceRefs())
}
