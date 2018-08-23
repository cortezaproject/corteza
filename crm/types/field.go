package types

/* If you want to edit this file by hand, remove codegen/[project]/types/index.php */

type (
	// Fields - CRM input field definitions
	Field struct {
		Name     string `json:"field_name" db:"field_name"`
		Type     string `json:"field_type" db:"field_type"`
		Template string `json:"field_template,omitempty" db:"field_template"`
	}
)
