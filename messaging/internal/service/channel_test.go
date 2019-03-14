package service

import (
	"context"
	"strings"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/messaging/types"
	systemTypes "github.com/crusttech/crust/system/types"
)

// func TestChannelCreation(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	chRpoMock := NewMockRepository(mockCtrl)
// 	chRpoMock.EXPECT().WithCtx(gomock.Any()).AnyTimes().Return(chRpoMock)
// 	chRpoMock.EXPECT().
// 		FindUserByID(usr.ID).
// 		Times(1).
// 		Return(usr, nil)
//
// 	svc := channel{
// 		channel:
// 	}
//
// 	svc.Create()
// }

func TestChannelNameTooShort(t *testing.T) {
	// mockCtrl := gomock.NewController(t)
	// defer mockCtrl.Finish()

	ctx := context.TODO()
	auth.SetIdentityToContext(ctx, &systemTypes.User{})

	svc := channel{db: &mockDB{}, ctx: ctx}
	e := func(out *types.Channel, err error) error { return err }

	longName := strings.Repeat("X", settingsChannelNameLength+1)

	assert(t, e(svc.Create(&types.Channel{})) != nil, "Should not allow to create unnamed channels")
	assert(t, e(svc.Create(&types.Channel{Name: longName})) != nil, "Should not allow to create channel with really long name")
}
