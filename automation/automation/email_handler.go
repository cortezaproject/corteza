package automation

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/mail"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	gomail "gopkg.in/mail.v2"
)

type (
	emailHandler struct {
		reg emailHandlerRegistry
	}

	messageArgs interface {
		GetReplyTo() (bool, string, *sysTypes.User)
		GetFrom() (bool, string, *sysTypes.User)
		GetTo() (bool, string, map[string]string, *sysTypes.User)
		GetCc() (bool, string, map[string]string, *sysTypes.User)
		GetHtml() (bool, string, io.Reader)
		GetPlain() (bool, string, io.Reader)
	}
)

func EmailHandler(reg emailHandlerRegistry) *emailHandler {
	h := &emailHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h emailHandler) send(_ context.Context, args *emailSendArgs) (err error) {
	msg := mail.New()

	_, s, r := args.GetSubject()
	if r != nil {
		aux, _ := ioutil.ReadAll(r)
		s = string(aux)
	}

	if err = h.procArgs(msg, s, args); err != nil {
		return
	}

	return mail.Send(msg)
}

func (h emailHandler) message(_ context.Context, args *emailMessageArgs) (*emailMessageResults, error) {
	var err error
	msg := mail.New()

	_, s, r := args.GetSubject()
	if r != nil {
		aux, _ := ioutil.ReadAll(r)
		s = string(aux)
	}

	if err = h.procArgs(msg, s, args); err != nil {
		return nil, err
	}

	return &emailMessageResults{Message: &emailMessage{msg: msg}}, nil
}

func (h emailHandler) sendMessage(_ context.Context, args *emailSendMessageArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	return mail.Send(args.Message.msg)
}

func (h emailHandler) setSubject(_ context.Context, args *emailSetSubjectArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	args.Message.msg.SetHeader("Subject", args.Subject)
	return nil
}

func (h emailHandler) setHeaders(_ context.Context, args *emailSetHeadersArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	args.Message.msg.SetHeaders(args.Headers)
	return nil
}

func (h emailHandler) setHeader(_ context.Context, args *emailSetHeaderArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	if args.hasValue {
		args.Message.msg.SetHeader(args.Name, append(args.Message.msg.GetHeader(args.Name), args.Value)...)
	} else {
		args.Message.msg.SetHeader(args.Name)
	}
	return nil
}

func (h emailHandler) setAddress(_ context.Context, args *emailSetAddressArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	args.Message.msg.SetAddressHeader(args.Type, args.Address, args.Name)
	return nil
}

func (h emailHandler) attach(_ context.Context, args *emailAttachArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	var (
		r = args.contentStream
	)

	if r == nil {
		r = strings.NewReader(args.contentString)
	}

	args.Message.msg.AttachReader(args.Name, r)
	return
}

func (h emailHandler) embed(_ context.Context, args *emailEmbedArgs) (err error) {
	if args.Message.msg == nil {
		return fmt.Errorf("email message not initialized")
	}

	args.Message.msg.EmbedReader(args.Name, args.Content)
	return
}

func (h emailHandler) procArgs(msg *gomail.Message, subject string, args messageArgs) (err error) {
	msg.SetHeader("Subject", subject)

	var (
		// if text/plain exists,
		// we'll add text/html as alternative body!
		hasPlain bool
	)

	if has, s, r := args.GetPlain(); has {
		if r != nil {
			aux, _ := ioutil.ReadAll(r)
			s = string(aux)
		}

		msg.SetBody("text/plain", s)
		hasPlain = true
	}

	if has, s, r := args.GetHtml(); has {
		if r != nil {
			aux, _ := ioutil.ReadAll(r)
			s = string(aux)
		}

		if hasPlain {
			msg.AddAlternative("text/html", s)
		} else {
			msg.SetBody("text/html", s)
		}
	}

	if has, s, m, u := args.GetTo(); has {
		if err = h.procEmailRecipients(msg, "To", s, m, u); err != nil {
			return
		}
	}

	if has, s, m, u := args.GetCc(); has {
		if err = h.procEmailRecipients(msg, "Cc", s, m, u); err != nil {
			return
		}
	}

	if has, s, u := args.GetReplyTo(); has {
		if err = h.procEmailRecipients(msg, "ReplyTo", s, nil, u); err != nil {
			return
		}
	}

	if has, s, u := args.GetFrom(); has {
		if err = h.procEmailRecipients(msg, "From", s, nil, u); err != nil {
			return
		}
	}

	return
}

func (h emailHandler) procEmailRecipients(msg *gomail.Message, field string, s string, mm map[string]string, u *sysTypes.User) (err error) {
	var (
		rr = make([]string, 0)
	)

	switch {
	case len(s) > 0:
		var (
			email string
			name  string
		)

		name = ""
		if spaceAt := strings.Index(s, " "); spaceAt > -1 {
			// proc <email> <name> ("foo@bar.baz foo baz")
			email, name = s[:spaceAt], strings.TrimSpace(s[spaceAt+1:])
		} else {
			// proc <email>
			email = s
		}

		rr = append(rr, msg.FormatAddress(email, name))

	case len(mm) > 0:
		for email, name := range mm {
			rr = append(rr, msg.FormatAddress(email, name))
		}

	case u != nil:
		rr = append(rr, msg.FormatAddress(u.Email, u.Name))

	}

	msg.SetHeader(field, rr...)
	return nil
}
