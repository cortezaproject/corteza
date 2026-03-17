package automation

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	agenticRuntime "github.com/cortezaproject/corteza/server/system/agentic/runtime"
	"github.com/cortezaproject/corteza/server/system/types"
)

type mockConversationStore struct {
	conv *types.AiConversation
	err  error
}

func (m *mockConversationStore) FindByID(_ context.Context, _ uint64) (*types.AiConversation, error) {
	return m.conv, m.err
}

func newNgHandler(rt ngAgentRuntime, convs ngAgentConversationStore) *ngAgentHandler {
	return &ngAgentHandler{
		runtime:       rt,
		conversations: convs,
	}
}

func TestNgAgentHandler_Prompt(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := newNgHandler(rt, &mockConversationStore{})

	res, err := h.prompt(context.Background(), &ngAgentPromptArgs{AgentID: 1, Input: "hello"})
	req.NoError(err)
	req.Equal("ok", res.Output)
	req.Equal(uint64(1), res.ConversationID)
	req.True(rt.called)
}

func TestNgAgentHandler_PromptRuntimeError(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntimeErr{}
	h := newNgHandler(rt, &mockConversationStore{})

	_, err := h.prompt(context.Background(), &ngAgentPromptArgs{AgentID: 1, Input: "hello"})
	req.Error(err)
}

func TestNgAgentHandler_ContinueConversation(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := newNgHandler(rt, &mockConversationStore{
		conv: &types.AiConversation{AgentID: 5},
	})

	res, err := h.continueConversation(context.Background(), &ngAgentContinueArgs{ConversationID: 10, Input: "follow up"})
	req.NoError(err)
	req.Equal("ok", res.Output)
	req.Equal(uint64(1), res.ConversationID)
	req.True(rt.called)
}

func TestNgAgentHandler_ContinueConversationNotFound(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := newNgHandler(rt, &mockConversationStore{
		err: fmt.Errorf("not found"),
	})

	_, err := h.continueConversation(context.Background(), &ngAgentContinueArgs{ConversationID: 99, Input: "hello"})
	req.Error(err)
	req.Contains(err.Error(), "could not load conversation")
	req.False(rt.called)
}

func TestNgAgentHandler_ContinueRuntimeError(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntimeErr{}
	h := newNgHandler(rt, &mockConversationStore{
		conv: &types.AiConversation{AgentID: 5},
	})

	_, err := h.continueConversation(context.Background(), &ngAgentContinueArgs{ConversationID: 10, Input: "hello"})
	req.Error(err)
}

type mockRuntimeErr struct{}

func (m *mockRuntimeErr) Run(_ context.Context, _ *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error) {
	return nil, fmt.Errorf("runtime error")
}
