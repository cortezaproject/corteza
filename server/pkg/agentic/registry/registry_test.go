package registry

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Add(t *testing.T) {
	reg := Registry()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		agent := &types.Agent{
			ID:     100,
			Handle: "test-agent",
		}

		created, err := reg.Add(ctx, agent)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), created.ID)
	})

	t.Run("error missing ID", func(t *testing.T) {
		agent := &types.Agent{
			Handle: "no-id",
		}
		_, err := reg.Add(ctx, agent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "agent ID is required")
	})

	t.Run("error duplicate ID", func(t *testing.T) {
		a1 := &types.Agent{ID: 200, Handle: "a1"}
		_, err := reg.Add(ctx, a1)
		require.NoError(t, err)

		a2 := &types.Agent{ID: 200, Handle: "a2"}
		_, err = reg.Add(ctx, a2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestRegistry_Get(t *testing.T) {
	reg := Registry()
	ctx := context.Background()

	agent := &types.Agent{ID: 100, Handle: "test-agent"}
	reg.Add(ctx, agent)

	t.Run("found", func(t *testing.T) {
		found, err := reg.Get(ctx, 100)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), found.ID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := reg.Get(ctx, 999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestRegistry_Remove(t *testing.T) {
	reg := Registry()
	ctx := context.Background()

	agent := &types.Agent{ID: 100, Handle: "test-agent"}
	reg.Add(ctx, agent)

	t.Run("success", func(t *testing.T) {
		err := reg.Remove(ctx, 100)
		require.NoError(t, err)

		_, err = reg.Get(ctx, 100)
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		err := reg.Remove(ctx, 999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("cannot remove active", func(t *testing.T) {
		a2 := &types.Agent{ID: 200, Handle: "active-agent", Status: "draft"}
		reg.Add(ctx, a2)
		reg.Activate(ctx, 200)

		err := reg.Remove(ctx, 200)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "active")
	})
}

func TestRegistry_Activate_Disable(t *testing.T) {
	reg := Registry()
	ctx := context.Background()
	agent := &types.Agent{ID: 100, Handle: "a1", Status: "draft"}
	reg.Add(ctx, agent)

	// Activate
	err := reg.Activate(ctx, 100)
	require.NoError(t, err)

	updated, _ := reg.Get(ctx, 100)
	assert.Equal(t, "active", updated.Status)
	assert.NotNil(t, updated.UpdatedAt)
	t1 := *updated.UpdatedAt

	// Disable
	time.Sleep(1 * time.Millisecond) // ensure timestamp change
	err = reg.Disable(ctx, 100)
	require.NoError(t, err)

	updated, _ = reg.Get(ctx, 100)
	assert.Equal(t, "disabled", updated.Status)
	assert.NotEqual(t, t1, *updated.UpdatedAt)

	// Not found checks
	require.Error(t, reg.Activate(ctx, 999))
	require.Error(t, reg.Disable(ctx, 999))
}
