package repository

import (
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestMessage(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
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
			assert(t, err == nil, "FFindMessageByID error: %v", err)
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
