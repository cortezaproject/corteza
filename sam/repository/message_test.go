package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestMessage(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Message()
	ctx := context.Background()
	msg := &types.Message{}

	var msg1, msg2 = "Test message v1", "Test message v2"

	var mm []*types.Message

	msg.Message = msg1

	msg, err = rpo.CreateMessage(ctx, msg)
	must(t, err)
	if msg.Message != msg1 {
		t.Fatal("Changes were not stored")
	}

	msg.Message = msg2

	msg, err = rpo.Update(ctx, msg)
	must(t, err)
	if msg.Message != msg2 {
		t.Fatal("Changes were not stored")
	}

	msg, err = rpo.FindMessageByID(ctx, msg.ID)
	must(t, err)
	if msg.Message != msg2 {
		t.Fatal("Changes were not stored")
	}

	mm, err = rpo.FindMessages(ctx, &types.MessageFilter{Query: msg2})
	must(t, err)
	if len(mm) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.Delete(ctx, msg.ID))
}
