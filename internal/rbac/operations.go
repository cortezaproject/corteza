package rbac

type (
	OperationGroup struct {
		Name       string
		Operations []Operation
	}

	Operation struct {
		Key   string
		Name  string
		Title string

		// Enabled will show/hide the operation in administration
		Enabled bool

		// The default value for managing operations on a role
		//
		// nil = unset (inherit),
		// true = checked (allow),
		// false = unchecked (deny)

		Default OperationState
	}

	OperationState string
)

const (
	OperationStateInherit OperationState = ""
	OperationStateAllow                  = "allow"
	OperationStateDeny                   = "deny"
)
