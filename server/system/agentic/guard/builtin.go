package guard

import (
	"context"

	"github.com/crusttech/human/server/pkg/inputguard"
	"github.com/crusttech/human/server/system/types"
)

// BuiltinGuard is a thin adapter over the inputguard package that implements
// GuardService. Swap out the underlying Check call to plug in a different implementation.
type BuiltinGuard struct {
	maxInputLength int
}

func NewBuiltinGuard() *BuiltinGuard {
	return &BuiltinGuard{maxInputLength: inputguard.DefaultMaxLength}
}

func (g *BuiltinGuard) CheckInput(_ context.Context, input string, _ []types.AiConversationMessage) (*GuardResult, error) {
	r := inputguard.CheckWithMaxLength(input, g.maxInputLength)
	if !r.Blocked {
		return safe(), nil
	}
	return blocked(r.Category, r.Reason), nil
}
