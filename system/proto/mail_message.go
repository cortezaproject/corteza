package proto

import (
	mail2 "net/mail"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/cortezaproject/corteza-server/system/types"
)

func NewMailMessage(mail *types.MailMessage) *MailMessage {
	if mail == nil {
		return nil
	}

	addrConv := func(aa []*mail2.Address) []*MailMessage_Header_MailAddress {
		out := make([]*MailMessage_Header_MailAddress, len(aa))

		for i, a := range aa {
			out[i] = &MailMessage_Header_MailAddress{
				Address: a.Address,
				Name:    a.Name,
			}
		}

		return out
	}

	hConv := func(hh mail2.Header) map[string]*MailMessage_Header_HeaderValues {
		out := make(map[string]*MailMessage_Header_HeaderValues, len(hh))

		for k, vv := range hh {
			out[k] = &MailMessage_Header_HeaderValues{Values: vv}
		}

		return out
	}

	var p = &MailMessage{
		Header: &MailMessage_Header{
			Date:    &timestamp.Timestamp{Seconds: mail.Date.Unix()},
			To:      addrConv(mail.Header.To),
			Cc:      addrConv(mail.Header.CC),
			Bcc:     addrConv(mail.Header.BCC),
			From:    addrConv(mail.Header.From),
			ReplyTo: addrConv(mail.Header.ReplyTo),
			Raw:     hConv(mail.Header.Raw),
		},
		RawBody: mail.RawBody,
	}

	return p
}
