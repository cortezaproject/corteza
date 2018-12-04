package rbac

type (
	OperationGroup struct {
		Title      string      `json:"title"`
		Operations []Operation `json:"operations"`
	}

	Operation struct {
		Key      string `json:"key"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`

		// Enabled will show/hide the operation in administration
		Enabled bool `json:"enabled"`

		// The default value for managing operations on a role
		//
		// nil = unset (inherit),
		// true = checked (allow),
		// false = unchecked (deny)

		Default OperationState `json:"default"`
	}

	OperationState string

	Permission struct {
		// Resource (organisation, team, channel)
		Resource string `json:"resource"`
		// Operation name (Operation.Key)
		Operation string `json:"operation"`
		// Operation state (inherit, allow, deny)
		Value OperationState `json:"value"`
	}
)

const (
	OperationStateInherit OperationState = ""
	OperationStateAllow                  = "allow"
	OperationStateDeny                   = "deny"
)
