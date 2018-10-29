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

	msgRpo := Message(context.Background(), factory.Database.MustGet())
	chRpo := Channel(context.Background(), factory.Database.MustGet())

	var msg1, msg2 = "Test message v1", "Test message v2"

	var mm types.MessageSet

	tx(t, func() error {
		ch := &types.Channel{}
		ch, err = chRpo.CreateChannel(ch)
		ch.Type = types.ChannelTypePublic

		msg := &types.Message{ChannelID: ch.ID}

		msg.Message = msg1
		msg, err = msgRpo.CreateMessage(msg)
		assert(t, err == nil, "CreateMessage error: %v", err)
		assert(t, msg.Message == msg1, "Changes were not stored")

		{
			msg.Message = msg2
			msg, err = msgRpo.UpdateMessage(msg)
			assert(t, err == nil, "UpdateMessage error: %v", err)
			assert(t, msg.Message == msg2, "Changes were not stored")
		}

		{
			msg, err = msgRpo.FindMessageByID(msg.ID)
			assert(t, err == nil, "FindMessageByID error: %v", err)
			assert(t, msg.Message == msg2, "Changes were not stored")
		}

		{
			mm, err = msgRpo.FindMessages(&types.MessageFilter{Query: msg2})
			assert(t, err == nil, "FindMessages error: %v", err)
			assert(t, len(mm) > 0, "No results found")
		}

		{
			err = msgRpo.DeleteMessageByID(msg.ID)
			assert(t, err == nil, "DeleteMessageByID error: %v", err)
		}

		return nil
	})
}

func TestReplies(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	msgRpo := Message(context.Background(), factory.Database.MustGet())
	chRpo := Channel(context.Background(), factory.Database.MustGet())

	var mm types.MessageSet

	tx(t, func() error {
		ch := &types.Channel{}
		ch, err = chRpo.CreateChannel(ch)
		ch.Type = types.ChannelTypePublic

		msg := &types.Message{ChannelID: ch.ID}
		rpl := &types.Message{ChannelID: ch.ID}

		msg, err = msgRpo.CreateMessage(msg)
		assert(t, err == nil, "CreateMessage error: %v", err)
		assert(t, msg.ID > 0, "Message did not get its ID")

		rpl.ReplyTo = msg.ID
		rpl, err = msgRpo.CreateMessage(rpl)
		assert(t, err == nil, "CreateMessage error: %v", err)
		assert(t, rpl.ID > 0, "Reply did not get its ID")

		// Let's increase this so that FindThreads
		// can include it into results
		msgRpo.IncReplyCount(msg.ID)

		{
			mm, err = msgRpo.FindMessages(&types.MessageFilter{
				RepliesTo: msg.ID,
				ChannelID: ch.ID,
			})

			assert(t, err == nil, "FindMessages error: %v", err)
			assert(t, len(mm) == 1, "Failed to fetch only reply, got: %d", len(mm))
			assert(t, mm[0].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = msgRpo.FindThreads(&types.MessageFilter{
				ChannelID: ch.ID,
			})

			assert(t, err == nil, "FindThreads error: %v", err)
			assert(t, len(mm) == 2, "Failed to fetch messages in threads (2 messages), got: %d", len(mm))
			assert(t, mm[0].ID == msg.ID, "Original message ID does not match")
			assert(t, mm[1].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = msgRpo.FindMessages(&types.MessageFilter{
				ChannelID: ch.ID,
			})

			assert(t, err == nil, "FindMessages error: %v", err)
			assert(t, len(mm) == 1, "Failed to fetch only original message")
			assert(t, mm[0].ID == msg.ID, "Reply ID does not match")
		}

		{

			assert(t, msgRpo.IncReplyCount(msg.ID) == nil, "IncReplyCount should not return an error")
			assert(t, msgRpo.IncReplyCount(msg.ID) == nil, "IncReplyCount should not return an error")
			// +1 that we have from before

			msg, err = msgRpo.FindMessageByID(msg.ID)
			assert(t, err == nil, "FindMessageByID error: %v", err)
			assert(t, msg.Replies == 3, "Reply counter check failed, expecting 3, got %v", msg.Replies)

			assert(t, msgRpo.DecReplyCount(msg.ID) == nil, "DecReplyCount should not return an error")

			msg, err = msgRpo.FindMessageByID(msg.ID)
			assert(t, err == nil, "FindMessageByID error: %v", err)
			assert(t, msg.Replies == 2, "Reply counter check failed, expecting 1, got %v", msg.Replies)
		}

		return nil
	})
}
