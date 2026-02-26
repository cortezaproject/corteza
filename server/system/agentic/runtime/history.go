package runtime

import (
	"context"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	// ConversationStore manages conversation persistence.
	// Satisfied by system/service.AiConversation().
	ConversationStore interface {
		FindByID(ctx context.Context, ID uint64) (*types.AiConversation, error)
		Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
		Update(ctx context.Context, upd *types.AiConversation) (*types.AiConversation, error)
		DeleteByID(ctx context.Context, ID uint64) error
	}
)
