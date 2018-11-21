package types

//go:generate go run ../../codegen/v2/type-set.go --no-pk-types Command,CommandParam --output command.gen.go

type (
	Command struct {
		Name        string          `db:"name"        json:"name"`
		Params      CommandParamSet `db:"params"      json:"params"`
		Description string          `db:"description" json:"description"`
	}

	CommandParam struct {
		Name     string `db:"name"     json:"name"`
		Type     string `db:"type"     json:"type"`
		Required bool   `db:"required" json:"required"`
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
