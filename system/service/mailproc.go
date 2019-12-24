package service

import (
	"context"
	"io"
	"io/ioutil"
	"net/mail"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	mailproc struct {
		sr mailprocScriptsRunner

		logger *zap.Logger
	}

	mailprocScriptsRunner interface {
		OnReceiveMailMessage(ctx context.Context, message *types.MailMessage) (err error)
	}
)

func Mailproc() *mailproc {
	return &mailproc{
		logger: DefaultLogger.Named("mailproc"),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc mailproc) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc mailproc) ContentProcessor(ctx context.Context, m io.Reader) error {
	if m, err := mailProcMessage(m); err != nil {
		return err
	} else if err = svc.sr.OnReceiveMailMessage(ctx, m); err != nil {
		return err
	}

	return nil
}

func mailProcMessage(r io.Reader) (out *types.MailMessage, err error) {
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

	out = &types.MailMessage{}

	out.Header.Raw = msg.Header

	out.Date, _ = msg.Header.Date()

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
