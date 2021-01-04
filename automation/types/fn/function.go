package fn

import "github.com/cortezaproject/corteza-server/pkg/wfexec"

type (
	// workflow functions are defined in the core code and through plugins
	Function struct {
		Ref        string        `json:"ref,omitempty"`
		Meta       *FunctionMeta `json:"meta,omitempty"`
		Parameters []*Param      `json:"parameters,omitempty"`
		Results    []*Param      `json:"results,omitempty"`

		Handler wfexec.ActivityHandler `json:"-"`
	}

	FunctionMeta struct {
		Label       string                 `json:"label,omitempty"`
		Description string                 `json:"description,omitempty"`
		Visual      map[string]interface{} `json:"visual,omitempty"`
	}
)
