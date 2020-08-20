package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"
	"testing"
)

func testSettings(t *testing.T, s store.Settings) {
	var (
		ctx = context.Background()
		_   = ctx
	)

	t.Run("create", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("update", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("delete/undelete", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("search", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("search by *", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
