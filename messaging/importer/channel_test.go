package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

func TestChannelImport_CastSet(t *testing.T) {
	impFixTester(t, "channels", func(t *testing.T, ri *Channel) {
		req := require.New(t)
		req.NotNil(ri.set)
		req.Len(ri.set, 3)

		req.NotNil(ri.set.FindByName("General"))
		req.Equal("Talk about anything", ri.set.FindByName("General").Topic)
		req.Equal(types.ChannelTypePublic, ri.set.FindByName("General").Type)

		req.NotNil(ri.set.FindByName("Random"))
		req.Equal("", ri.set.FindByName("Random").Topic)
		req.Equal(types.ChannelTypePublic, ri.set.FindByName("Random").Type)

		req.NotNil(ri.set.FindByName("Secret"))
		req.Equal(types.ChannelTypePrivate, ri.set.FindByName("Secret").Type)
	})
}
