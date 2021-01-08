package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMessageLength(t *testing.T) {
	var (
		ctx = context.Background()
		svc = message{}
		e   = func(out *types.Message, err error) error { return err }
	)

	require.True(t, e(svc.Create(ctx, &types.Message{})) != nil, "Should not allow to create empty message")

	if settingsMessageBodyLength > 0 {
		longText := strings.Repeat("X", settingsMessageBodyLength+1)
		require.True(t, e(svc.Create(ctx, &types.Message{Message: longText})) != nil, "Should not allow to create message with really long text")
	}
}

func TestMentionsExtraction(t *testing.T) {
	var (
		svc   = message{}
		mm    types.MentionSet
		cases = []struct {
			text string
			ids  []uint64
		}{
			{"abcde",
				[]uint64{}},
			{"<@4095834095>",
				[]uint64{4095834095}},
			{"<@4095834095> <@4095834095>",
				[]uint64{4095834095}},
			{"<@4095834095> <@4095834097>",
				[]uint64{4095834095, 4095834097}},
			{"dfsf<@4095834095>dsfsd<@4095834097>sdfs",
				[]uint64{4095834095, 4095834097}},
			{"dfsf<@4095834095>dsfsd<@40958340dfsZ",
				[]uint64{4095834095}},
			{"<@4095834095 label> <@4095834097>",
				[]uint64{4095834095, 4095834097}},
		}
	)

	for _, c := range cases {
		mm = svc.extractMentions(&types.Message{Message: c.text})

		require.True(t, len(mm) == len(c.ids), "Number of extracted (%d) and expected (%d) user IDs do not match (%s)", len(mm), len(c.ids), c.text)

		for _, id := range c.ids {
			require.True(t, len(mm.FindByUserID(id)) == 1, "Owner ID (%d) was not extracted (%s)", id, c.text)
		}
	}
}
