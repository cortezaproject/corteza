package proto

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

func FromMailMessage(mail *types.MailMessage) *MailMessage {
	if mail == nil {
		return nil
	}

	var p = &MailMessage{}

	panic("@todo implement mailmessage => proto conv!")
	return p
}
