package service

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

func loadLlmProvider(ctx context.Context, s store.LlmProviders, id uint64) (*types.LlmProvider, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid LLM provider ID")
	}
	return store.LookupLlmProviderByID(ctx, s, id)
}
