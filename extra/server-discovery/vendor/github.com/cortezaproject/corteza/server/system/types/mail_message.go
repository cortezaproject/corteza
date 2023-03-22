package types

import (
	"io"
	"io/ioutil"
	"net/mail"
	"time"
)

type (
	MailMessage struct {
		Date time.Time `json:"date"`

		Subject string `json:"subject"`

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

func NewMailMessage(r io.Reader) (out *MailMessage, err error) {
	var (
		aa  []*mail.Address
		msg *mail.Message

		addrKeys = []string{
			"from",
			"to",
			"cc",
			"bcc",
			"reply-to",
		}
	)

	if msg, err = mail.ReadMessage(r); err != nil {
		return
	}

	out = &MailMessage{}

	out.Header.Raw = msg.Header

	out.Date, _ = msg.Header.Date()
	out.Subject = msg.Header.Get("subject")

	for _, key := range addrKeys {
		aa, err = msg.Header.AddressList(key)

		if err != nil && err != mail.ErrHeaderNotPresent {
			return
		}

		if len(aa) > 0 {
			switch key {
			case "from":
				out.Header.From = aa
			case "to":
				out.Header.To = aa
			case "cc":
				out.Header.CC = aa
			case "bcc":
				out.Header.BCC = aa
			case "reply-to":
				out.Header.ReplyTo = aa
			}
		}
	}

	if out.RawBody, err = ioutil.ReadAll(msg.Body); err != nil {
		return
	}

	return
}
