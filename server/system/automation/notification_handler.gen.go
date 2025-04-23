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
	"github.com/cortezaproject/corteza/server/system/types"
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
		h.SendRecord(),
	)
}

type (
	notificationSendArgs struct {
		hasRecipient    bool
		Recipient       interface{}
		recipientID     uint64
		recipientHandle string
		recipientEmail  string
		recipientRes    *types.User

		hasTitle bool
		Title    string

		hasDescription bool
		Description    string
	}
)

func (a notificationSendArgs) GetRecipient() (bool, uint64, string, string, *types.User) {
	return a.hasRecipient, a.recipientID, a.recipientHandle, a.recipientEmail, a.recipientRes
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
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
			{
				Name:  "title",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "description",
				Types: []string{"String"},
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
				case h.reg.Type("User").Type():
					args.recipientRes = aux.Get().(*types.User)
				}
			}

			return out, h.send(ctx, args)
		},
	}
}

type (
	notificationSendRecordArgs struct {
		hasRecipient    bool
		Recipient       interface{}
		recipientID     uint64
		recipientHandle string
		recipientEmail  string
		recipientRes    *types.User

		hasTitle bool
		Title    string

		hasDescription bool
		Description    string

		hasModuleID bool
		ModuleID    uint64

		hasNamespaceID bool
		NamespaceID    uint64

		hasRecordID bool
		RecordID    uint64

		hasOpenMode bool
		OpenMode    string

		hasEdit bool
		Edit    bool
	}
)

func (a notificationSendRecordArgs) GetRecipient() (bool, uint64, string, string, *types.User) {
	return a.hasRecipient, a.recipientID, a.recipientHandle, a.recipientEmail, a.recipientRes
}

// SendRecord function Send record notification
//
// expects implementation of sendRecord function:
//
//	func (h notificationHandler) sendRecord(ctx context.Context, args *notificationSendRecordArgs) (err error) {
//	   return
//	}
func (h notificationHandler) SendRecord() *atypes.Function {
	return &atypes.Function{
		Ref:    "notificationSendRecord",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short:       "Send record notification",
			Description: "Sends a notification that links to a specific record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recipient",
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
			{
				Name:  "title",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "description",
				Types: []string{"String"},
			},
			{
				Name:  "moduleID",
				Types: []string{"ID"}, Required: true,
			},
			{
				Name:  "namespaceID",
				Types: []string{"ID"}, Required: true,
			},
			{
				Name:  "recordID",
				Types: []string{"ID"},
			},
			{
				Name:  "openMode",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Visual: map[string]interface{}{"input": map[string]interface{}{"default": "modal", "properties": map[string]interface{}{"options": []interface{}{map[string]interface{}{"text": "Open link in the same tab", "value": "sameTab"}, map[string]interface{}{"text": "Open link in a new tab", "value": "newTab"}, map[string]interface{}{"text": "Open in a modal", "value": "modal"}}}, "type": "select"}},
				},
			},
			{
				Name:  "edit",
				Types: []string{"Boolean"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &notificationSendRecordArgs{
					hasRecipient:   in.Has("recipient"),
					hasTitle:       in.Has("title"),
					hasDescription: in.Has("description"),
					hasModuleID:    in.Has("moduleID"),
					hasNamespaceID: in.Has("namespaceID"),
					hasRecordID:    in.Has("recordID"),
					hasOpenMode:    in.Has("openMode"),
					hasEdit:        in.Has("edit"),
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
				case h.reg.Type("User").Type():
					args.recipientRes = aux.Get().(*types.User)
				}
			}

			return out, h.sendRecord(ctx, args)
		},
	}
}
