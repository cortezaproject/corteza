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
	responses       []LLMResponse
	callErrors      []error
	callIdx         int
	captured        [][]types.AiConversationMessage
	capturedPrompts []string
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
	// inGroup answers ToolNamesIn, keyed "group/maxRisk"; calls records what
	// was asked, so a test can assert the risk ceiling was passed on.
	inGroup map[string][]string
	calls   []string
}

func (m *mockMCP) GetTools(_ context.Context, _ []string) ([]Tool, error) {
	return m.tools, m.toolsErr
}

func (m *mockMCP) ExecuteTool(_ context.Context, _ string, _ map[string]any) (any, error) {
	return m.result, m.execErr
}

func (m *mockMCP) ToolNamesIn(group, maxRisk string) []string {
	key := group + "/" + maxRisk
	m.calls = append(m.calls, key)
	return m.inGroup[key]
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

// activeAgent runs its tool without asking. Most of these tests are about what
// happens once a tool has run; an unstated mode resolves to "ask" for anything
// the registry does not classify as read-only, and the mock classifies nothing.
func activeAgent() *types.Agent {
	return &types.Agent{
		ID:     42,
		Status: "active",
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{Name: "test_tool", Permission: "always"}},
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

// askingAgent grants one tool the operator wants to be asked about.
func askingAgent() *types.Agent {
	a := activeAgent()
	a.Access.Tools = []types.AgentAccessTool{{Name: "test_tool", Permission: "ask"}}
	return a
}

func toolCallingLLM() *mockLLM {
	return &mockLLM{responses: []LLMResponse{
		{Text: "using tool", ToolCalls: []ToolCall{{ID: "1", Name: "test_tool", Args: map[string]any{}}}},
		{Text: "final answer"},
	}}
}

// A tool set to "ask" stops the run rather than refusing it: the conversation
// is saved as it stands and resumes once the caller comes back with it
// approved.
func TestRun_StopsForApproval(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: askingAgent()}, toolCallingLLM(), &mockMCP{})

	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "do it"})
	require.NoError(t, err)
	require.Equal(t, StatusAwaitingApproval, resp.Status)
	require.NotNil(t, resp.PendingApproval)
	assert.Equal(t, "test_tool", resp.PendingApproval.Tool)

	// The tool must not have run.
	require.Len(t, resp.ToolCalls, 1)
	assert.Contains(t, resp.ToolCalls[0].Error, "waiting for the user to approve")
}

// The same run with the tool approved goes through and never asks.
func TestRun_ApprovedToolRuns(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: askingAgent()}, toolCallingLLM(), &mockMCP{})

	resp, err := rt.Run(testCtx(), &AgentRequest{
		AgentID:       42,
		Input:         "do it",
		ApprovedTools: []string{"test_tool"},
	})
	require.NoError(t, err)
	assert.Equal(t, StatusComplete, resp.Status)
	assert.Nil(t, resp.PendingApproval)
	assert.Equal(t, "final answer", resp.Output)
}

// A resume carries the approval and nothing to say. The turn still needs a user
// message: the history it resumes ends in tool results saying the work was not
// done, and something has to tell the agent it may now do it.
func TestRun_ResumeWithoutInput(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: askingAgent()}, toolCallingLLM(), &mockMCP{})

	resp, err := rt.Run(testCtx(), &AgentRequest{
		AgentID:       42,
		ApprovedTools: []string{"test_tool"},
	})
	require.NoError(t, err)
	assert.Equal(t, StatusComplete, resp.Status)
}

// A run that never stops still says so, so a caller can branch on one field.
func TestRun_ReportsComplete(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: activeAgent()}, toolCallingLLM(), &mockMCP{})

	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "hi"})
	require.NoError(t, err)
	assert.Equal(t, StatusComplete, resp.Status)
}

// A chatbot widget or an automation-triggered run has nobody to put an approval
// to. Stalling there produces no output and no error, so "ask" is a refusal
// that says which setting would let it run.
func TestRun_UnattendedRefusesInsteadOfAsking(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: askingAgent()}, toolCallingLLM(), &mockMCP{})

	resp, err := rt.Run(testCtx(), &AgentRequest{AgentID: 42, Input: "do it", Unattended: true})
	require.NoError(t, err)
	assert.Equal(t, StatusComplete, resp.Status, "an unattended run never pauses")
	assert.Nil(t, resp.PendingApproval)

	require.Len(t, resp.ToolCalls, 1)
	assert.Contains(t, resp.ToolCalls[0].Error, "nobody to ask")
	assert.Contains(t, resp.ToolCalls[0].Error, "always")
}

// An unattended run with the tool approved up front still goes through, so an
// automation can carry an approval its operator set.
func TestRun_UnattendedRunsAnApprovedTool(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: askingAgent()}, toolCallingLLM(), &mockMCP{})

	resp, err := rt.Run(testCtx(), &AgentRequest{
		AgentID:       42,
		Input:         "do it",
		Unattended:    true,
		ApprovedTools: []string{"test_tool"},
	})
	require.NoError(t, err)
	assert.Equal(t, StatusComplete, resp.Status)
	assert.Empty(t, resp.ToolCalls[0].Error)
}

// A per-TAQ tool is minted as "automation_<id>", which names nothing anyone
// recognises and is in no group. Both had to be answered for separately, or the
// prompt reads "may I use automation_510775573053898753 (deletes)".
func TestPendingApprovalDescribesAMintedAutomation(t *testing.T) {
	rt := newRuntime(&mockRegistry{agent: activeAgent()}, nil, &mockMCP{})

	assert.Equal(t, "write", rt.toolRisk("automation_510775573053898753"),
		"running an automation writes; it is not a deletion")
	assert.Equal(t, "destructive", rt.toolRisk("compose_record_delete"),
		"a name the write ceiling does not cover deletes")

	titles := toolTitles([]Tool{
		{Name: "automation_510775573053898753", Title: "Refresh every holding's totals"},
		{Name: "compose_record_create"},
	})
	assert.Equal(t, "Refresh every holding's totals", titles["automation_510775573053898753"])
	assert.NotContains(t, titles, "compose_record_create", "no title is better than an empty one")
}

func TestDynamicAutomationRef(t *testing.T) {
	for name, want := range map[string]bool{
		"automation_123":           true,
		"automation_taq_exec":      false,
		"automation_workflow_exec": false,
		"compose_record_create":    false,
		"automation_":              false,
	} {
		_, ok := dynamicAutomationRef(name)
		assert.Equal(t, want, ok, "for %q", name)
	}
}
