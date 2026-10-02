package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/notification_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	notificationHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h notificationHandler) register() {
	h.reg.AddFunctions(
		h.Send(),
	)
}

type (
	notificationSendArgs struct {
		hasRecipient    bool
		Recipient       interface{}
		recipientID     uint64
		recipientHandle string
		recipientEmail  string

		hasTitle bool
		Title    string

		hasDescription bool
		Description    string
	}
)

func (a notificationSendArgs) GetRecipient() (bool, uint64, string, string) {
	return a.hasRecipient, a.recipientID, a.recipientHandle, a.recipientEmail
}

// Send function Send simple notification
//
// expects implementation of send function:
//
//	func (h notificationHandler) send(ctx context.Context, args *notificationSendArgs) (err error) {
//	   return
//	}
func (h notificationHandler) Send() *atypes.Function {
	return &atypes.Function{
		Ref:    "notificationSend",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short:       "Send simple notification",
			Description: "Sends a simple notification with title and description to a user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recipient",
				Types: []string{"ID", "Handle", "String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Recipient",
				},
			},
			{
				Name:  "title",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Title",
				},
			},
			{
				Name:  "description",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Label:  "Description",
					Visual: map[string]interface{}{"input": map[string]interface{}{"type": "richtext"}},
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &notificationSendArgs{
					hasRecipient:   in.Has("recipient"),
					hasTitle:       in.Has("title"),
					hasDescription: in.Has("description"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Recipient argument
			if args.hasRecipient {
				aux := expr.Must(expr.Select(in, "recipient"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.recipientID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.recipientHandle = aux.Get().(string)
				case h.reg.Type("String").Type():
					args.recipientEmail = aux.Get().(string)
				}
			}

			return out, h.send(ctx, args)
		},
	}
}
