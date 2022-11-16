package automation

import (
	"fmt"
	"io"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/http"
	"gopkg.in/mail.v2"
)

type (
	emailMessage struct {
		// only basic implementation for now
		// we can manipulate mail.Message internals through
		// specialized wf functions (message, setSubject, setHeaders, ...)
		msg *mail.Message
	}
)

func CastToEmailMessage(val interface{}) (out *emailMessage, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *emailMessage:
		if val.msg == nil {
			val.msg = mail.NewMessage()
		}

		return val, nil
	case nil:
		return &emailMessage{msg: mail.NewMessage()}, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func ReadRequestBody(in interface{}) (s string) {
	var (
		b   []byte
		err error
	)

	switch val := in.(type) {
	case *http.Request:
		b, err = io.ReadAll(val.Body)
	case io.Reader:
		b, err = io.ReadAll(val)
	default:
		b = []byte{}
	}

	if err != nil {
		return ""
	}

	return string(b)
}
