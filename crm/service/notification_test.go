package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	gomail "gopkg.in/mail.v2"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/system/types"
)

func TestUserRefExpanding(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ntf := &notification{}

	// msg := &gomail.Message{}
	usr := &types.User{ID: 72932592256548967, Email: "user@mock.ed", Name: "Mocked User"}

	usrSvc := NewMocknotificationUserService(mockCtrl)
	usrSvc.EXPECT().FindByID(usr.ID).Times(1).Return(usr, nil)
	ntf.userSvc = usrSvc

	input := []string{"sample@domain.tld", "sample@domain.tld Name", "72932592256548967"}
	rcpts, err := ntf.expandUserRefs(usrSvc, input)
	test.ErrNil(t, err, "expandUserRefs returned an error: %v")
	test.Assert(t, len(rcpts) == len(input), "Expecting %d headers, got %d", len(input), len(rcpts))
	test.Assert(t, rcpts[2] == usr.Email+" "+usr.Name, "Expecting %d headers, got %d", len(input), len(rcpts))
}

func TestAttachEmailRecipients(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ntf := &notification{}

	usrSvc := NewMocknotificationUserService(mockCtrl)
	usrSvc.EXPECT().With(gomock.Any())
	ntf.userSvc = usrSvc

	msg := gomail.NewMessage()

	input := []string{"sample@domain.tld", "sample2@domain.tld First Name"}
	err := ntf.AttachEmailRecipients(msg, "To", input...)
	to := msg.GetHeader("To")

	test.ErrNil(t, err, "AttachEmailRecipients returned an error: %v")

	test.Assert(t, len(input) == len(to), "Expecting %d headers, got %d", len(input), len(to))

	test.Assert(t, to[0] == "sample@domain.tld", "Expecting address to match, got %v", to[0])

	test.Assert(t, to[1] == "\"First Name\" <sample2@domain.tld>", "Expecting address to match, got %v", to[1])
}
