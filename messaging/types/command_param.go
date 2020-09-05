package types

type (
	CommandParam struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Required bool   `json:"required"`
	}
)
