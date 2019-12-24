package corredor

type (
	Script struct {
		Name        string     `json:"name"`
		Label       string     `json:"label"`
		Description string     `json:"description"`
		Errors      []string   `json:"errors,omitempty"`
		Triggers    []*Trigger `json:"triggers"`

		// If bundle or type is set, consider
		// this a frontend script
		Bundle string `json:"bundle,omitempty"`
		Type   string `json:"type,omitempty"`
	}
)
