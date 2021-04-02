package automation

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/expr"
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
