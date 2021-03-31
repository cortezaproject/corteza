package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/templates_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ wfexec.ExecResponse

type (
	templatesHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h templatesHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Search(),
		h.Each(),
		h.Create(),
		h.Update(),
		h.Delete(),
		h.Recover(),
		h.Render(),
	)
}

type (
	templatesLookupArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Template
	}

	templatesLookupResults struct {
		Template *types.Template
	}
)

func (a templatesLookupArgs) GetLookup() (bool, uint64, string, *types.Template) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Lookup function Template lookup
//
// expects implementation of lookup function:
// func (h templatesHandler) lookup(ctx context.Context, args *templatesLookupArgs) (results *templatesLookupResults, err error) {
//    return
// }
func (h templatesHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesLookup",
		Kind:   "function",
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Template lookup",
			Description: "Find specific template by ID or handle",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Template"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "template",
				Types: []string{"Template"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesLookupArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			var results *templatesLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Template (*types.Template) to Template
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Template").Cast(results.Template); err != nil {
					return
				} else if err = expr.Assign(out, "template", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	templatesSearchArgs struct {
		hasHandle bool
		Handle    string

		hasType bool
		Type    string

		hasOwnerID bool
		OwnerID    uint64

		hasPartial bool
		Partial    bool

		hasLabels bool
		Labels    map[string]string

		hasSort bool
		Sort    string

		hasLimit bool
		Limit    uint64

		hasIncTotal bool
		IncTotal    bool

		hasIncPageNavigation bool
		IncPageNavigation    bool

		hasPageCursor bool
		PageCursor    string
	}

	templatesSearchResults struct {
		Templates []*types.Template
		Total     uint64
	}
)

// Search function Templates search
//
// expects implementation of search function:
// func (h templatesHandler) search(ctx context.Context, args *templatesSearchArgs) (results *templatesSearchResults, err error) {
//    return
// }
func (h templatesHandler) Search() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesSearch",
		Kind:   "function",
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Templates search",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "handle",
				Types: []string{"String"},
			},
			{
				Name:  "type",
				Types: []string{"String"},
			},
			{
				Name:  "ownerID",
				Types: []string{"ID"},
			},
			{
				Name:  "partial",
				Types: []string{"Boolean"},
			},
			{
				Name:  "labels",
				Types: []string{"KV"},
			},
			{
				Name:  "sort",
				Types: []string{"String"},
			},
			{
				Name:  "limit",
				Types: []string{"UnsignedInteger"},
			},
			{
				Name:  "incTotal",
				Types: []string{"Boolean"},
			},
			{
				Name:  "incPageNavigation",
				Types: []string{"Boolean"},
			},
			{
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:    "templates",
				Types:   []string{"Template"},
				IsArray: true,
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesSearchArgs{
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
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *templatesSearchResults
			if results, err = h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Templates (*types.Template) to Array (of Template)
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Templates))
				)

				for i := range results.Templates {
					if tarr[i], err = h.reg.Type("Template").Cast(results.Templates[i]); err != nil {
						return
					}
				}

				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "templates", tval); err != nil {
					return
				}
			}

			{
				// converting results.Total (uint64) to UnsignedInteger
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("UnsignedInteger").Cast(results.Total); err != nil {
					return
				} else if err = expr.Assign(out, "total", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	templatesEachArgs struct {
		hasHandle bool
		Handle    string

		hasType bool
		Type    string

		hasOwnerID bool
		OwnerID    uint64

		hasPartial bool
		Partial    bool

		hasLabels bool
		Labels    map[string]string

		hasSort bool
		Sort    string

		hasLimit bool
		Limit    uint64

		hasIncTotal bool
		IncTotal    bool

		hasIncPageNavigation bool
		IncPageNavigation    bool

		hasPageCursor bool
		PageCursor    string
	}

	templatesEachResults struct {
		Template *types.Template
		Total    uint64
	}
)

// Each function Tempates
//
// expects implementation of each function:
// func (h templatesHandler) each(ctx context.Context, args *templatesEachArgs) (results *templatesEachResults, err error) {
//    return
// }
func (h templatesHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesEach",
		Kind:   "iterator",
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Tempates",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "handle",
				Types: []string{"String"},
			},
			{
				Name:  "type",
				Types: []string{"String"},
			},
			{
				Name:  "ownerID",
				Types: []string{"ID"},
			},
			{
				Name:  "partial",
				Types: []string{"Boolean"},
			},
			{
				Name:  "labels",
				Types: []string{"KV"},
			},
			{
				Name:  "sort",
				Types: []string{"String"},
			},
			{
				Name:  "limit",
				Types: []string{"UnsignedInteger"},
			},
			{
				Name:  "incTotal",
				Types: []string{"Boolean"},
			},
			{
				Name:  "incPageNavigation",
				Types: []string{"Boolean"},
			},
			{
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "template",
				Types: []string{"Template"},
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &templatesEachArgs{
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
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return h.each(ctx, args)
		},
	}
}

