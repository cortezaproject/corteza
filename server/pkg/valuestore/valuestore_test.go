package valuestore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSettingEnv(t *testing.T) {
	s := New()
	s.SetEnv(map[string]any{
		"k1": "v1",
		"k2": "v2",
	})

	require.Equal(t, "v1", s.Env("k1"))
	require.Equal(t, "v2", s.Env("k2"))
}

func TestPreventResettingEnv(t *testing.T) {
	s := New()
	s.SetEnv(map[string]any{
		"k1": "v1",
		"k2": "v2",
	})

	require.Panics(t, func() {
		s.SetEnv(map[string]any{
			"k3": "v3",
			"k4": "v4",
		})
	})
}

func TestPreventUninitializedAccess(t *testing.T) {
	s := New()

	require.Panics(t, func() {
		s.Env("k1")
	})
}
