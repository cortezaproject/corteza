package automation

import (
	"context"
	"io"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	ngAttachmentHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *attachmentHandler
	}
)

func NgAttachmentHandler(reg constructSvc, tReg typeRegistry, svc attachmentService) *ngAttachmentHandler {
	ngh := &ngAttachmentHandler{
		reg:  reg,
		tReg: tReg,
		h:    &attachmentHandler{svc: svc}, // note: we don't pass registry so it doesn't have it, but internal methods only use svc
	}

	ngh.register()
	return ngh
}

func (h ngAttachmentHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Create(),
		h.Delete(),
		h.OpenOriginal(),
		h.OpenPreview(),
	)
}

func (h ngAttachmentHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "attachmentLookup",
		Kind:   "function",
		Groups: []string{"Compose Attachment"},
		Labels: map[string]string{"attachment": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Attachment lookup",
			Description: "Find specific attachment by ID",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "paperclip"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "attachment",
				Name:         "",
				Types:        []string{"ID"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "attachment",
				Name:         "",
				Types:        []string{"Attachment"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Attachment",
						Argument: "attachment",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentLookupArgs{
					hasAttachment: in.Has("attachment"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *attachmentLookupResults
			if results, err = h.h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Attachment (*types.Attachment) to Attachment
				var (
					tval expr.TypedValue
				)

				if tval, err = h.tReg.Type("Attachment").Cast(results.Attachment); err != nil {
					return
				} else if err = expr.Assign(out, "attachment", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

func (h ngAttachmentHandler) Create() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "attachmentCreate",
		Kind:   "function",
		Groups: []string{"Compose Attachment"},
		Labels: map[string]string{"attachment": "step,workflow", "create": "step"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Create file and attach it to a resource",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "paperclip"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "name",
				Name:         "",
				Types:        []string{"String"},
			},
			{
				ArgumentName: "resource",
				Name:         "",
				Types:        []string{"ComposeRecord"}, Required: true,
			},
			{
				ArgumentName: "fieldName",
				Name:         "",
				Types:        []string{"String"},
			},
			{
				ArgumentName: "content",
				Name:         "",
				Types:        []string{"String", "Reader", "Bytes"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "attachment",
				Name:         "",
				Types:        []string{"Attachment"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Name",
						Argument: "name",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Resource",
						Argument: "resource",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Field Name",
						Argument: "fieldName",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Content",
						Argument: "content",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentCreateArgs{
					hasName:      in.Has("name"),
					hasResource:  in.Has("resource"),
					hasFieldName: in.Has("fieldName"),
					hasContent:   in.Has("content"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Content argument
			if args.hasContent {
				aux := expr.Must(expr.Select(in, "content"))
				switch aux.Type() {
				case h.tReg.Type("String").Type():
					args.contentString = aux.Get().(string)
				case h.tReg.Type("Reader").Type():
					args.contentStream = aux.Get().(io.Reader)
				case h.tReg.Type("Bytes").Type():
					args.contentBytes = aux.Get().([]byte)
				}
			}

			var results *attachmentCreateResults
			if results, err = h.h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Attachment (*types.Attachment) to Attachment
				var (
					tval expr.TypedValue
				)

				if tval, err = h.tReg.Type("Attachment").Cast(results.Attachment); err != nil {
					return
				} else if err = expr.Assign(out, "attachment", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

func (h ngAttachmentHandler) Delete() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "attachmentDelete",
		Kind:   "function",
		Groups: []string{"Compose Attachment"},
		Labels: map[string]string{"attachment": "step,workflow", "delete": "step"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Delete attachment",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "paperclip"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "attachment",
				Name:         "",
				Types:        []string{"ID"}, Required: true,
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Attachment",
						Argument: "attachment",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentDeleteArgs{
					hasAttachment: in.Has("attachment"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.h.delete(ctx, args)
		},
	}
}

func (h ngAttachmentHandler) OpenOriginal() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "attachmentOpenOriginal",
		Kind:   "function",
		Groups: []string{"Compose Attachment"},
		Labels: map[string]string{"attachment": "step,workflow", "original-attachment": "step"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Open original attachment",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "paperclip"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "attachment",
				Name:         "",
				Types:        []string{"ID", "Attachment"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "content",
				Name:         "",
				Types:        []string{"Reader"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Attachment",
						Argument: "attachment",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentOpenOriginalArgs{
					hasAttachment: in.Has("attachment"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Attachment argument
			if args.hasAttachment {
				aux := expr.Must(expr.Select(in, "attachment"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.attachmentID = aux.Get().(uint64)
				case h.tReg.Type("Attachment").Type():
					args.attachmentAttachment = aux.Get().(*types.Attachment)
				}
			}

			var results *attachmentOpenOriginalResults
			if results, err = h.h.openOriginal(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Content (io.Reader) to Reader
				var (
					tval expr.TypedValue
				)

				if tval, err = h.tReg.Type("Reader").Cast(results.Content); err != nil {
					return
				} else if err = expr.Assign(out, "content", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

func (h ngAttachmentHandler) OpenPreview() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "attachmentOpenPreview",
		Kind:   "function",
		Groups: []string{"Compose Attachment"},
		Labels: map[string]string{"attachment": "step,workflow", "preview-attachment": "step"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Open attachment preview",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "paperclip"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "attachment",
				Name:         "",
				Types:        []string{"ID", "Attachment"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "content",
				Name:         "",
				Types:        []string{"Reader"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "Expression",
						Label:    "Attachment",
						Argument: "attachment",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentOpenPreviewArgs{
					hasAttachment: in.Has("attachment"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Attachment argument
			if args.hasAttachment {
				aux := expr.Must(expr.Select(in, "attachment"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.attachmentID = aux.Get().(uint64)
				case h.tReg.Type("Attachment").Type():
					args.attachmentAttachment = aux.Get().(*types.Attachment)
				}
			}

			var results *attachmentOpenPreviewResults
			if results, err = h.h.openPreview(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Content (io.Reader) to Reader
				var (
					tval expr.TypedValue
				)

				if tval, err = h.tReg.Type("Reader").Cast(results.Content); err != nil {
					return
				} else if err = expr.Assign(out, "content", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
