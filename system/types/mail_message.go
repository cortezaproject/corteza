package types

import (
	"net/mail"
	"time"
)

type (
	MailMessage struct {
		Date time.Time `json:"date"`

		Header MailMessageHeader `json:"header"`

		// RawBody will be base64 encoded!
		// (might contain binary data)
		RawBody []byte `json:"rawBody,string"`

		// @todo parts
		// Parts []...
	}

	MailMessageHeader struct {
		// extract common addresses
		To      []*mail.Address `json:"to"`
		CC      []*mail.Address `json:"cc"`
		BCC     []*mail.Address `json:"bcc"`
		From    []*mail.Address `json:"from"`
		ReplyTo []*mail.Address `json:"replyTo"`

		Raw mail.Header `json:"raw"`
	}
)
