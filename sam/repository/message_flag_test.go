package repository

import (
	"context"

	"github.com/titpetric/factory"

	"testing"

	"github.com/crusttech/crust/sam/types"
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

		assert(t, err == nil, "Should create message flag without an error, got: %v", err)

		f, err = rpo.FindByID(f.ID)
		assert(t, err == nil, "Should fetch message flag without an error, got: %v", err)
		assert(t, f != nil && f.ChannelID == chID, "fetch should return valid type struct")

		ff, err = rpo.FindByMessageIDs(msgID)
		assert(t, err == nil, "Should fetch message flag by range without an error, got: %v", err)
		assert(t, len(ff) == 1, "fetch by range should return 1 message")

		err = rpo.DeleteByID(f.ID)
		assert(t, err == nil, "Should delete message flag without an error, got: %v", err)

		return nil
	})
}
