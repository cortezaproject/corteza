package types

type (
	// Fields - CRM input field definitions
	Field struct {
		Name string `json:"name" db:"name"`
		Type string `json:"type" db:"type"`
	}
)
