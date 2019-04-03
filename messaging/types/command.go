package types

type (
	Command struct {
		Name        string          `db:"name"        json:"name"`
		Params      CommandParamSet `db:"params"      json:"params"`
		Description string          `db:"description" json:"description"`
	}
)

var (
	Preset CommandSet // @todo move this to someplace safe
)

func init() {
	Preset = CommandSet{
		&Command{
			Name:        "echo",
			Description: "It does exactly what it says on the tin"},
		&Command{
			Name:        "shrug",
			Description: "It does exactly what it says on the tin"},
	}
}
