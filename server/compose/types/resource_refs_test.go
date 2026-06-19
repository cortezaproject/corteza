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
