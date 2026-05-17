package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/system/agentic/skills"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id.Init(cli.Context())
}

// mocks

type mockRegistry struct {
	agent *types.Agent
	err   error
}

func (m *mockRegistry) Get(_ context.Context, _ uint64) (*types.Agent, error) {
	return m.agent, m.err
}

type mockLLM struct {
	responses        []LLMResponse
	callErrors       []error
	callIdx          int
	captured         [][]types.AiConversationMessage
	capturedPrompts  []string
}

func (m *mockLLM) Chat(_ context.Context, prompt string, msgs []types.AiConversationMessage, _ []Tool, _ LLMConfig) (*LLMResponse, error) {
	m.captured = append(m.captured, msgs)
	m.capturedPrompts = append(m.capturedPrompts, prompt)
	i := m.callIdx
	m.callIdx++
	if i < len(m.callErrors) && m.callErrors[i] != nil {
		return nil, m.callErrors[i]
	}
	if i < len(m.responses) {
		r := m.responses[i]
		return &r, nil
	}
	return &LLMResponse{Text: "done"}, nil
}

type mockSkillRegistry struct {
	skills map[string][]*skills.Skill
}

func (m *mockSkillRegistry) ForTool(name string) []*skills.Skill { return m.skills[name] }
func (m *mockSkillRegistry) All() []*skills.Skill {
	var out []*skills.Skill
	for _, ss := range m.skills {
		out = append(out, ss...)
	}
	return out
}

type mockMCP struct {
	tools    []Tool
	toolsErr error
	result   any
	execErr  error
}

func (m *mockMCP) GetTools(_ context.Context, _ []string) ([]Tool, error) {
	return m.tools, m.toolsErr
}

func (m *mockMCP) ExecuteTool(_ context.Context, _ string, _ map[string]any) (any, error) {
	return m.result, m.execErr
}

type mockStore struct{}

func (m *mockStore) FindByID(_ context.Context, _ uint64) (*types.AiConversation, error) {
	return nil, errors.New("not found")
}

func (m *mockStore) Create(_ context.Context, c *types.AiConversation) (*types.AiConversation, error) {
	c.ID = 1
	return c, nil
}

func (m *mockStore) Update(_ context.Context, c *types.AiConversation) (*types.AiConversation, error) {
	return c, nil
}

func (m *mockStore) DeleteByID(_ context.Context, _ uint64) error {
	return nil
}

// helpers 

func testCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), auth.Authenticated(1))
}

func activeAgent() *types.Agent {
	return &types.Agent{
		ID:     42,
		Status: "active",
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{Name: "test_tool"}},
		},
	}
}

// errorCode extracts meta.code from a pkg/errors Error.
func errorCode(err error) string {
	type jsonMarshaler interface {
		MarshalJSON() ([]byte, error)
	}
	jm, ok := err.(jsonMarshaler)
	if !ok {
		return ""
	}
	data, _ := jm.MarshalJSON()
	var v struct {
		Meta struct {
			Code string `json:"code"`
		} `json:"meta"`
	}
	json.Unmarshal(data, &v)
	return v.Meta.Code
}

func newRuntime(reg Registry, llm LLMClient, mcp MCPClient) *runtime {
	return Runtime(reg, llm, mcp, &mockStore{}, nil, nil, nil, nil)
}

//  tests 

func TestRun_AgentNotFound(t *testing.T) {
	rt := newRuntime(&mockRegistry{err: errors.New("not found")}, nil, nil)
	_, err := rt.Run(testCtx(), &AgentRequest{AgentID: 99, Input: "hi"})
	require.Error(t, err)
	assert.Equal(t, "agent_not_found", errorCode(err))
}

func TestRun_AgentDisabled(t *testing.T) {
	agent := activeAgent()
	agent.Status = "disabled"
	rt := newRuntime(&mockRegistry{agent: agent}, nil, nil)
	_, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.Error(t, err)
	assert.Equal(t, "agent_disabled", errorCode(err))
}

func TestRun_MCPError(t *testing.T) {
	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		nil,
		&mockMCP{toolsErr: errors.New("mcp failure")},
	)
	_, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.Error(t, err)
	assert.Equal(t, "mcp_error", errorCode(err))
}

func TestRun_LLMError(t *testing.T) {
	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		&mockLLM{callErrors: []error{errors.New("provider unavailable")}},
		&mockMCP{},
	)
	_, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.Error(t, err)
	assert.Equal(t, "llm_error", errorCode(err))
}

func TestRun_TokenLimitExceeded(t *testing.T) {
	agent := activeAgent()
	agent.Execution.Limits.ContextWindow = 10

	rt := newRuntime(
		&mockRegistry{agent: agent},
		&mockLLM{responses: []LLMResponse{
			{
				Text:      "thinking",
				ToolCalls: []ToolCall{{ID: "1", Name: "test_tool", Args: map[string]any{}}},
				Usage:     Usage{InputTokens: 6, OutputTokens: 6, ContextWindow: 12},
			},
		}},
		&mockMCP{result: "ok"},
	)
	_, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.Error(t, err)
	assert.Equal(t, "limit_exceeded", errorCode(err))
}

