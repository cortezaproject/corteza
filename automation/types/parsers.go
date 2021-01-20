package types

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

func ParseWorkflowVariables(ss []string) (p *expr.Vars, err error) {
	p = &expr.Vars{}
	return p, parseStringsInput(ss, &p)
}

func ParseWorkflowMeta(ss []string) (p *WorkflowMeta, err error) {
	p = &WorkflowMeta{}
	return p, parseStringsInput(ss, p)
}

func ParseWorkflowStepSet(ss []string) (p WorkflowStepSet, err error) {
	p = WorkflowStepSet{}
	return p, parseStringsInput(ss, &p)
}

func ParseWorkflowPathSet(ss []string) (p WorkflowPathSet, err error) {
	p = WorkflowPathSet{}
	return p, parseStringsInput(ss, &p)
}

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
}
