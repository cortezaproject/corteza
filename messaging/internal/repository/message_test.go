// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
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
		ch, err = chRpo.Create(ch)
		ch.Type = types.ChannelTypePublic

		msg := &types.Message{ChannelID: ch.ID}

		msg.Message = msg1
		msg, err = msgRpo.Create(msg)
		test.Assert(t, err == nil, "CreateMessage error: %+v", err)
		test.Assert(t, msg.Message == msg1, "Changes were not stored")

		{
			msg.Message = msg2
			msg, err = msgRpo.Update(msg)
			test.Assert(t, err == nil, "UpdateMessage error: %+v", err)
			test.Assert(t, msg.Message == msg2, "Changes were not stored")
		}

		{
			msg, err = msgRpo.FindByID(msg.ID)
			test.Assert(t, err == nil, "FindMessageByID error: %+v", err)
			test.Assert(t, msg.Message == msg2, "Changes were not stored")
		}

		{
			mm, err = msgRpo.Find(&types.MessageFilter{Query: msg2})
			test.Assert(t, err == nil, "FindMessages error: %+v", err)
			test.Assert(t, len(mm) > 0, "No results found")
		}

		{
			err = msgRpo.DeleteByID(msg.ID)
			test.Assert(t, err == nil, "DeleteMessageByID error: %+v", err)
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
		ch, err = chRpo.Create(ch)
		ch.Type = types.ChannelTypePublic

		msg := &types.Message{ChannelID: ch.ID}
		rpl := &types.Message{ChannelID: ch.ID}

		msg, err = msgRpo.Create(msg)
		test.Assert(t, err == nil, "CreateMessage error: %+v", err)
		test.Assert(t, msg.ID > 0, "Message did not get its ID")

		rpl.ReplyTo = msg.ID
		rpl, err = msgRpo.Create(rpl)
		test.Assert(t, err == nil, "CreateMessage error: %+v", err)
		test.Assert(t, rpl.ID > 0, "Reply did not get its ID")

		// Let's increase this so that FindThreads
		// can include it into results
		msgRpo.IncReplyCount(msg.ID)

		{
			mm, err = msgRpo.Find(&types.MessageFilter{
				RepliesTo: msg.ID,
				ChannelID: ch.ID,
			})

			test.Assert(t, err == nil, "FindMessages error: %+v", err)
			test.Assert(t, len(mm) == 1, "Failed to fetch only reply, got: %d", len(mm))
			test.Assert(t, mm[0].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = msgRpo.FindThreads(&types.MessageFilter{
				ChannelID: ch.ID,
			})

			test.Assert(t, err == nil, "FindThreads error: %+v", err)
			test.Assert(t, len(mm) == 2, "Failed to fetch messages in threads (2 messages), got: %d", len(mm))
			test.Assert(t, mm[0].ID == msg.ID, "Original message ID does not match")
			test.Assert(t, mm[1].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = msgRpo.Find(&types.MessageFilter{
				ChannelID: ch.ID,
			})

			test.Assert(t, err == nil, "FindMessages error: %+v", err)
			test.Assert(t, len(mm) == 1, "Failed to fetch only original message")
			test.Assert(t, mm[0].ID == msg.ID, "Reply ID does not match")
		}

		{

			test.Assert(t, msgRpo.IncReplyCount(msg.ID) == nil, "IncReplyCount should not return an error")
			test.Assert(t, msgRpo.IncReplyCount(msg.ID) == nil, "IncReplyCount should not return an error")
			// +1 that we have from before

			msg, err = msgRpo.FindByID(msg.ID)
			test.Assert(t, err == nil, "FindMessageByID error: %+v", err)
			test.Assert(t, msg.Replies == 3, "Reply counter check failed, expecting 3, got %d", msg.Replies)

			test.Assert(t, msgRpo.DecReplyCount(msg.ID) == nil, "DecReplyCount should not return an error")

			msg, err = msgRpo.FindByID(msg.ID)
			test.Assert(t, err == nil, "FindMessageByID error: %+v", err)
			test.Assert(t, msg.Replies == 2, "Reply counter check failed, expecting 1, got %d", msg.Replies)
		}

		return nil
	})
}
