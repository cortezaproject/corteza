package service

import (
	"context"
	"strings"
	"testing"

	authTypes "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/sam/types"
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

func TesMessageLength(t *testing.T) {
	// mockCtrl := gomock.NewController(t)
	// defer mockCtrl.Finish()

	ctx := context.TODO()
	auth.SetIdentityToContext(ctx, &authTypes.User{})

	svc := message{db: &mockDB{}, ctx: ctx}
	e := func(out *types.Message, err error) error { return err }

	longText := strings.Repeat("X", settingsMessageBodyLength+1)

	assert(t, e(svc.Create(&types.Message{})) != nil, "Should not allow to create unnamed channels")
	assert(t, e(svc.Create(&types.Message{Message: longText})) != nil, "Should not allow to create channel with really long name")
}