func TestRun_Timeout(t *testing.T) {
	ctx, cancel := context.WithCancel(testCtx())
	cancel() // immediately cancelled

	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		&mockLLM{},
		&mockMCP{},
	)
	_, err := rt.Run(ctx, &AgentRequest{AgentID: 42, Input: "hi"})
	require.Error(t, err)
	assert.Equal(t, "timeout", errorCode(err))
}

func TestRun_ToolSoftError(t *testing.T) {
	llm := &mockLLM{responses: []LLMResponse{
		{Text: "using tool", ToolCalls: []ToolCall{{ID: "1", Name: "test_tool", Args: map[string]any{}}}},
		{Text: "final answer"},
	}}
	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		llm,
		&mockMCP{execErr: errors.New("tool broke")},
	)
	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.NoError(t, err)
	assert.Equal(t, "final answer", resp.Output)
	require.Len(t, resp.ToolCalls, 1)
	assert.Equal(t, "tool broke", resp.ToolCalls[0].Error)
}

func TestRun_WindDownInjected(t *testing.T) {
	agent := activeAgent()
	agent.Execution.Limits.ContextWindow = 100
	agent.Execution.Limits.SoftLimitRatio = 0.5

	llm := &mockLLM{responses: []LLMResponse{
		{
			Text:      "using tool",
			ToolCalls: []ToolCall{{ID: "1", Name: "test_tool", Args: map[string]any{}}},
			Usage:     Usage{InputTokens: 30, OutputTokens: 30, ContextWindow: 60},
		},
		{Text: "final answer"},
	}}
	rt := newRuntime(
		&mockRegistry{agent: agent},
		llm,
		&mockMCP{result: "ok"},
	)
	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.NoError(t, err)
	assert.Equal(t, "final answer", resp.Output)

	// second LLM call should have the wind-down message in its history
	require.Len(t, llm.captured, 2)
	var found bool
	for _, m := range llm.captured[1] {
		if m.Content == "You are reaching the maximum token limit. Please finish up and give your final answer." {
			found = true
		}
	}
	assert.True(t, found, "expected wind-down message in conversation history before second LLM call")
}

func TestRun_SkillActivatedAfterToolCall(t *testing.T) {
	llm := &mockLLM{responses: []LLMResponse{
		{Text: "using tool", ToolCalls: []ToolCall{{ID: "1", Name: "test_tool", Args: map[string]any{}}}},
		{Text: "final answer"},
	}}
	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		llm,
		&mockMCP{result: "ok"},
	)
	rt.SetSkillRegistry(&mockSkillRegistry{skills: map[string][]*skills.Skill{
		"test_tool": {{
			Name:     "demo_skill",
			Triggers: []string{"test_tool"},
			Body:     "DEMO_SKILL_BODY_MARKER",
		}},
	}})

	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.NoError(t, err)
	assert.Equal(t, "final answer", resp.Output)

	require.Len(t, llm.capturedPrompts, 2)
	assert.NotContains(t, llm.capturedPrompts[0], "DEMO_SKILL_BODY_MARKER", "skill must not appear in the first prompt — no tool called yet")
	assert.Contains(t, llm.capturedPrompts[1], "DEMO_SKILL_BODY_MARKER", "skill must appear in the second prompt after the tool was called")
	assert.Contains(t, llm.capturedPrompts[1], "## SKILL: demo_skill", "skill header must be present")
}

func TestRun_SkillInjectedOnlyOnce(t *testing.T) {
	llm := &mockLLM{responses: []LLMResponse{
		{Text: "tool 1", ToolCalls: []ToolCall{{ID: "1", Name: "test_tool", Args: map[string]any{}}}},
		{Text: "tool 2", ToolCalls: []ToolCall{{ID: "2", Name: "test_tool", Args: map[string]any{}}}},
		{Text: "final answer"},
	}}
	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		llm,
		&mockMCP{result: "ok"},
	)
	rt.SetSkillRegistry(&mockSkillRegistry{skills: map[string][]*skills.Skill{
		"test_tool": {{
			Name:     "demo_skill",
			Triggers: []string{"test_tool"},
			Body:     "DEMO_SKILL_BODY_MARKER",
		}},
	}})

	_, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.NoError(t, err)

	require.Len(t, llm.capturedPrompts, 3)
	count := strings.Count(llm.capturedPrompts[2], "DEMO_SKILL_BODY_MARKER")
	assert.Equal(t, 1, count, "skill body should appear exactly once even after multiple matching tool calls")
}

func TestRun_HappyPath(t *testing.T) {
	rt := newRuntime(
		&mockRegistry{agent: activeAgent()},
		&mockLLM{responses: []LLMResponse{{Text: "Hello, how can I help?"}}},
		&mockMCP{},
	)
	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.NoError(t, err)
	assert.Equal(t, "Hello, how can I help?", resp.Output)
	assert.NotZero(t, resp.ConversationID)
}
