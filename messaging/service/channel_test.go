package service

import (
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/stretchr/testify/require"
)

func TestChannelNameTooShort(t *testing.T) {
	svc := channel{}
	e := func(out *types.Channel, err error) error { return err }

	require.True(t, e(svc.Create(&types.Channel{})) != nil, "Should not allow to create unnamed channels")

	if settingsChannelNameLength > 0 {
		longName := strings.Repeat("X", settingsChannelNameLength+1)
		require.True(t, e(svc.Create(&types.Channel{Name: longName})) != nil, "Should not allow to create channel with really long name")
	}
}
