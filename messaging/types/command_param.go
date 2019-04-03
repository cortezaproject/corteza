package types

type (
	CommandParam struct {
		Name     string `db:"name"     json:"name"`
		Type     string `db:"type"     json:"type"`
		Required bool   `db:"required" json:"required"`
	}
)
