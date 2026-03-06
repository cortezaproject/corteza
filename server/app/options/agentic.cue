package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

agentic: schema.#optionsGroup & {
	handle: "agentic"
	title:  "Agentic"

	options: {
		mcp_server_name: {
			description:  "MCP server name reported to clients."
			defaultValue: "Corteza MCP"
		}
		mcp_server_version: {
			description:  "MCP server version reported to clients."
			defaultValue: "v1"
		}
		anthropic_api_version: {
			description:  "Anthropic API version header sent with every request."
			defaultValue: "2023-06-01"
		}
	}
}
