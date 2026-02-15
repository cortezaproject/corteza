package runtime

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/system/types"
)

// mockLLM is a no-op LLM client for development/testing.
type mockLLM struct{}

func MockLLM() *mockLLM { return &mockLLM{} }

func (m *mockLLM) Chat(_ context.Context, _ string, _ []types.AiConversationMessage, _ []Tool, _ LLMConfig) (*LLMResponse, error) {
	return nil, fmt.Errorf("LLM not configured: using mock client")
}
