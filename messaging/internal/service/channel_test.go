// +build unit

package service

import (
	"context"
	"strings"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestChannelNameTooShort(t *testing.T) {
	ctx := context.Background()
	ctx = auth.SetIdentityToContext(ctx, &systemTypes.User{})

	svc := channel{db: &mockDB{}, ctx: ctx}
	e := func(out *types.Channel, err error) error { return err }

	longName := strings.Repeat("X", settingsChannelNameLength+1)

	test.Assert(t, e(svc.Create(&types.Channel{})) != nil, "Should not allow to create unnamed channels")
	test.Assert(t, e(svc.Create(&types.Channel{Name: longName})) != nil, "Should not allow to create channel with really long name")
}
