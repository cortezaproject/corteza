package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/attachment_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"io"
)

var _ wfexec.ExecResponse

type (
	attachmentHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h attachmentHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Create(),
		h.Delete(),
		h.OpenOriginal(),
		h.OpenPreview(),
	)
}

type (
	attachmentLookupArgs struct {
		hasAttachment bool
		Attachment    uint64
	}

	attachmentLookupResults struct {
		Attachment *types.Attachment
	}
)

// Lookup function Attachment lookup
//
// expects implementation of lookup function:
// func (h attachmentHandler) lookup(ctx context.Context, args *attachmentLookupArgs) (results *attachmentLookupResults, err error) {
//    return
// }
func (h attachmentHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "attachmentLookup",
		Kind:   "function",
		Labels: map[string]string{"attachment": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Attachment lookup",
			Description: "Find specific attachment by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "attachment",
				Types: []string{"ID"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "attachment",
				Types: []string{"Attachment"},
			},
		},

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
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Attachment (*types.Attachment) to Attachment
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Attachment").Cast(results.Attachment); err != nil {
					return
				} else if err = expr.Assign(out, "attachment", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	attachmentCreateArgs struct {
		hasName bool
		Name    string

		hasResource bool
		Resource    *types.Record

		hasContent    bool
		Content       interface{}
		contentString string
		contentStream io.Reader
		contentBytes  []byte
	}

	attachmentCreateResults struct {
		Attachment *types.Attachment
	}
)

func (a attachmentCreateArgs) GetContent() (bool, string, io.Reader, []byte) {
	return a.hasContent, a.contentString, a.contentStream, a.contentBytes
}

// Create function Create file and attach it to a resource
//
// expects implementation of create function:
// func (h attachmentHandler) create(ctx context.Context, args *attachmentCreateArgs) (results *attachmentCreateResults, err error) {
//    return
// }
func (h attachmentHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:    "attachmentCreate",
		Kind:   "function",
		Labels: map[string]string{"attachment": "step,workflow", "create": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Create file and attach it to a resource",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "name",
				Types: []string{"String"},
			},
			{
				Name:  "resource",
				Types: []string{"ComposeRecord"}, Required: true,
			},
			{
				Name:  "content",
				Types: []string{"String", "Reader", "Bytes"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "attachment",
				Types: []string{"Attachment"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentCreateArgs{
					hasName:     in.Has("name"),
					hasResource: in.Has("resource"),
					hasContent:  in.Has("content"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Content argument
			if args.hasContent {
				aux := expr.Must(expr.Select(in, "content"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.contentString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.contentStream = aux.Get().(io.Reader)
				case h.reg.Type("Bytes").Type():
					args.contentBytes = aux.Get().([]byte)
				}
			}

			var results *attachmentCreateResults
			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Attachment (*types.Attachment) to Attachment
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Attachment").Cast(results.Attachment); err != nil {
					return
				} else if err = expr.Assign(out, "attachment", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	attachmentDeleteArgs struct {
		hasAttachment bool
		Attachment    uint64
	}
)

// Delete function Delete attachment
//
// expects implementation of delete function:
// func (h attachmentHandler) delete(ctx context.Context, args *attachmentDeleteArgs) (err error) {
//    return
// }
func (h attachmentHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:    "attachmentDelete",
		Kind:   "function",
		Labels: map[string]string{"attachment": "step,workflow", "delete": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Delete attachment",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "attachment",
				Types: []string{"ID"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &attachmentDeleteArgs{
					hasAttachment: in.Has("attachment"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.delete(ctx, args)
		},
	}
}

type (
	attachmentOpenOriginalArgs struct {
		hasAttachment        bool
		Attachment           interface{}
		attachmentID         uint64
		attachmentAttachment *types.Attachment
	}

	attachmentOpenOriginalResults struct {
		Content io.Reader
	}
)

func (a attachmentOpenOriginalArgs) GetAttachment() (bool, uint64, *types.Attachment) {
	return a.hasAttachment, a.attachmentID, a.attachmentAttachment
}

// OpenOriginal function Open original attachment
//
// expects implementation of openOriginal function:
// func (h attachmentHandler) openOriginal(ctx context.Context, args *attachmentOpenOriginalArgs) (results *attachmentOpenOriginalResults, err error) {
//    return
// }
func (h attachmentHandler) OpenOriginal() *atypes.Function {
	return &atypes.Function{
		Ref:    "attachmentOpenOriginal",
		Kind:   "function",
		Labels: map[string]string{"attachment": "step,workflow", "original-attachment": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Open original attachment",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "attachment",
				Types: []string{"ID", "Attachment"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "content",
				Types: []string{"Reader"},
			},
		},

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
				case h.reg.Type("ID").Type():
					args.attachmentID = aux.Get().(uint64)
				case h.reg.Type("Attachment").Type():
					args.attachmentAttachment = aux.Get().(*types.Attachment)
				}
			}

			var results *attachmentOpenOriginalResults
			if results, err = h.openOriginal(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Content (io.Reader) to Reader
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Reader").Cast(results.Content); err != nil {
					return
				} else if err = expr.Assign(out, "content", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	attachmentOpenPreviewArgs struct {
		hasAttachment        bool
		Attachment           interface{}
		attachmentID         uint64
		attachmentAttachment *types.Attachment
	}

	attachmentOpenPreviewResults struct {
		Content io.Reader
	}
)

func (a attachmentOpenPreviewArgs) GetAttachment() (bool, uint64, *types.Attachment) {
	return a.hasAttachment, a.attachmentID, a.attachmentAttachment
}

// OpenPreview function Open attachment preview
//
// expects implementation of openPreview function:
// func (h attachmentHandler) openPreview(ctx context.Context, args *attachmentOpenPreviewArgs) (results *attachmentOpenPreviewResults, err error) {
//    return
// }
func (h attachmentHandler) OpenPreview() *atypes.Function {
	return &atypes.Function{
		Ref:    "attachmentOpenPreview",
		Kind:   "function",
		Labels: map[string]string{"attachment": "step,workflow", "preview-attachment": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Open attachment preview",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "attachment",
				Types: []string{"ID", "Attachment"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "content",
				Types: []string{"Reader"},
			},
		},

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
				case h.reg.Type("ID").Type():
					args.attachmentID = aux.Get().(uint64)
				case h.reg.Type("Attachment").Type():
					args.attachmentAttachment = aux.Get().(*types.Attachment)
				}
			}

			var results *attachmentOpenPreviewResults
			if results, err = h.openPreview(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Content (io.Reader) to Reader
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Reader").Cast(results.Content); err != nil {
					return
				} else if err = expr.Assign(out, "content", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
