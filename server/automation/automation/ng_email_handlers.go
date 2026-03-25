package automation

import (
	"context"
	"io"
	"strings"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/mail"
	gomail "gopkg.in/mail.v2"
)

type (
	ngEmailConstructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
	}

	ngEmailHandler struct {
		reg ngEmailConstructSvc
	}
)

func NgEmailHandler(reg ngEmailConstructSvc) *ngEmailHandler {
	h := &ngEmailHandler{reg: reg}
	h.register()
	return h
}

func (h ngEmailHandler) register() {
	h.reg.AddFunctions(
		h.SendEmail(),
	)
}

func (h ngEmailHandler) SendEmail() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "ngEmailSend",
		Kind:   "function",
		Groups: []string{"Email"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Send Email",
			Description: "Send and configure an email message. Combines functionality of Message, Setting properties, and Sending.",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "envelope"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "message", Types: []string{"EmailMessage"}, Meta: &atypes.ParamMeta{Label: "Base Email Message (Optional)"}},
			{ArgumentName: "to", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "To Address"}},
			{ArgumentName: "cc", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "CC Address"}},
			{ArgumentName: "bcc", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "BCC Address"}},
			{ArgumentName: "replyTo", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "Reply-To Address"}},
			{ArgumentName: "from", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "From Address"}},
			{ArgumentName: "subject", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "Subject"}},
			{ArgumentName: "html", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "HTML Body"}},
			{ArgumentName: "plain", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "Plain Body"}},
			{ArgumentName: "headers", Types: []string{"KVV"}, Meta: &atypes.ParamMeta{Label: "Custom Headers (KVV)"}},
			{ArgumentName: "attach", Types: []string{"Reader", "String"}, Meta: &atypes.ParamMeta{Label: "Attachment Content"}},
			{ArgumentName: "attachName", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "Attachment Name"}},
			{ArgumentName: "embed", Types: []string{"Reader"}, Meta: &atypes.ParamMeta{Label: "Embed Content"}},
			{ArgumentName: "embedName", Types: []string{"String"}, Meta: &atypes.ParamMeta{Label: "Embed Name"}},
		},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Base Message", Argument: "message"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "To", Argument: "to"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "CC", Argument: "cc"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "BCC", Argument: "bcc"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Reply-To", Argument: "replyTo"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "From", Argument: "from"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Subject", Argument: "subject"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Custom Headers (KVV)", Argument: "headers"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "HTML Body", Argument: "html"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Plain Body", Argument: "plain"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Attachment Content", Argument: "attach"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Attachment Name", Argument: "attachName"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Embed Content", Argument: "embed"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Embed Name", Argument: "embedName"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			out = &expr.Vars{}

			var msg *gomail.Message

			if in != nil && in.Has("message") {
				v := expr.Must(expr.Select(in, "message"))
				if em, ok := v.Get().(*emailMessage); ok && em != nil && em.msg != nil {
					msg = em.msg
				}
			}

			if msg == nil {
				msg = mail.New()
			}

			if in != nil {
				if in.Has("headers") {
					v := expr.Must(expr.Select(in, "headers"))
					if hdrs, ok := v.Get().(map[string][]string); ok {
						msg.SetHeaders(hdrs)
					}
				}

				if in.Has("to") {
					v := expr.Must(expr.Select(in, "to"))
					if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.SetHeader("To", s)
					}
				}

				if in.Has("cc") {
					v := expr.Must(expr.Select(in, "cc"))
					if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.SetHeader("Cc", s)
					}
				}

				if in.Has("bcc") {
					v := expr.Must(expr.Select(in, "bcc"))
					if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.SetHeader("Bcc", s)
					}
				}

				if in.Has("replyTo") {
					v := expr.Must(expr.Select(in, "replyTo"))
					if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.SetHeader("Reply-To", s)
					}
				}

				if in.Has("from") {
					v := expr.Must(expr.Select(in, "from"))
					if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.SetHeader("From", s)
					}
				}

				if in.Has("subject") {
					v := expr.Must(expr.Select(in, "subject"))
					if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.SetHeader("Subject", s)
					}
				}

				var htmlStr, plainStr string
				if in.Has("html") {
					v := expr.Must(expr.Select(in, "html"))
					htmlStr, _ = expr.CastToString(v.Get())
				}
				if in.Has("plain") {
					v := expr.Must(expr.Select(in, "plain"))
					plainStr, _ = expr.CastToString(v.Get())
				}

				if plainStr != "" {
					msg.SetBody("text/plain", plainStr)
				}
				if htmlStr != "" {
					if plainStr != "" {
						msg.AddAlternative("text/html", htmlStr)
					} else {
						msg.SetBody("text/html", htmlStr)
					}
				}

				var attachName, embedName string
				if in.Has("attachName") {
					v := expr.Must(expr.Select(in, "attachName"))
					attachName, _ = expr.CastToString(v.Get())
				}
				if in.Has("embedName") {
					v := expr.Must(expr.Select(in, "embedName"))
					embedName, _ = expr.CastToString(v.Get())
				}

				if in.Has("attach") {
					v := expr.Must(expr.Select(in, "attach"))
					if r, ok := v.Get().(io.Reader); ok {
						msg.AttachReader(attachName, r)
					} else if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.AttachReader(attachName, strings.NewReader(s))
					}
				}

				if in.Has("embed") {
					v := expr.Must(expr.Select(in, "embed"))
					if r, ok := v.Get().(io.Reader); ok {
						msg.EmbedReader(embedName, r)
					} else if s, err := expr.CastToString(v.Get()); err == nil && s != "" {
						msg.EmbedReader(embedName, strings.NewReader(s))
					}
				}
			}

			// Assign the configured message to output just in case someone wants to chain it
			out.Set("message", &emailMessage{msg: msg})

			return out, mail.Send(msg)
		},
	}
}
