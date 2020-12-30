package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	// workflow functions are defined in the core code and through plugins
	WorkflowFunction struct {
		Ref        string                 `json:"ref"`
		Meta       WorkflowFunctionMeta   `json:"meta"`
		Handler    wfexec.ActivityHandler `json:"-"`
		Parameters []*WorkflowParameter   `json:"parameters"`
		Results    []*WorkflowParameter   `json:"results"`
	}

	WorkflowFunctionMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowParameter struct {
		Name string                `json:"name"`
		Type string                `json:"type"`
		Meta WorkflowParameterMeta `json:"meta"`
	}

	WorkflowParameterMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// Used for expression steps, arguments and results mapping
	WorkflowExpression struct {
		Name string `json:"name"`
		Expr string `json:"expr"`
	}

	Variables map[string]interface{}
)

func ParseWorkflowVariables(ss []string) (p Variables, err error) {
	p = Variables{}
	return p, parseStringsInput(ss, &p)
}

func (vv *Variables) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = Variables{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into Variables: %w", string(b), err)
		}
	}

	return nil
}

func (vv Variables) Value() (driver.Value, error) {
	return json.Marshal(vv)
}
