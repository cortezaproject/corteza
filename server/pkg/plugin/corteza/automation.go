package corteza

import "github.com/cortezaproject/corteza/server/automation/types"

type (
	CortezaPluginAutomationFunction interface {
		AutomationFunctions() []*types.Function
	}
)
