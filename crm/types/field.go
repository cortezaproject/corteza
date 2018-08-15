package types

/* If you want to edit this file by hand, remove codegen/[project]/types/index.php */

type (
	// Fields - CRM input field definitions
	Field struct {
		Name string `json:"name" db:"name"`
		Type string `json:"type" db:"type"`
	}
)
