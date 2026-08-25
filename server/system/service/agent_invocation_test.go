package service

import (
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
)

// An agent that names no way of being invoked cannot run at all, and the only
// sign of it is the exec call refusing much later. The editor's own model
// starts with user invocation on; this keeps the API agreeing with it.
func TestDefaultInvocation(t *testing.T) {
	t.Run("an unspecified invocation enables the user", func(t *testing.T) {
		inv := types.AgentInvocation{}
		defaultInvocation(&inv)
		assert.True(t, inv.User.Enabled)
		assert.False(t, inv.System.Enabled)
	})

	t.Run("a system-only agent is left alone", func(t *testing.T) {
		inv := types.AgentInvocation{System: types.AgentInvocationSystem{Enabled: true}}
		defaultInvocation(&inv)
		assert.False(t, inv.User.Enabled, "a deliberate system-only agent must not become user-invocable")
		assert.True(t, inv.System.Enabled)
	})

	t.Run("an explicit user invocation is untouched", func(t *testing.T) {
		inv := types.AgentInvocation{User: types.AgentInvocationUser{Enabled: true}}
		defaultInvocation(&inv)
		assert.True(t, inv.User.Enabled)
	})
}
