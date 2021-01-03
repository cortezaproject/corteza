package fn

import "github.com/cortezaproject/corteza-server/pkg/wfexec"

type (
	// workflow functions are defined in the core code and through plugins
	Function struct {
		Ref        string       `json:"ref"`
		Meta       FunctionMeta `json:"meta"`
		Parameters []*Param     `json:"parameters"`
		Results    []*Param     `json:"results"`

		Handler wfexec.ActivityHandler `json:"-"`
	}

	FunctionMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}
)
