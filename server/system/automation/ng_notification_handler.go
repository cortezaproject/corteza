package automation

import (
	"context"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"go.uber.org/zap"
)

type (
	constructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
		AddTriggers(tt ...atypes.ConstructTrigger)
	}

	ngNotificationHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *notificationHandler
	}
)

func NgNotificationHandler(reg constructSvc, tReg typeRegistry, ntfSvc notificationService, uSvc notificationUserService, log *zap.Logger) *ngNotificationHandler {
	ngh := &ngNotificationHandler{
		reg:  reg,
		tReg: tReg,
		h: &notificationHandler{
			svc:    ntfSvc,
			uSvc:   uSvc,
			logger: log.Named("notification"),
		},
	}

	ngh.register()
	return ngh
}

func (h ngNotificationHandler) register() {
	h.reg.AddFunctions(
		h.Send(),
	)
}

func (h ngNotificationHandler) Send() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "notificationSend",
		Kind:   "function",
		Groups: []string{"Notifications"},
		Labels: map[string]string(nil),
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Send Notification",
			Description: "Sends a simple notification with title and description to a user",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "bell"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "recipient",
				Types:        []string{"ID", "Handle", "String"},
				Required:     true,
				Meta: &atypes.ParamMeta{
					Label:       "Recipient",
					Description: "User to send the notification to. Can be an ID, handle, or email address.",
				},
			},
			{
				ArgumentName: "title",
				Types:        []string{"String"},
				Required:     true,
			},
			{
				ArgumentName: "description",
				Types:        []string{"String"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "UserSelector",
						Label:    "Recipient",
						Argument: "recipient",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "String",
						Label:    "Title",
						Argument: "title",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "String",
						Label:    "Description",
						Argument: "description",
					},
				}},
			}},
		}},

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
				case h.tReg.Type("ID").Type():
					args.recipientID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.recipientHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.recipientEmail = aux.Get().(string)
				}
			}

			return out, h.h.send(ctx, args)
		},
	}
}
