package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

workflow: schema.#optionsGroup & {
	handle: "workflow"
	options: {
		register: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Registers enabled and valid workflows and executes them when triggered"
		}
		exec_debug: {
			type:        "bool"
			description: "Enables verbose logging for workflow execution"
		}
		call_stack_size: {
			type:          "int"
			defaultGoExpr: "16"
			description:   "Defines the maximum call stack size between workflows"
		}
	}
	title: "Workflow"
}
