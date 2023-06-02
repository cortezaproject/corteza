package corteza

// CortezaPlugin should know how to:
//   - fetch meta for a specific automation function
//   - execute an automation function via parameters
//
// Do we create a wrapper for that?
type (
	CortezaPlugin interface {
		// we already have an AutomationFunction plugin
		// we need support for multiple of them
		// and also support for other types of handlers
		// Meta() *types.Function
		// Exec(in *expr.Vars) (out *expr.Vars, err error)

		// fetch all automation functions
		CortezaPluginAutomationFunction
		CortezaPluginExpressionType
	}
)
