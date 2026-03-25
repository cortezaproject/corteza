package runtime

import "github.com/cortezaproject/corteza/server/pkg/errors"

func errAgentNotFound(id uint64) error {
	return errors.New(errors.KindNotFound, "agent not found",
		errors.Meta("code", "agent_not_found"),
		errors.Meta("agentID", id),
	)
}

func errAgentDisabled(id uint64) error {
	return errors.New(errors.KindUnauthenticated, "agent is disabled",
		errors.Meta("code", "agent_disabled"),
		errors.Meta("agentID", id),
	)
}

func errRBACDenied() error {
	return errors.New(errors.KindUnauthenticated, "access denied",
		errors.Meta("code", "rbac_denied"),
	)
}

func errConversationNotFound(id uint64) error {
	return errors.New(errors.KindNotFound, "conversation not found",
		errors.Meta("code", "conversation_not_found"),
		errors.Meta("conversationID", id),
	)
}


func errLLM(err error) error {
	return errors.New(errors.KindExternal, "LLM request failed: "+err.Error(),
		errors.Meta("code", "llm_error"),
		errors.Wrap(err),
	)
}

func errMCP(err error) error {
	return errors.New(errors.KindExternal, "MCP error: "+err.Error(),
		errors.Meta("code", "mcp_error"),
		errors.Wrap(err),
	)
}

func errLimitExceeded(reason string) error {
	return errors.New(errors.KindInvalidData, "limit exceeded: "+reason,
		errors.Meta("code", "limit_exceeded"),
	)
}

func errTimeout() error {
	return errors.New(errors.KindExternal, "agent execution timed out",
		errors.Meta("code", "timeout"),
	)
}
