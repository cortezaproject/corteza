package gig

type (
	taskDefParam struct {
		Name        string `json:"name"`
		Kind        string `json:"kind"`
		Required    bool   `json:"required"`
		Description string `json:"description"`
	}
	TaskDef struct {
		Ref         string         `json:"ref"`
		Kind        string         `json:"kind"`
		Description string         `json:"description,omitempty"`
		Params      []taskDefParam `json:"params,omitempty"`
	}
	TaskDefSet []TaskDef
)

var (
	TaskDecoder       = "decoder"
	TaskPreprocessor  = "preprocessor"
	TaskPostprocessor = "postprocessor"
)
