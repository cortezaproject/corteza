package automation

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ngTemplatesHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *templatesHandler
	}
)

func NgTemplatesHandler(reg constructSvc, tReg typeRegistry, tSvc templateService) *ngTemplatesHandler {
	ngh := &ngTemplatesHandler{
		reg:  reg,
		tReg: tReg,
		h: &templatesHandler{
			tSvc: tSvc,
		},
	}

	ngh.register()
	return ngh
}

func (h ngTemplatesHandler) register() {}

func (h ngTemplatesHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "templatesLookup",
		Kind:   "function",
		Groups: []string{"Templates"},
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Find Template",
			Description: "Look up a template by ID or handle",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "file"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Template"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "template", Types: []string{"Template"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Template)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &templatesLookupArgs{
				hasLookup: in.Has("lookup"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			var results *templatesLookupResults
			if results, err = h.h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Template").Cast(results.Template); err == nil {
				_ = expr.Assign(out, "template", tval)
			}
			return
		},
	}
}

func (h ngTemplatesHandler) Search() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "templatesSearch",
		Kind:   "function",
		Groups: []string{"Templates"},
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Search Templates",
			Description: "Search templates with filters",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "file"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "handle", Types: []string{"String"}},
			{ArgumentName: "type", Types: []string{"String"}},
			{ArgumentName: "ownerID", Types: []string{"ID"}},
			{ArgumentName: "partial", Types: []string{"Boolean"}},
			{ArgumentName: "labels", Types: []string{"LabelValue"}},
			{ArgumentName: "sort", Types: []string{"String"}},
			{ArgumentName: "limit", Types: []string{"UnsignedInteger"}},
			{ArgumentName: "incTotal", Types: []string{"Boolean"}},
			{ArgumentName: "incPageNavigation", Types: []string{"Boolean"}},
			{ArgumentName: "pageCursor", Types: []string{"String"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "templates", Types: []string{"Template"}, IsArray: true},
			{ArgumentName: "total", Types: []string{"UnsignedInteger"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Handle", Argument: "handle"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Type", Argument: "type"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Owner ID", Argument: "ownerID"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Partial", Argument: "partial"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Sort", Argument: "sort"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Limit", Argument: "limit"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Include Total", Argument: "incTotal"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Include Page Navigation", Argument: "incPageNavigation"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Page Cursor", Argument: "pageCursor"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &templatesSearchArgs{
				hasHandle:            in.Has("handle"),
				hasType:              in.Has("type"),
				hasOwnerID:           in.Has("ownerID"),
				hasPartial:           in.Has("partial"),
				hasLabels:            in.Has("labels"),
				hasSort:              in.Has("sort"),
				hasLimit:             in.Has("limit"),
				hasIncTotal:          in.Has("incTotal"),
				hasIncPageNavigation: in.Has("incPageNavigation"),
				hasPageCursor:        in.Has("pageCursor"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *templatesSearchResults
			if results, err = h.h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			var tarr = make([]expr.TypedValue, len(results.Templates))
			for i := range results.Templates {
				if tarr[i], err = h.tReg.Type("Template").Cast(results.Templates[i]); err != nil {
					return
				}
			}
			if tval, err := expr.NewArray(tarr); err == nil {
				_ = expr.Assign(out, "templates", tval)
			}

			if tval, err := h.tReg.Type("UnsignedInteger").Cast(results.Total); err == nil {
				_ = expr.Assign(out, "total", tval)
			}
			return
		},
	}
}

func (h ngTemplatesHandler) Create() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "templatesCreate",
		Kind:   "function",
		Groups: []string{"Templates"},
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Create Template",
			Description: "Create a new template",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "file"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "handle", Types: []string{"String"}, Required: true},
			{ArgumentName: "language", Types: []string{"String"}},
			{ArgumentName: "type", Types: []string{"String"}},
			{ArgumentName: "partial", Types: []string{"Boolean"}},
			{ArgumentName: "template", Types: []string{"String"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "template", Types: []string{"Template"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Handle", Argument: "handle"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Language", Argument: "language"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Type", Argument: "type"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Partial", Argument: "partial"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Template", Argument: "template"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &templatesCreateArgs{
				hasTemplate: true,
				Template:    &types.Template{},
			}

			if in.Has("handle") {
				aux := expr.Must(expr.Select(in, "handle"))
				args.Template.Handle, _ = expr.CastToString(aux.Get())
			}
			if in.Has("language") {
				aux := expr.Must(expr.Select(in, "language"))
				args.Template.Language, _ = expr.CastToString(aux.Get())
			}
			if in.Has("type") {
				aux := expr.Must(expr.Select(in, "type"))
				if typeStr, castErr := expr.CastToString(aux.Get()); castErr == nil {
					args.Template.Type = types.DocumentType(typeStr)
				}
			}
			if in.Has("partial") {
				aux := expr.Must(expr.Select(in, "partial"))
				if pBool, castErr := expr.CastToBoolean(aux.Get()); castErr == nil {
					args.Template.Partial = pBool
				}
			}
			if in.Has("template") {
				aux := expr.Must(expr.Select(in, "template"))
				args.Template.Template, _ = expr.CastToString(aux.Get())
			}

			var results *templatesCreateResults
			if results, err = h.h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Template").Cast(results.Template); err == nil {
				_ = expr.Assign(out, "template", tval)
			}
			return
		},
	}
}

func (h ngTemplatesHandler) Update() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "templatesUpdate",
		Kind:   "function",
		Groups: []string{"Templates"},
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Update Template",
			Description: "Update an existing template",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "file"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "template", Types: []string{"Template"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "template", Types: []string{"Template"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Template", Argument: "template"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &templatesUpdateArgs{
				hasTemplate: in.Has("template"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *templatesUpdateResults
			if results, err = h.h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Template").Cast(results.Template); err == nil {
				_ = expr.Assign(out, "template", tval)
			}
			return
		},
	}
}

func (h ngTemplatesHandler) Delete() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "templatesDelete",
		Kind:   "function",
		Groups: []string{"Templates"},
		Labels: map[string]string{"delete": "step", "templates": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Delete Template",
			Description: "Delete a template",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "file"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Template"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Template)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &templatesDeleteArgs{
				hasLookup: in.Has("lookup"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			return out, h.h.delete(ctx, args)
		},
	}
}

func (h ngTemplatesHandler) Render() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "templatesRender",
		Kind:   "function",
		Groups: []string{"Templates"},
		Labels: map[string]string{"render": "step", "templates": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Render Template",
			Description: "Render a template with given variables",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "file"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Template"}, Required: true},
			{ArgumentName: "documentName", Types: []string{"String"}},
			{ArgumentName: "documentType", Types: []string{"String"}},
			{ArgumentName: "variables", Types: []string{"Vars"}},
			{ArgumentName: "options", Types: []string{"RenderOptions"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "document", Types: []string{"RenderedDocument"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Template)", Argument: "lookup"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Document Name", Argument: "documentName"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Document Type", Argument: "documentType"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Variables", Argument: "variables"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Options", Argument: "options"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &templatesRenderArgs{
				hasLookup:       in.Has("lookup"),
				hasDocumentName: in.Has("documentName"),
				hasDocumentType: in.Has("documentType"),
				hasVariables:    in.Has("variables"),
				hasOptions:      in.Has("options"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			var results *templatesRenderResults
			if results, err = h.h.render(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("RenderedDocument").Cast(results.Document); err == nil {
				_ = expr.Assign(out, "document", tval)
			}
			return
		},
	}
}
