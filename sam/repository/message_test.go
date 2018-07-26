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
	att := types.Message{}.New()

	var msg1, msg2 = "Test message v1", "Test message v2"

	var aa []*types.Message

	att.SetMessage(msg1)

	att, err = rpo.Create(ctx, att)
	must(t, err)
	if att.Message != msg1 {
		t.Fatal("Changes were not stored")
	}

	att.SetMessage(msg2)

	att, err = rpo.Update(ctx, att)
	must(t, err)
	if att.Message != msg2 {
		t.Fatal("Changes were not stored")
	}

	att, err = rpo.FindByID(ctx, att.ID)
	must(t, err)
	if att.Message != msg2 {
		t.Fatal("Changes were not stored")
	}

	aa, err = rpo.Find(ctx, &types.MessageFilter{Query: msg2})
	must(t, err)
	if len(aa) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.Delete(ctx, att.ID))
}
