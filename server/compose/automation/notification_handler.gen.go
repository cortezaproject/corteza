package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/notification_handler.yaml

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
		h.SendRecord(),
	)
}

type (
	notificationSendRecordArgs struct {
		hasRecipient    bool
		Recipient       interface{}
		recipientID     uint64
		recipientHandle string
		recipientEmail  string

		hasTitle bool
		Title    string

		hasDescription bool
		Description    string

		hasModule    bool
		Module       interface{}
		moduleID     uint64
		moduleHandle string

		hasNamespace    bool
		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string

		hasRecord bool
		Record    uint64

		hasOpenMode bool
		OpenMode    string

		hasEdit bool
		Edit    bool
	}
)

func (a notificationSendRecordArgs) GetRecipient() (bool, uint64, string, string) {
	return a.hasRecipient, a.recipientID, a.recipientHandle, a.recipientEmail
}

func (a notificationSendRecordArgs) GetModule() (bool, uint64, string) {
	return a.hasModule, a.moduleID, a.moduleHandle
}

func (a notificationSendRecordArgs) GetNamespace() (bool, uint64, string) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle
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
				Types: []string{"ID", "Handle", "String"}, Required: true,
			},
			{
				Name:  "title",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "description",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Visual: map[string]interface{}{"input": map[string]interface{}{"type": "richtext"}},
				},
			},
			{
				Name:  "module",
				Types: []string{"ID", "Handle"}, Required: true,
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle"}, Required: true,
			},
			{
				Name:  "record",
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
					hasModule:      in.Has("module"),
					hasNamespace:   in.Has("namespace"),
					hasRecord:      in.Has("record"),
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
				}
			}

			// Converting Module argument
			if args.hasModule {
				aux := expr.Must(expr.Select(in, "module"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.moduleID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.moduleHandle = aux.Get().(string)
				}
			}

			// Converting Namespace argument
			if args.hasNamespace {
				aux := expr.Must(expr.Select(in, "namespace"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.namespaceID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.namespaceHandle = aux.Get().(string)
				}
			}

			return out, h.sendRecord(ctx, args)
		},
	}
}
