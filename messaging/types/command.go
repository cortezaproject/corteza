package types

type (
	Command struct {
		Name        string          `db:"name"        json:"name"`
		Params      CommandParamSet `db:"params"      json:"params"`
		Description string          `db:"description" json:"description"`
	}
)
