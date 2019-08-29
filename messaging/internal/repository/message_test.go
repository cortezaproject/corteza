// +build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/titpetric/factory"
	dbLogger "github.com/titpetric/factory/logger"

	"github.com/cortezaproject/corteza-server/internal/test"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func TestMessage(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	msgRpo := Message(context.Background(), factory.Database.MustGet("messaging"))
	chRpo := Channel(context.Background(), factory.Database.MustGet("messaging"))

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

func TestBeforeMessageID(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet("messaging")
	msgRpo := Message(context.Background(), db)
	chRpo := Channel(context.Background(), db)

	tx(t, func() error {
		// insert 1 channel
		ch := &types.Channel{}
		ch, err = chRpo.Create(ch)
		ch.Type = types.ChannelTypePublic

		// insert 100 messages
		db.SetLogger(dbLogger.Silent{})
		messages := make([]*types.Message, 100)
		for k, _ := range messages {
			messages[k], err = msgRpo.Create(&types.Message{
				ChannelID: ch.ID,
				Message:   fmt.Sprintf("#%d: Lorem ipsum dolor sit amet", k),
			})

			test.Assert(t, err == nil, "CreateMessage error: %+v", err)
		}
		db.SetLogger(dbLogger.Default{})

		// request last 10 messages from channel
		lastPageRequest := &types.MessageFilter{
			ChannelID: []uint64{ch.ID},
			Limit:     10,
		}

		var lastPage types.MessageSet
		lastPage, err = msgRpo.Find(lastPageRequest)

		test.Assert(t, err == nil, "lastPageRequest error: %+v", err)
		test.Assert(t, len(lastPage) > 0, "No results found (last page)")

		// request previous 10 messages from channel
		prevPageRequest := &types.MessageFilter{
			ChannelID: []uint64{ch.ID},
			Limit:     10,
			BeforeID:  lastPage[9].ID,
		}

		var prevPage types.MessageSet
		prevPage, err = msgRpo.Find(prevPageRequest)

		test.Assert(t, err == nil, "prevPageRequest error: %+v", err)
		test.Assert(t, prevPage[0].ID != messages[0].ID, "We have 100 IDs, second page shouldn't start with first ID")
		test.Assert(t, prevPage[0].ID == messages[89].ID, "ID should match index 89 (max index - 10), but %d != %d", prevPage[0].ID, messages[89].ID)
		test.Assert(t, len(prevPage) > 0, "No results found (previous page)")

		return nil
	})
}

func TestReplies(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	msgRpo := Message(context.Background(), factory.Database.MustGet("messaging"))
	chRpo := Channel(context.Background(), factory.Database.MustGet("messaging"))

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
				ThreadID:  []uint64{msg.ID},
				ChannelID: []uint64{ch.ID},
			})

			test.Assert(t, err == nil, "FindMessages error: %+v", err)
			test.Assert(t, len(mm) == 1, "Failed to fetch only reply, got: %d", len(mm))
			test.Assert(t, mm[0].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = msgRpo.FindThreads(&types.MessageFilter{
				ChannelID: []uint64{ch.ID},
			})

			test.Assert(t, err == nil, "FindThreads error: %+v", err)
			test.Assert(t, len(mm) == 2, "Failed to fetch messages in threads (2 messages), got: %d", len(mm))
			test.Assert(t, mm[0].ID == msg.ID, "Original message ID does not match")
			test.Assert(t, mm[1].ID == rpl.ID, "Reply ID does not match")
		}

		{
			mm, err = msgRpo.Find(&types.MessageFilter{
				ChannelID: []uint64{ch.ID},
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
