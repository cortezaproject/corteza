package automation

import (
	"context"
	"fmt"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"

	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	constructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
		AddTriggers(tt ...atypes.ConstructTrigger)
	}

	typeRegistry interface {
		Type(ref string) expr.Type
	}

	ngRecordsHandler struct {
		reg  constructSvc
		tReg typeRegistry
		ns   namespaceService
		mod  moduleService
		rec  recordService
	}
)

func NgRecordsHandler(reg constructSvc, tReg typeRegistry, ns namespaceService, mod moduleService, rec recordService) *ngRecordsHandler {
	h := &ngRecordsHandler{
		reg:  reg,
		tReg: tReg,
		ns:   ns,
		mod:  mod,
		rec:  rec,
	}

	h.register()
	return h
}

func (h ngRecordsHandler) register() {
	// @todo right place
	h.reg.AddTriggers(atypes.ConstructTrigger{
		ResourceType: "system",
		EventType:    "onManual",
	})

	h.reg.AddFunctions(
		h.Lookup(),
	)
}

func (h ngRecordsHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:  "composeRecordsLookup",
		Kind: "function",
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Compose record lookup",
			Description: "Find specific record by ID",
		},

		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},

		Parameters: []*atypes.Param{
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "record",
				Types: []string{"ID", "ComposeRecord"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "record",
				Types: []string{"ComposeRecord"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{{
					Input: atypes.SectionElementInput{
						Type:     "NamespaceSelector",
						Label:    "Namespace",
						Argument: "namespace",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "ModuleSelector",
						Label:    "Module",
						Argument: "module",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "RecordSelector",
						Label:    "Record",
						Argument: "record",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &recordsLookupArgs{
					hasNamespace: in.Has("namespace"),
					hasModule:    in.Has("module"),
					hasRecord:    in.Has("record"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Namespace argument
			if args.hasNamespace {
				aux := expr.Must(expr.Select(in, "namespace"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.namespaceID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.tReg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
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
				case h.tReg.Type("ComposeModule").Type():
					args.moduleRes = aux.Get().(*types.Module)
				}
			}

			// Converting Record argument
			if args.hasRecord {
				aux := expr.Must(expr.Select(in, "record"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.recordID = aux.Get().(uint64)
				case h.tReg.Type("ComposeRecord").Type():
					args.recordRes = aux.Get().(*types.Record)
				}
			}

			var results *recordsLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Record (*types.Record) to ComposeRecord
				var (
					tval expr.TypedValue
				)

				if tval, err = h.tReg.Type("ComposeRecord").Cast(results.Record); err != nil {
					return
				} else if err = expr.Assign(out, "record", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

func (h ngRecordsHandler) lookup(ctx context.Context, args *recordsLookupArgs) (results *recordsLookupResults, err error) {
	results = &recordsLookupResults{}
	results.Record, err = h.lookupRecord(ctx, args)
	return
}

func (h ngRecordsHandler) lookupRecord(ctx context.Context, args recordLookup) (record *types.Record, err error) {
	var (
		namespace *types.Namespace
		module    *types.Module
		recordID  uint64
	)

	if _, recordID, record = args.GetRecord(); record != nil {
		return
	}

	if namespace, module, err = h.loadCombo(ctx, args); err != nil {
		return
	}

	record, _, err = h.rec.FindByID(ctx, namespace.ID, module.ID, recordID)
	return
}

func (h ngRecordsHandler) loadCombo(ctx context.Context, args interface{}) (namespace *types.Namespace, module *types.Module, err error) {
	if lkp, is := args.(namespaceLookup); is {
		if namespace, err = lookupNamespace(ctx, h.ns, lkp); err != nil {
			err = fmt.Errorf("could not load namespace: %w", err)
			return
		}
	} else {
		err = fmt.Errorf("could not extract namespace lookup arguments")
		return
	}

	if lkp, is := args.(moduleLookup); is {
		if module, err = lookupModule(ctx, h.ns, h.mod, lkp); err != nil {
			err = fmt.Errorf("could not load module: %w", err)
			return
		}
	} else {
		err = fmt.Errorf("could not extract module lookup arguments")
		return
	}

	return
}
