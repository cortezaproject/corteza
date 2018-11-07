package service

import (
	"context"
	"strings"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/sam/types"
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

func TestMessageLength(t *testing.T) {
	// mockCtrl := gomock.NewController(t)
	// defer mockCtrl.Finish()

	ctx := context.TODO()
	auth.SetIdentityToContext(ctx, &systemTypes.User{})

	svc := message{db: &mockDB{}, ctx: ctx}
	e := func(out *types.Message, err error) error { return err }

	longText := strings.Repeat("X", settingsMessageBodyLength+1)

	assert(t, e(svc.Create(&types.Message{})) != nil, "Should not allow to create unnamed channels")
	assert(t, e(svc.Create(&types.Message{Message: longText})) != nil, "Should not allow to create channel with really long name")
}

func TestMentionsExtraction(t *testing.T) {
	var (
		svc   = message{}
		mm    types.MentionSet
		cases = []struct {
			text string
			ids  []uint64
		}{
			{"abcde",
				[]uint64{}},
			{"<@4095834095>",
				[]uint64{4095834095}},
			{"<@4095834095> <@4095834095>",
				[]uint64{4095834095}},
			{"<@4095834095> <@4095834097>",
				[]uint64{4095834095, 4095834097}},
			{"dfsf<@4095834095>dsfsd<@4095834097>sdfs",
				[]uint64{4095834095, 4095834097}},
			{"dfsf<@4095834095>dsfsd<@40958340dfsZ",
				[]uint64{4095834095}},
			{"<@4095834095 label> <@4095834097>",
				[]uint64{4095834095, 4095834097}},
		}
	)

	for _, c := range cases {
		mm = svc.extractMentions(&types.Message{Message: c.text})

		assert(t, len(mm) == len(c.ids), "Number of extracted (%d) and expected (%d) user IDs do not match (%s)", len(mm), len(c.ids), c.text)

		for _, id := range c.ids {
			assert(t, len(mm.FindByUserID(id)) == 1, "User ID (%d) was not extracted (%s)", id, c.text)
		}
	}
}
