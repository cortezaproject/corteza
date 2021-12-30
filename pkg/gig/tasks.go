package gig

type (
	task string

	taskDefParam struct {
		Name        string `json:"name"`
		Kind        string `json:"kind"`
		Required    bool   `json:"required"`
		Description string `json:"description"`
	}
	TaskDef struct {
		Ref    string         `json:"ref"`
		Kind   task           `json:"kind"`
		Params []taskDefParam `json:"params"`
	}
	TaskDefSet []TaskDef
)

var (
	TaskDecoder       task = "decoder"
	TaskPreprocessor  task = "preprocessor"
	TaskPostprocessor task = "postprocessor"
)
