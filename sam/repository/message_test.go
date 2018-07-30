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

	var mm []*types.Message

	msg.Message = msg1

	msg, err = rpo.CreateMessage(msg)
	must(t, err)
	if msg.Message != msg1 {
		t.Fatal("Changes were not stored")
	}

	msg.Message = msg2

	msg, err = rpo.UpdateMessage(msg)
	must(t, err)
	if msg.Message != msg2 {
		t.Fatal("Changes were not stored")
	}

	msg, err = rpo.FindMessageByID(msg.ID)
	must(t, err)
	if msg.Message != msg2 {
		t.Fatal("Changes were not stored")
	}

	mm, err = rpo.FindMessages(&types.MessageFilter{Query: msg2})
	must(t, err)
	if len(mm) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.DeleteMessageByID(msg.ID))
}
