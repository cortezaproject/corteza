package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/email_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"io"
)

var _ wfexec.ExecResponse

type (
	emailHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h emailHandler) register() {
	h.reg.AddFunctions(
		h.Send(),
		h.Message(),
		h.SendMessage(),
		h.SetSubject(),
		h.SetHeaders(),
		h.SetHeader(),
		h.SetAddress(),
		h.Attach(),
		h.Embed(),
	)
}

type (
	emailSendArgs struct {
		hasSubject bool
		Subject    string

		hasReplyTo    bool
		ReplyTo       interface{}
		replyToString string
		replyToUser   *sysTypes.User

		hasFrom    bool
		From       interface{}
		fromString string
		fromUser   *sysTypes.User

		hasTo    bool
		To       interface{}
		toString string
		toKV     map[string]string
		toUser   *sysTypes.User

		hasCc    bool
		Cc       interface{}
		ccString string
		ccKV     map[string]string
		ccUser   *sysTypes.User

		hasHtml    bool
		Html       interface{}
		htmlString string
		htmlStream io.Reader

		hasPlain    bool
		Plain       interface{}
		plainString string
		plainStream io.Reader
	}
)

func (a emailSendArgs) GetReplyTo() (bool, string, *sysTypes.User) {
	return a.hasReplyTo, a.replyToString, a.replyToUser
}

func (a emailSendArgs) GetFrom() (bool, string, *sysTypes.User) {
	return a.hasFrom, a.fromString, a.fromUser
}

func (a emailSendArgs) GetTo() (bool, string, map[string]string, *sysTypes.User) {
	return a.hasTo, a.toString, a.toKV, a.toUser
}

func (a emailSendArgs) GetCc() (bool, string, map[string]string, *sysTypes.User) {
	return a.hasCc, a.ccString, a.ccKV, a.ccUser
}

func (a emailSendArgs) GetHtml() (bool, string, io.Reader) {
	return a.hasHtml, a.htmlString, a.htmlStream
}

func (a emailSendArgs) GetPlain() (bool, string, io.Reader) {
	return a.hasPlain, a.plainString, a.plainStream
}

