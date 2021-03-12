package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/workflow.yaml

type (
	WorkflowOpt struct {
		Enabled bool `env:"WORKFLOW_ENABLED"`
	}
)

// Workflow initializes and returns a WorkflowOpt with default values
func Workflow() (o *WorkflowOpt) {
	o = &WorkflowOpt{
		Enabled: true,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Workflow) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
