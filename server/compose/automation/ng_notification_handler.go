package automation

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"go.uber.org/zap"
)

type (
	ngNotificationHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *notificationHandler
	}
)

func NgNotificationHandler(reg constructSvc, tReg typeRegistry, ntfSvc notificationService, uSvc notificationUserService, nSvc notificationNamespaceService, mSvc notificationModuleService, log *zap.Logger) *ngNotificationHandler {
	ngh := &ngNotificationHandler{
		reg:  reg,
		tReg: tReg,
		h: &notificationHandler{
			svc:    ntfSvc,
			uSvc:   uSvc,
			nSvc:   nSvc,
			mSvc:   mSvc,
			logger: log.Named("notification"),
		},
	}

	ngh.register()
	return ngh
}

func (h ngNotificationHandler) register() {
	h.reg.AddFunctions(
		h.SendRecord(),
	)
}

func (h ngNotificationHandler) SendRecord() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "notificationSendRecord",
		Kind:   "function",
		Groups: []string{"Compose Notification"},
		Labels: map[string]string(nil),
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Send record notification",
			Description: "Sends a notification that links to a specific record",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "bell"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "recipient", Types: []string{"ID", "Handle", "String"}, Required: true},
			{ArgumentName: "title", Types: []string{"String"}, Required: true},
			{ArgumentName: "description", Types: []string{"String"}},
			{ArgumentName: "module", Types: []string{"ID", "Handle"}, Required: true},
			{ArgumentName: "namespace", Types: []string{"ID", "Handle"}, Required: true},
			{ArgumentName: "record", Types: []string{"ID"}},
			{ArgumentName: "openMode", Types: []string{"String"}},
			{ArgumentName: "edit", Types: []string{"Boolean"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Recipient", Argument: "recipient"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Title", Argument: "title"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Description", Argument: "description"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Module", Argument: "module"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Namespace", Argument: "namespace"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Record ID", Argument: "record"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Open Mode", Argument: "openMode"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Edit", Argument: "edit"}},
				},
			}},
		}},

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
				case h.tReg.Type("ID").Type():
					args.recipientID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.recipientHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.recipientEmail = aux.Get().(string)
				}
			}

			// Converting Module argument
			if args.hasModule {
				aux := expr.Must(expr.Select(in, "module"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.moduleID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.moduleHandle = aux.Get().(string)
				}
			}

			// Converting Namespace argument
			if args.hasNamespace {
				aux := expr.Must(expr.Select(in, "namespace"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.namespaceID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.namespaceHandle = aux.Get().(string)
				}
			}

			return out, h.h.sendRecord(ctx, args)
		},
	}
}
