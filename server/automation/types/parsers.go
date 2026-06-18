package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/expr"
)

func ParseWorkflowVariables(ss []string) (p *expr.Vars, err error) {
	p = &expr.Vars{}
	return p, parseStringsInput(ss, &p)
}

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
}
