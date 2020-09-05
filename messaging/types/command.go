package types

type (
	Command struct {
		Name        string          `json:"name"`
		Params      CommandParamSet `json:"params"`
		Description string          `json:"description"`
	}
)
