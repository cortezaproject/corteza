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
	},

		atypes.ConstructTrigger{
			ResourceType: "compose",
			EventType:    "onInterval",
			Constraints: []atypes.ConstructTriggerConstraint{{
				Name:  "interval",
				Types: []string{"String"},
				Meta:  atypes.ConstructTriggerConstraintMeta{},
			}},
			Properties: []atypes.ConstructTriggerProperty{},
		},

		atypes.ConstructTrigger{
			ResourceType: "compose",
			EventType:    "onTimestamp",
			Constraints: []atypes.ConstructTriggerConstraint{{
				Name:  "timestamp",
				Types: []string{"String"},
				Meta:  atypes.ConstructTriggerConstraintMeta{},
			}},
			Properties: []atypes.ConstructTriggerProperty{},
		},

		atypes.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterCreate",
			Segments: []atypes.ConstructSegment{
				{
					Sections: []atypes.ConstructSection{
						{
							Elements: []atypes.SectionElement{
								{
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
								},
							},
						},
					},
				},
			},
			Properties: []atypes.ConstructTriggerProperty{
				{
					Name: "record",
					Type: "ComposeRecord",
				},

				{
					Name: "oldRecord",
					Type: "ComposeRecord",
				},

				{
					Name: "module",
					Type: "ComposeModule",
				},

				{
					Name: "namespace",
					Type: "ComposeNamespace",
				},

				{
					Name: "recordValueErrors",
					Type: "ComposeRecordValueErrorSet",
				},

				{
					Name: "selected",
					Type: "",
				},
			},
			Constraints: []atypes.ConstructTriggerConstraint{
				{
					Name:  "namespace",
					Types: []string{"Handle", "String"},
				},
				{
					Name:  "module",
					Types: []string{"Handle", "String"},
				},
			},
		},

		atypes.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterUpdate",
			Properties: []atypes.ConstructTriggerProperty{
				{
					Name: "record",
					Type: "ComposeRecord",
				},

				{
					Name: "oldRecord",
					Type: "ComposeRecord",
				},

				{
					Name: "module",
					Type: "ComposeModule",
				},

				{
					Name: "namespace",
					Type: "ComposeNamespace",
				},

				{
					Name: "recordValueErrors",
					Type: "ComposeRecordValueErrorSet",
				},

				{
					Name: "selected",
					Type: "",
				},
			},
			Segments: []atypes.ConstructSegment{
				{
					Sections: []atypes.ConstructSection{
						{
							Elements: []atypes.SectionElement{
								{
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
								},
							},
						},
					},
				},
			},
			Constraints: []atypes.ConstructTriggerConstraint{
				{
					Name:  "namespace",
					Types: []string{"Handle", "String"},
				},
				{
					Name:  "module",
					Types: []string{"Handle", "String"},
				},
			},
		},

		atypes.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterDelete",
			Properties: []atypes.ConstructTriggerProperty{

				{
					Name: "record",
					Type: "ComposeRecord",
				},

				{
					Name: "oldRecord",
					Type: "ComposeRecord",
				},

				{
					Name: "module",
					Type: "ComposeModule",
				},

				{
					Name: "namespace",
					Type: "ComposeNamespace",
				},

				{
					Name: "recordValueErrors",
					Type: "ComposeRecordValueErrorSet",
				},

				{
					Name: "selected",
					Type: "",
				},
			},
			Segments: []atypes.ConstructSegment{
				{
					Sections: []atypes.ConstructSection{
						{
							Elements: []atypes.SectionElement{
								{
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
								},
							},
						},
					},
				},
			},
			Constraints: []atypes.ConstructTriggerConstraint{
				{
					Name:  "namespace",
					Types: []string{"Handle", "String"},
				},
				{
					Name:  "module",
					Types: []string{"Handle", "String"},
				},
			},
		})

	h.reg.AddFunctions(
		h.Lookup(),
		h.Create(),
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
				ArgumentName: "namespace",
				Name:         "",
				Types:        []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				ArgumentName: "module",
				Name:         "",
				Types:        []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				ArgumentName: "record",
				Name:         "",
				Types:        []string{"ID", "ComposeRecord"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				ArgumentName: "record",
				Name:         "",
				Types:        []string{"ComposeRecord"},
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
						Context: atypes.SectionElementInputContext{
							DependsOn: map[string]string{
								"namespaceID": "namespace",
							},
						},
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "RecordSelector",
						Label:    "Record",
						Argument: "record",
						Context: atypes.SectionElementInputContext{
							DependsOn: map[string]string{
								"namespaceID": "namespace",
								"moduleID":    "module",
							},
						},
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

type (
	ngRecordsCreateArgs struct {
		hasNamespace    bool
		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceRes    *types.Namespace

		hasModule    bool
		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleRes    *types.Module

		hasValues  bool
		Values     interface{}
		ValuesKV   map[string]string
		ValuesKVV  map[string][]string
		ValuesVars *expr.Vars
	}
)

func (a ngRecordsCreateArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

func (a ngRecordsCreateArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (h ngRecordsHandler) Create() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "composeRecordsCreate",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "create": "step", "record": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Compose record create",
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "namespace",
				Name:         "",
				Types:        []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				ArgumentName: "module",
				Name:         "",
				Types:        []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},

			{
				ArgumentName: "values",
				Name:         "",
				Types:        []string{"KV", "KVV", "Any"}, Required: true,
				Aggregate: true,
			},
		},

		Results: []*atypes.Param{

			{
				ArgumentName: "record",
				Name:         "",
				Types:        []string{"ComposeRecord"},
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
						Context: atypes.SectionElementInputContext{
							DependsOn: map[string]string{
								"namespaceID": "namespace",
							},
						},
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:     "FieldValueMap",
						Label:    "Values",
						Argument: "__fieldValues",
						Context: atypes.SectionElementInputContext{
							DependsOn: map[string]string{
								"namespaceID": "namespace",
								"moduleID":    "module",
							},
						},
					},
				}},
			}},
		}},

		ArgsMerger: func(ctx context.Context, args atypes.ExprSet, raw []expr.TypedValue) (out *expr.Vars, err error) {
			aux := make(map[string]any, 2)
			aux[args[0].ArgumentName] = raw[0]
			aux[args[1].ArgumentName] = raw[1]

			auxVals := make(map[string][]string, 4)
			for i := 2; i < len(raw); i++ {
				rv := raw[i]
				switch rv.Type() {
				case h.tReg.Type("KV").Type():
					for k, v := range rv.Get().(map[string]string) {
						auxVals[k] = append(auxVals[k], v)
					}
				case h.tReg.Type("KVV").Type():
					for k, v := range rv.Get().(map[string][]string) {
						auxVals[k] = append(auxVals[k], v...)
					}

				default:
					var s string
					s, err = expr.CastToString(rv.Get())
					if err != nil {
						return
					}

					auxVals[args[i].Target] = append(auxVals[args[i].Target], s)
				}
			}

			aux[args[2].ArgumentName], err = expr.NewKVV(auxVals)
			if err != nil {
				return
			}

			return expr.NewVars(aux)
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &ngRecordsCreateArgs{
					hasNamespace: in.Has("namespace"),
					hasModule:    in.Has("module"),
					hasValues:    in.Has("values"),
				}
			)

			if err = in.Decode(in); err != nil {
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

			// Converting Values argument
			if args.hasValues {
				aux := expr.Must(expr.Select(in, "values"))
				switch aux.Type() {
				case h.tReg.Type("KV").Type():
					args.ValuesKV = aux.Get().(map[string]string)
				case h.tReg.Type("KVV").Type():
					args.ValuesKVV = aux.Get().(map[string][]string)
				}
			}

			var results *recordsCreateResults
			if results, err = h.create(ctx, args); err != nil {
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

func (h ngRecordsHandler) create(ctx context.Context, args *ngRecordsCreateArgs) (results *recordsCreateResults, err error) {
	results = &recordsCreateResults{}

	namespace, module, err := h.loadCombo(ctx, args)
	if err != nil {
		return nil, err
	}

	// @todo kvv
	values := make(types.RecordValueSet, 0, len(args.ValuesKV))
	for k, v := range args.ValuesKV {
		values = append(values, &types.RecordValue{
			Name:  k,
			Value: v,
		})
	}

	for k, vv := range args.ValuesKVV {
		for i, v := range vv {
			values = append(values, &types.RecordValue{
				Name:  k,
				Place: uint(i),
				Value: v,
			})
		}
	}

	record := &types.Record{
		NamespaceID: namespace.ID,
		ModuleID:    module.ID,
		Values:      values,
	}

	results.Record, err = wrapRecordValueErrorSet(h.rec.Create(ctx, record))

	return
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