type (
	templatesCreateArgs struct {
		hasTemplate bool
		Template    *types.Template
	}

	templatesCreateResults struct {
		Template *types.Template
	}
)

// Create function Template create
//
// expects implementation of create function:
// func (h templatesHandler) create(ctx context.Context, args *templatesCreateArgs) (results *templatesCreateResults, err error) {
//    return
// }
func (h templatesHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesCreate",
		Kind:   "function",
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Template create",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "template",
				Types: []string{"Template"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "template",
				Types: []string{"Template"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesCreateArgs{
					hasTemplate: in.Has("template"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *templatesCreateResults
			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Template (*types.Template) to Template
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Template").Cast(results.Template); err != nil {
					return
				} else if err = expr.Assign(out, "template", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	templatesUpdateArgs struct {
		hasTemplate bool
		Template    *types.Template
	}

	templatesUpdateResults struct {
		Template *types.Template
	}
)

// Update function Template update
//
// expects implementation of update function:
// func (h templatesHandler) update(ctx context.Context, args *templatesUpdateArgs) (results *templatesUpdateResults, err error) {
//    return
// }
func (h templatesHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesUpdate",
		Kind:   "function",
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Template update",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "template",
				Types: []string{"Template"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "template",
				Types: []string{"Template"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesUpdateArgs{
					hasTemplate: in.Has("template"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *templatesUpdateResults
			if results, err = h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Template (*types.Template) to Template
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Template").Cast(results.Template); err != nil {
					return
				} else if err = expr.Assign(out, "template", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	templatesDeleteArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Template
	}
)

func (a templatesDeleteArgs) GetLookup() (bool, uint64, string, *types.Template) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Delete function Template delete
//
// expects implementation of delete function:
// func (h templatesHandler) delete(ctx context.Context, args *templatesDeleteArgs) (err error) {
//    return
// }
func (h templatesHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesDelete",
		Kind:   "function",
		Labels: map[string]string{"delete": "step", "templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Template delete",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Template"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesDeleteArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			return out, h.delete(ctx, args)
		},
	}
}

type (
	templatesRecoverArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Template
	}
)

func (a templatesRecoverArgs) GetLookup() (bool, uint64, string, *types.Template) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Recover function Template recover
//
// expects implementation of recover function:
// func (h templatesHandler) recover(ctx context.Context, args *templatesRecoverArgs) (err error) {
//    return
// }
func (h templatesHandler) Recover() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesRecover",
		Kind:   "function",
		Labels: map[string]string{"recover": "step", "templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Template recover",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Template"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesRecoverArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			return out, h.recover(ctx, args)
		},
	}
}

type (
	templatesRenderArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Template

		hasDocumentName bool
		DocumentName    string

		hasDocumentType bool
		DocumentType    string

		hasVariables bool
		Variables    expr.RVars

		hasOptions bool
		Options    map[string]string
	}

	templatesRenderResults struct {
		Document *renderedDocument
	}
)

func (a templatesRenderArgs) GetLookup() (bool, uint64, string, *types.Template) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Render function Render template
//
// expects implementation of render function:
// func (h templatesHandler) render(ctx context.Context, args *templatesRenderArgs) (results *templatesRenderResults, err error) {
//    return
// }
func (h templatesHandler) Render() *atypes.Function {
	return &atypes.Function{
		Ref:    "templatesRender",
		Kind:   "function",
		Labels: map[string]string{"render": "step", "templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Render template",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Template"}, Required: true,
			},
			{
				Name:  "documentName",
				Types: []string{"String"},
			},
			{
				Name:  "documentType",
				Types: []string{"String"},
			},
			{
				Name:  "variables",
				Types: []string{"Vars"},
			},
			{
				Name:  "options",
				Types: []string{"RenderOptions"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "document",
				Types: []string{"RenderedDocument"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &templatesRenderArgs{
					hasLookup:       in.Has("lookup"),
					hasDocumentName: in.Has("documentName"),
					hasDocumentType: in.Has("documentType"),
					hasVariables:    in.Has("variables"),
					hasOptions:      in.Has("options"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Template").Type():
					args.lookupRes = aux.Get().(*types.Template)
				}
			}

			var results *templatesRenderResults
			if results, err = h.render(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Document (*renderedDocument) to RenderedDocument
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("RenderedDocument").Cast(results.Document); err != nil {
					return
				} else if err = expr.Assign(out, "document", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
