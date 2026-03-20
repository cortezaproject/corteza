package automation

import (
	"context"
	"fmt"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
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
		Groups:       []string{"Manual"},
		Meta: &atypes.ConstructTriggerMeta{},
	},

		atypes.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onInterval",
			Groups:       []string{"Schedule"},
			Meta:         &atypes.ConstructTriggerMeta{},
			Segments: []atypes.ConstructSegment{
				{
					Sections: []atypes.ConstructSection{
						{
							Elements: []atypes.SectionElement{
								{
									Input: atypes.SectionElementInput{
										Type:        "String",
										Label:       "Interval (Cron)",
										Argument:    "interval",
										Placeholder: "* * * * *",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []atypes.ConstructTriggerConstraint{{
				Name:  "interval",
				Types: []string{"String"},
				Meta:  atypes.ConstructTriggerConstraintMeta{},
			}},
			Properties: []atypes.ConstructTriggerProperty{},
		},

		atypes.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onTimestamp",
			Groups:       []string{"Schedule"},
			Meta:         &atypes.ConstructTriggerMeta{},
			Segments: []atypes.ConstructSegment{
				{
					Sections: []atypes.ConstructSection{
						{
							Elements: []atypes.SectionElement{
								{
									Input: atypes.SectionElementInput{
										Type:        "DateTime",
										Label:       "Timestamp (RFC3339)",
										Argument:    "timestamp",
										Placeholder: "2026-01-02T15:04:05Z",
									},
								},
							},
						},
					},
				},
			},
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
			Groups:       []string{"Records"},
			Meta:         &atypes.ConstructTriggerMeta{},
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
										Context: atypes.SectionElementInputContext{
											DependsOn: map[string]string{
												"namespaceID": "namespace",
											},
										},
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
					Types: []string{"ID", "Handle", "String"},
				},
				{
					Name:  "module",
					Types: []string{"ID", "Handle", "String"},
				},
			},
		},

		atypes.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterUpdate",
			Groups:       []string{"Records"},
			Meta:         &atypes.ConstructTriggerMeta{},
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
										Context: atypes.SectionElementInputContext{
											DependsOn: map[string]string{
												"namespaceID": "namespace",
											},
										},
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
					Types: []string{"ID", "Handle", "String"},
				},
				{
					Name:  "module",
					Types: []string{"ID", "Handle", "String"},
				},
			},
		},

		atypes.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterDelete",
			Groups:       []string{"Records"},
			Meta:         &atypes.ConstructTriggerMeta{},
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
										Context: atypes.SectionElementInputContext{
											DependsOn: map[string]string{
												"namespaceID": "namespace",
											},
										},
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
					Types: []string{"ID", "Handle", "String"},
				},
				{
					Name:  "module",
					Types: []string{"ID", "Handle", "String"},
				},
			},
		})

	h.reg.AddFunctions(
		h.Lookup(),
		h.Create(),
		h.Each(),
	)
}

