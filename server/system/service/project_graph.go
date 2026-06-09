package service

import (
	"context"

	"github.com/crusttech/human/server/system/types"
)

type projectGraphService struct{}

var DefaultProjectGraph = &projectGraphService{}

// Graph returns a mock project graph until resource scoping (project_id columns) is in place.
func (s *projectGraphService) Graph(_ context.Context, _ uint64) (*types.ProjectGraph, error) {
	nodes := []*types.ProjectGraphNode{
		{ID: 1001, Kind: "module", Name: "Lead", Sensitivity: "internal"},
		{ID: 1002, Kind: "module", Name: "Opportunity"},
		{ID: 1003, Kind: "module", Name: "Quote"},
		{ID: 2001, Kind: "page", Name: "Lead List"},
		{ID: 2002, Kind: "page", Name: "Opportunity Board"},
		{ID: 3001, Kind: "chart", Name: "Pipeline Forecast"},
		{ID: 4001, Kind: "connection", Name: "Google Sheets", Sensitivity: "confidential"},
		{ID: 4002, Kind: "connection", Name: "Stripe", Sensitivity: "restricted"},
		{ID: 5001, Kind: "automation", Name: "Lead Scoring"},
		{ID: 6001, Kind: "agent", Name: "Sales Assistant"},
		{ID: 7001, Kind: "chatbot", Name: "Support Bot"},
		{ID: 8001, Kind: "role", Name: "Sales Rep"},
		{ID: 8002, Kind: "role", Name: "Sales Manager"},
		{ID: 8003, Kind: "role", Name: "Customer"},
	}

	edges := []*types.ProjectGraphEdge{
		// module-field-ref: Opportunity.Lead → Lead
		{SourceID: 1002, TargetID: 1001, Reason: "module-field-ref"},
		// module-field-ref: Quote.Opportunity → Opportunity
		{SourceID: 1003, TargetID: 1002, Reason: "module-field-ref"},
		// page-module
		{SourceID: 2001, TargetID: 1001, Reason: "page-module"},
		{SourceID: 2002, TargetID: 1002, Reason: "page-module"},
		// chart reads module
		{SourceID: 3001, TargetID: 1002, Reason: "page-module"},
		// trigger-module: Lead Scoring → Lead
		{SourceID: 5001, TargetID: 1001, Reason: "trigger-module"},
		// workflow-agent: Sales Assistant reads Opportunity
		{SourceID: 6001, TargetID: 1002, Reason: "workflow-agent"},
		// chatbot-agent: Support Bot reads Lead
		{SourceID: 7001, TargetID: 1001, Reason: "chatbot-agent"},
	}

	return &types.ProjectGraph{Nodes: nodes, Edges: edges}, nil
}