// Send function Sends email message with basic parameters
//
// expects implementation of send function:
// func (h emailHandler) send(ctx context.Context, args *emailSendArgs) (err error) {
//    return
// }
func (h emailHandler) Send() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailSend",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Sends email message with basic parameters",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "subject",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label: "Subject",
				},
			},
			{
				Name:  "replyTo",
				Types: []string{"String", "User"},
				Meta: &atypes.ParamMeta{
					Label: "Reply to",
				},
			},
			{
				Name:  "from",
				Types: []string{"String", "User"},
				Meta: &atypes.ParamMeta{
					Label: "Sender",
				},
			},
			{
				Name:  "to",
				Types: []string{"String", "KV", "User"},
				Meta: &atypes.ParamMeta{
					Label: "Recipients",
				},
			},
			{
				Name:  "cc",
				Types: []string{"String", "KV", "User"},
				Meta: &atypes.ParamMeta{
					Label: "CC",
				},
			},
			{
				Name:  "html",
				Types: []string{"String", "Reader"},
				Meta: &atypes.ParamMeta{
					Label: "HTML message body",
				},
			},
			{
				Name:  "plain",
				Types: []string{"String", "Reader"},
				Meta: &atypes.ParamMeta{
					Label: "Plain text message body",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailSendArgs{
					hasSubject: in.Has("subject"),
					hasReplyTo: in.Has("replyTo"),
					hasFrom:    in.Has("from"),
					hasTo:      in.Has("to"),
					hasCc:      in.Has("cc"),
					hasHtml:    in.Has("html"),
					hasPlain:   in.Has("plain"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting ReplyTo argument
			if args.hasReplyTo {
				aux := expr.Must(expr.Select(in, "replyTo"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.replyToString = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.replyToUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting From argument
			if args.hasFrom {
				aux := expr.Must(expr.Select(in, "from"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.fromString = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.fromUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting To argument
			if args.hasTo {
				aux := expr.Must(expr.Select(in, "to"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.toString = aux.Get().(string)
				case h.reg.Type("KV").Type():
					args.toKV = aux.Get().(map[string]string)
				case h.reg.Type("User").Type():
					args.toUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting Cc argument
			if args.hasCc {
				aux := expr.Must(expr.Select(in, "cc"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.ccString = aux.Get().(string)
				case h.reg.Type("KV").Type():
					args.ccKV = aux.Get().(map[string]string)
				case h.reg.Type("User").Type():
					args.ccUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting Html argument
			if args.hasHtml {
				aux := expr.Must(expr.Select(in, "html"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.htmlString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.htmlStream = aux.Get().(io.Reader)
				}
			}

			// Converting Plain argument
			if args.hasPlain {
				aux := expr.Must(expr.Select(in, "plain"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.plainString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.plainStream = aux.Get().(io.Reader)
				}
			}

			return out, h.send(ctx, args)
		},
	}
}

type (
	emailMessageArgs struct {
		hasSubject bool
		Subject    string

		hasReplyTo    bool
		ReplyTo       interface{}
		replyToString string
		replyToUser   *sysTypes.User

		hasFrom    bool
		From       interface{}
		fromString string
		fromUser   *sysTypes.User

		hasTo    bool
		To       interface{}
		toString string
		toKV     map[string]string
		toUser   *sysTypes.User

		hasCc    bool
		Cc       interface{}
		ccString string
		ccKV     map[string]string
		ccUser   *sysTypes.User

		hasHtml    bool
		Html       interface{}
		htmlString string
		htmlStream io.Reader

		hasPlain    bool
		Plain       interface{}
		plainString string
		plainStream io.Reader
	}

	emailMessageResults struct {
		Message *emailMessage
	}
)

func (a emailMessageArgs) GetReplyTo() (bool, string, *sysTypes.User) {
	return a.hasReplyTo, a.replyToString, a.replyToUser
}

func (a emailMessageArgs) GetFrom() (bool, string, *sysTypes.User) {
	return a.hasFrom, a.fromString, a.fromUser
}

func (a emailMessageArgs) GetTo() (bool, string, map[string]string, *sysTypes.User) {
	return a.hasTo, a.toString, a.toKV, a.toUser
}

func (a emailMessageArgs) GetCc() (bool, string, map[string]string, *sysTypes.User) {
	return a.hasCc, a.ccString, a.ccKV, a.ccUser
}

func (a emailMessageArgs) GetHtml() (bool, string, io.Reader) {
	return a.hasHtml, a.htmlString, a.htmlStream
}

func (a emailMessageArgs) GetPlain() (bool, string, io.Reader) {
	return a.hasPlain, a.plainString, a.plainStream
}

// Message function Constructs new email message
//
// expects implementation of message function:
// func (h emailHandler) message(ctx context.Context, args *emailMessageArgs) (results *emailMessageResults, err error) {
//    return
// }
func (h emailHandler) Message() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailMessage",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Constructs new email message",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "subject",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label: "Subject",
				},
			},
			{
				Name:  "replyTo",
				Types: []string{"String", "User"},
				Meta: &atypes.ParamMeta{
					Label: "Reply to",
				},
			},
			{
				Name:  "from",
				Types: []string{"String", "User"},
				Meta: &atypes.ParamMeta{
					Label: "Sender",
				},
			},
			{
				Name:  "to",
				Types: []string{"String", "KV", "User"},
				Meta: &atypes.ParamMeta{
					Label: "Recipients",
				},
			},
			{
				Name:  "cc",
				Types: []string{"String", "KV", "User"},
				Meta: &atypes.ParamMeta{
					Label: "CC",
				},
			},
			{
				Name:  "html",
				Types: []string{"String", "Reader"},
				Meta: &atypes.ParamMeta{
					Label: "HTML message body",
				},
			},
			{
				Name:  "plain",
				Types: []string{"String", "Reader"},
				Meta: &atypes.ParamMeta{
					Label: "Plain text message body",
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "message",
				Types: []string{"EmailMessage"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailMessageArgs{
					hasSubject: in.Has("subject"),
					hasReplyTo: in.Has("replyTo"),
					hasFrom:    in.Has("from"),
					hasTo:      in.Has("to"),
					hasCc:      in.Has("cc"),
					hasHtml:    in.Has("html"),
					hasPlain:   in.Has("plain"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting ReplyTo argument
			if args.hasReplyTo {
				aux := expr.Must(expr.Select(in, "replyTo"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.replyToString = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.replyToUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting From argument
			if args.hasFrom {
				aux := expr.Must(expr.Select(in, "from"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.fromString = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.fromUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting To argument
			if args.hasTo {
				aux := expr.Must(expr.Select(in, "to"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.toString = aux.Get().(string)
				case h.reg.Type("KV").Type():
					args.toKV = aux.Get().(map[string]string)
				case h.reg.Type("User").Type():
					args.toUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting Cc argument
			if args.hasCc {
				aux := expr.Must(expr.Select(in, "cc"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.ccString = aux.Get().(string)
				case h.reg.Type("KV").Type():
					args.ccKV = aux.Get().(map[string]string)
				case h.reg.Type("User").Type():
					args.ccUser = aux.Get().(*sysTypes.User)
				}
			}

			// Converting Html argument
			if args.hasHtml {
				aux := expr.Must(expr.Select(in, "html"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.htmlString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.htmlStream = aux.Get().(io.Reader)
				}
			}

			// Converting Plain argument
			if args.hasPlain {
				aux := expr.Must(expr.Select(in, "plain"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.plainString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.plainStream = aux.Get().(io.Reader)
				}
			}

			var results *emailMessageResults
			if results, err = h.message(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Message (*emailMessage) to EmailMessage
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("EmailMessage").Cast(results.Message); err != nil {
					return
				} else if err = expr.Assign(out, "message", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	emailSendMessageArgs struct {
		hasMessage bool
		Message    *emailMessage
	}
)

// SendMessage function Sends email message
//
// expects implementation of sendMessage function:
// func (h emailHandler) sendMessage(ctx context.Context, args *emailSendMessageArgs) (err error) {
//    return
// }
func (h emailHandler) SendMessage() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailSendMessage",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Sends email message",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailSendMessageArgs{
					hasMessage: in.Has("message"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.sendMessage(ctx, args)
		},
	}
}

type (
	emailSetSubjectArgs struct {
		hasMessage bool
		Message    *emailMessage

		hasSubject bool
		Subject    string
	}
)

// SetSubject function Sets message subject
//
// expects implementation of setSubject function:
// func (h emailHandler) setSubject(ctx context.Context, args *emailSetSubjectArgs) (err error) {
//    return
// }
func (h emailHandler) SetSubject() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailSetSubject",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Sets message subject",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
			{
				Name:  "subject",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Subject",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailSetSubjectArgs{
					hasMessage: in.Has("message"),
					hasSubject: in.Has("subject"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.setSubject(ctx, args)
		},
	}
}

type (
	emailSetHeadersArgs struct {
		hasMessage bool
		Message    *emailMessage

		hasHeaders bool
		Headers    map[string][]string
	}
)

// SetHeaders function Sets message headers (overrides any existing headers, subject, recipients)
//
// expects implementation of setHeaders function:
// func (h emailHandler) setHeaders(ctx context.Context, args *emailSetHeadersArgs) (err error) {
//    return
// }
func (h emailHandler) SetHeaders() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailSetHeaders",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Sets message headers (overrides any existing headers, subject, recipients)",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
			{
				Name:  "headers",
				Types: []string{"KVV"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Headers",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailSetHeadersArgs{
					hasMessage: in.Has("message"),
					hasHeaders: in.Has("headers"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.setHeaders(ctx, args)
		},
	}
}

type (
	emailSetHeaderArgs struct {
		hasMessage bool
		Message    *emailMessage

		hasName bool
		Name    string

		hasValue bool
		Value    string
	}
)

// SetHeader function Appends value or removes specific header,
//
// expects implementation of setHeader function:
// func (h emailHandler) setHeader(ctx context.Context, args *emailSetHeaderArgs) (err error) {
//    return
// }
func (h emailHandler) SetHeader() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailSetHeader",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Appends value or removes specific header,",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
			{
				Name:  "name",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Value",
				},
			},
			{
				Name:  "value",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label:       "Value",
					Description: "Raw header value. Omiting value will remove header.",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailSetHeaderArgs{
					hasMessage: in.Has("message"),
					hasName:    in.Has("name"),
					hasValue:   in.Has("value"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.setHeader(ctx, args)
		},
	}
}

type (
	emailSetAddressArgs struct {
		hasMessage bool
		Message    *emailMessage

		hasType bool
		Type    string

		hasAddress bool
		Address    string

		hasName bool
		Name    string
	}
)

// SetAddress function Adds new recipient, sender or reply-to address
//
// expects implementation of setAddress function:
// func (h emailHandler) setAddress(ctx context.Context, args *emailSetAddressArgs) (err error) {
//    return
// }
func (h emailHandler) SetAddress() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailSetAddress",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Adds new recipient, sender or reply-to address",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
			{
				Name:  "type",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Type",
					Description: "One of From",
				},
			},
			{
				Name:  "address",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Address",
				},
			},
			{
				Name:  "name",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label: "Name",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailSetAddressArgs{
					hasMessage: in.Has("message"),
					hasType:    in.Has("type"),
					hasAddress: in.Has("address"),
					hasName:    in.Has("name"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.setAddress(ctx, args)
		},
	}
}

type (
	emailAttachArgs struct {
		hasMessage bool
		Message    *emailMessage

		hasContent    bool
		Content       interface{}
		contentStream io.Reader
		contentString string

		hasName bool
		Name    string
	}
)

func (a emailAttachArgs) GetContent() (bool, io.Reader, string) {
	return a.hasContent, a.contentStream, a.contentString
}

// Attach function Attach content to an email message
//
// expects implementation of attach function:
// func (h emailHandler) attach(ctx context.Context, args *emailAttachArgs) (err error) {
//    return
// }
func (h emailHandler) Attach() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailAttach",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Attach content to an email message",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
			{
				Name:  "content",
				Types: []string{"Reader", "String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Content",
				},
			},
			{
				Name:  "name",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label: "Name",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailAttachArgs{
					hasMessage: in.Has("message"),
					hasContent: in.Has("content"),
					hasName:    in.Has("name"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Content argument
			if args.hasContent {
				aux := expr.Must(expr.Select(in, "content"))
				switch aux.Type() {
				case h.reg.Type("Reader").Type():
					args.contentStream = aux.Get().(io.Reader)
				case h.reg.Type("String").Type():
					args.contentString = aux.Get().(string)
				}
			}

			return out, h.attach(ctx, args)
		},
	}
}

type (
	emailEmbedArgs struct {
		hasMessage bool
		Message    *emailMessage

		hasContent bool
		Content    io.Reader

		hasName bool
		Name    string
	}
)

// Embed function Embed file to an email message
//
// expects implementation of embed function:
// func (h emailHandler) embed(ctx context.Context, args *emailEmbedArgs) (err error) {
//    return
// }
func (h emailHandler) Embed() *atypes.Function {
	return &atypes.Function{
		Ref:    "emailEmbed",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Embed file to an email message",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "message",
				Types: []string{"EmailMessage"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Message to be sent",
				},
			},
			{
				Name:  "content",
				Types: []string{"Reader"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Content",
				},
			},
			{
				Name:  "name",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label: "Name",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &emailEmbedArgs{
					hasMessage: in.Has("message"),
					hasContent: in.Has("content"),
					hasName:    in.Has("name"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.embed(ctx, args)
		},
	}
}