func (h ngRecordsHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:  "composeRecordsLookup",
		Kind: "function",
		Groups: []string{"Records"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Find Record",
			Description: "Lookup record by ID",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "database"},
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
		Groups: []string{"Records"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Create Record",
			Description: "Add new record to module",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "database"},
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
						Argument: "values",
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

func (h ngRecordsHandler) Each() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "composeRecordsEach",
		Kind:   "iterator",
		Groups: []string{"Loops"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Each Record",
			Description: "Iterate over records matching a query",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "refresh"},
		},

		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "namespace",
				Types:        []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				ArgumentName: "module",
				Types:        []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to iterate records from",
					Description: "Records are fetched from the specified module.",
				},
			},
			{
				ArgumentName: "query",
				Types:        []string{"String"},
				Meta: &atypes.ParamMeta{
					Label:       "Filter query",
					Description: "Optional filter string to match records.",
				},
			},
			{
				ArgumentName: "sort",
				Types:        []string{"String"},
				Meta: &atypes.ParamMeta{
					Label:       "Sort",
					Description: "Sort expression (e.g. \"createdAt DESC\").",
				},
			},
			{
				ArgumentName: "limit",
				Types:        []string{"Integer"},
				Meta: &atypes.ParamMeta{
					Label:       "Limit",
					Description: "Maximum number of records to iterate.",
				},
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "record",
				Types:        []string{"ComposeRecord"},
			},
			{
				ArgumentName: "index",
				Types:        []string{"Integer"},
			},
			{
				ArgumentName: "total",
				Types:        []string{"Integer"},
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
						Type:        "String",
						Label:       "Query",
						Argument:    "query",
						Placeholder: "Filter expression",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:        "String",
						Label:       "Sort",
						Argument:    "sort",
						Placeholder: "createdAt DESC",
					},
				}, {
					Input: atypes.SectionElementInput{
						Type:        "Number",
						Label:       "Limit",
						Argument:    "limit",
						Placeholder: "0 (no limit)",
					},
				}},
			}},
		}},

		Iterator: func(ctx context.Context, in *expr.Vars) (wfexec.IteratorHandler, error) {
			var (
				i = &recordSetIterator{}
				f = types.RecordFilter{}
			)

			// Resolve namespace + module
			type nsModArgs struct {
				hasNamespace    bool
				namespaceID     uint64
				namespaceHandle string
				namespaceRes    *types.Namespace

				hasModule    bool
				moduleID     uint64
				moduleHandle string
				moduleRes    *types.Module
			}
			args := &nsModArgs{
				hasNamespace: in.Has("namespace"),
				hasModule:    in.Has("module"),
			}

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

			// Resolve namespace
			var namespace *types.Namespace
			if args.namespaceRes != nil {
				namespace = args.namespaceRes
			} else if args.namespaceID > 0 {
				ns, err := h.ns.FindByID(ctx, args.namespaceID)
				if err != nil {
					return nil, fmt.Errorf("could not load namespace: %w", err)
				}
				namespace = ns
			} else if args.namespaceHandle != "" {
				ns, err := h.ns.FindByHandle(ctx, args.namespaceHandle)
				if err != nil {
					return nil, fmt.Errorf("could not load namespace: %w", err)
				}
				namespace = ns
			} else {
				return nil, fmt.Errorf("namespace is required")
			}

			// Resolve module
			var module *types.Module
			if args.moduleRes != nil {
				module = args.moduleRes
			} else if args.moduleID > 0 {
				mod, err := h.mod.FindByID(ctx, namespace.ID, args.moduleID)
				if err != nil {
					return nil, fmt.Errorf("could not load module: %w", err)
				}
				module = mod
			} else if args.moduleHandle != "" {
				mod, err := h.mod.FindByHandle(ctx, namespace.ID, args.moduleHandle)
				if err != nil {
					return nil, fmt.Errorf("could not load module: %w", err)
				}
				module = mod
			} else {
				return nil, fmt.Errorf("module is required")
			}

			f.ModuleID = module.ID
			f.NamespaceID = namespace.ID
			f.IncTotal = true

			// Optional: query
			if in.Has("query") {
				v := expr.Must(expr.Select(in, "query"))
				s, err := expr.CastToString(v.Get())
				if err != nil {
					return nil, fmt.Errorf("each: invalid query: %w", err)
				}
				f.Query = s
			}

			// Optional: sort
			if in.Has("sort") {
				v := expr.Must(expr.Select(in, "sort"))
				s, err := expr.CastToString(v.Get())
				if err != nil {
					return nil, fmt.Errorf("each: invalid sort: %w", err)
				}
				if err = f.Sort.Set(s); err != nil {
					return nil, fmt.Errorf("each: invalid sort expression: %w", err)
				}
			}

			// Optional: limit
			if in.Has("limit") {
				v := expr.Must(expr.Select(in, "limit"))
				lim, err := expr.CastToInteger(v.Get())
				if err != nil {
					return nil, fmt.Errorf("each: invalid limit: %w", err)
				}
				if lim > 0 {
					i.useIterLimit = true
					i.iterLimit = uint(lim)
					f.Limit = uint(lim)
					if lim > int64(wfexec.MaxIteratorBufferSize) {
						f.Limit = wfexec.MaxIteratorBufferSize
					}
				}
			}

			if f.Limit == 0 {
				f.Limit = wfexec.MaxIteratorBufferSize
			}

			i.filter = f
			i.loader = func() (err error) {
				if i.filter.PageCursor != nil && i.filter.NextPage == nil {
					return
				}

				i.total += i.ptr
				i.ptr = 0

				i.filter.PageCursor = i.filter.NextPage
				i.filter.NextPage = nil
				i.buffer, i.filter, err = h.rec.Find(ctx, i.filter)

				return
			}

			// Initial load
			return i, i.loader()
		},
	}
}

// ngRecordSetIterator wraps recordSetIterator to use the Ng type registry
// (reuses the existing recordSetIterator from the old workflow system)

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
