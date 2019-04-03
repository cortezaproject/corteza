// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
)

func TestReaction(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := MessageFlag(context.Background(), factory.Database.MustGet())

	tx(t, func() error {
		var chID = factory.Sonyflake.NextID()
		var msgID = factory.Sonyflake.NextID()
		var f *types.MessageFlag
		var ff types.MessageFlagSet
		f, err = rpo.Create(&types.MessageFlag{
			ChannelID: chID,
			MessageID: msgID,
			UserID:    3,
			Flag:      "success",
		})

		test.Assert(t, err == nil, "Should create message flag without an error, got: %+v", err)

		f, err = rpo.FindByID(f.ID)
		test.Assert(t, err == nil, "Should fetch message flag without an error, got: %+v", err)
		test.Assert(t, f != nil && f.ChannelID == chID, "fetch should return valid type struct")

		ff, err = rpo.FindByMessageIDs(msgID)
		test.Assert(t, err == nil, "Should fetch message flag by range without an error, got: %+v", err)
		test.Assert(t, len(ff) == 1, "fetch by range should return 1 message")

		err = rpo.DeleteByID(f.ID)
		test.Assert(t, err == nil, "Should delete message flag without an error, got: %+v", err)

		return nil
	})
}
