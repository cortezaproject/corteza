package repository

import (
	"context"

	"github.com/titpetric/factory"

	"testing"

	"github.com/crusttech/crust/sam/types"
)

func TestMessage(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Message(context.Background(), factory.Database.MustGet())
	msg := &types.Message{}

	var msg1, msg2 = "Test message v1", "Test message v2"

	var mm types.MessageSet

	{
		msg.Message = msg1
		msg, err = rpo.CreateMessage(msg)
		assert(t, err == nil, "CreateMessage error: %v", err)
		assert(t, msg.Message == msg1, "Changes were not stored")

		{
			msg.Message = msg2
			msg, err = rpo.UpdateMessage(msg)
			assert(t, err == nil, "UpdateMessage error: %v", err)
			assert(t, msg.Message == msg2, "Changes were not stored")
		}

		{
			msg, err = rpo.FindMessageByID(msg.ID)
			assert(t, err == nil, "FindMessageByID error: %v", err)
			assert(t, msg.Message == msg2, "Changes were not stored")
		}

		{
			mm, err = rpo.FindMessages(&types.MessageFilter{Query: msg2})
			assert(t, err == nil, "FindMessages error: %v", err)
			assert(t, len(mm) > 0, "No results found")
		}

		{
			err = rpo.DeleteMessageByID(msg.ID)
			assert(t, err == nil, "DeleteMessageByID error: %v", err)
		}
	}
}

func TestReplies(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	chID := factory.Sonyflake.NextID()

	rpo := Message(context.Background(), factory.Database.MustGet())
	msg := &types.Message{ChannelID: chID}
	rpl := &types.Message{ChannelID: chID}

	var mm types.MessageSet

	tx(t, func() error {
		msg, err = rpo.CreateMessage(msg)
		assert(t, err == nil, "CreateMessage error: %v", err)
		assert(t, msg.ID > 0, "Message did not get its ID")

		rpl.ReplyTo = msg.ID
		rpl, err = rpo.CreateMessage(rpl)
		assert(t, err == nil, "CreateMessage error: %v", err)
		assert(t, rpl.ID > 0, "Reply did not get its ID")

		{
			mm, err = rpo.FindMessages(&types.MessageFilter{
				RepliesTo: msg.ID,
				ChannelID: chID,
			})

			assert(t, err == nil, "FindMessages error: %v", err)
			assert(t, len(mm) == 1, "Failed to fetch only reply")
			assert(t, mm[0].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = rpo.FindMessages(&types.MessageFilter{
				ChannelID: chID,
			})

			assert(t, err == nil, "FindMessages error: %v", err)
			assert(t, len(mm) == 1, "Failed to fetch only original message")
			assert(t, mm[0].ID == msg.ID, "Reply ID does not match")
		}

		{
			rpo.IncReplyCount(msg.ID)
			rpo.IncReplyCount(msg.ID)
			rpo.IncReplyCount(msg.ID)

			msg, err = rpo.FindMessageByID(msg.ID)
			assert(t, err == nil, "FindMessageByID error: %v", err)
			assert(t, msg.Replies == 3, "Reply counter check failed, expecting 3, got %v", msg.Replies)

			rpo.DecReplyCount(msg.ID)
			rpo.DecReplyCount(msg.ID)

			msg, err = rpo.FindMessageByID(msg.ID)
			assert(t, err == nil, "FindMessageByID error: %v", err)
			assert(t, msg.Replies == 1, "Reply counter check failed, expecting 1, got %v", msg.Replies)
		}

		return nil
	})
}
