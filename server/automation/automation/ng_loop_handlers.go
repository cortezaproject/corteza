package automation

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

type (
	ngLoopConstructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
		AddTriggers(tt ...atypes.ConstructTrigger)
	}

	ngLoopHandler struct {
		reg    ngLoopConstructSvc
		parser expr.Parsable
	}
)

func NgLoopHandler(reg ngLoopConstructSvc, parser expr.Parsable) *ngLoopHandler {
	h := &ngLoopHandler{
		reg:    reg,
		parser: parser,
	}

	h.register()
	return h
}

func (h ngLoopHandler) register() {
	h.reg.AddFunctions(
		h.Sequence(),
		h.Do(),
		// h.Each(),
		// h.Lines(),
	)
}

// Sequence iterates from first to last by step.
func (h ngLoopHandler) Sequence() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "loopSequence",
		Kind:   "iterator",
		Groups: []string{"Loops"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Count",
			Description: "Iterate from first to last by step",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "refresh"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "first",
				Types:        []string{"Integer"},
			},
			{
				ArgumentName: "last",
				Types:        []string{"Integer"},
				Required:     true,
			},
			{
				ArgumentName: "step",
				Types:        []string{"Integer"},
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "counter",
				Types:        []string{"Integer"},
			},
			{
				ArgumentName: "isFirst",
				Types:        []string{"Boolean"},
			},
			{
				ArgumentName: "isLast",
				Types:        []string{"Boolean"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{
						Input: atypes.SectionElementInput{
							Type:        "Number",
							Label:       "First",
							Argument:    "first",
							Placeholder: "0",
						},
					},
					{
						Input: atypes.SectionElementInput{
							Type:     "Number",
							Label:    "Last",
							Argument: "last",
						},
					},
					{
						Input: atypes.SectionElementInput{
							Type:        "Number",
							Label:       "Step",
							Argument:    "step",
							Placeholder: "1",
						},
					},
				},
			}},
		}},

		Iterator: func(ctx context.Context, in *expr.Vars) (wfexec.IteratorHandler, error) {
			var (
				first int64 = 0
				last  int64 = 1
				step  int64 = 1
				err   error
			)

			if in != nil {
				if in.Has("first") {
					v := expr.Must(expr.Select(in, "first"))
					if first, err = expr.CastToInteger(v.Get()); err != nil {
						return nil, fmt.Errorf("sequence: invalid first: %w", err)
					}
				}

				if in.Has("last") {
					v := expr.Must(expr.Select(in, "last"))
					if last, err = expr.CastToInteger(v.Get()); err != nil {
						return nil, fmt.Errorf("sequence: invalid last: %w", err)
					}
				}

				if in.Has("step") {
					v := expr.Must(expr.Select(in, "step"))
					if step, err = expr.CastToInteger(v.Get()); err != nil {
						return nil, fmt.Errorf("sequence: invalid step: %w", err)
					}
				}
			}

			sign := signOf(step)
			if first*sign >= last*sign {
				return nil, fmt.Errorf("sequence: first must be less than last (given direction)")
			}

			return &sequenceIterator{
				counter: first,
				cFirst:  first,
				cLast:   last,
				cStep:   step,
			}, nil
		},
	}
}

// Do iterates while a boolean expression is true.
func (h ngLoopHandler) Do() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "loopDo",
		Kind:   "iterator",
		Groups: []string{"Loops"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "While",
			Description: "Iterate while condition is true",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "refresh"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "while",
				Types:        []string{"String"},
				Required:     true,
				Meta: &atypes.ParamMeta{
					Label:       "Condition",
					Description: "Boolean expression evaluated each iteration",
				},
			},
		},

		Results: []*atypes.Param{},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{
						Input: atypes.SectionElementInput{
							Type:        "String",
							Label:       "While",
							Argument:    "while",
							Placeholder: "someVar > 0",
						},
					},
				},
			}},
		}},

		Iterator: func(ctx context.Context, in *expr.Vars) (wfexec.IteratorHandler, error) {
			var whileExpr string

			if in != nil && in.Has("while") {
				v := expr.Must(expr.Select(in, "while"))
				s, err := expr.CastToString(v.Get())
				if err != nil {
					return nil, fmt.Errorf("do: invalid while expression: %w", err)
				}
				whileExpr = s
			}

			if whileExpr == "" {
				return nil, fmt.Errorf("do: while expression is required")
			}

			i := &conditionIterator{}
			var err error
			if i.expr, err = h.parser.Parse(whileExpr); err != nil {
				return nil, fmt.Errorf("do: failed to parse while expression: %w", err)
			}

			return i, nil
		},
	}
}

// Each iterates over items in a collection (Array).
func (h ngLoopHandler) Each() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "loopEach",
		Kind:   "iterator",
		Groups: []string{"Loops"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Each",
			Description: "Iterate over items in a collection",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "refresh"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "items",
				Types:        []string{"Array"},
				Required:     true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "item",
				Types:        []string{"Any"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{
						Input: atypes.SectionElementInput{
							Type:     "Expression",
							Label:    "Items",
							Argument: "items",
						},
					},
				},
			}},
		}},

		Iterator: func(ctx context.Context, in *expr.Vars) (wfexec.IteratorHandler, error) {
			if in == nil || !in.Has("items") {
				return nil, fmt.Errorf("each: items argument is required")
			}

			v := expr.Must(expr.Select(in, "items"))
			arr, ok := v.Get().([]expr.TypedValue)
			if !ok {
				return nil, fmt.Errorf("each: expected Array, got %T", v.Get())
			}

			return &collectionIterator{set: arr}, nil
		},
	}
}

// Lines iterates over lines of a string (newline-separated).
func (h ngLoopHandler) Lines() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "loopLines",
		Kind:   "iterator",
		Groups: []string{"Loops"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Lines",
			Description: "Iterate over lines of a string",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "refresh"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "text",
				Types:        []string{"String"},
				Required:     true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "line",
				Types:        []string{"String"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{
						Input: atypes.SectionElementInput{
							Type:     "Expression",
							Label:    "Text",
							Argument: "text",
						},
					},
				},
			}},
		}},

		Iterator: func(ctx context.Context, in *expr.Vars) (wfexec.IteratorHandler, error) {
			if in == nil || !in.Has("text") {
				return nil, fmt.Errorf("lines: text argument is required")
			}

			v := expr.Must(expr.Select(in, "text"))
			s, err := expr.CastToString(v.Get())
			if err != nil {
				return nil, fmt.Errorf("lines: invalid text: %w", err)
			}

			return &lineIterator{s: bufio.NewScanner(strings.NewReader(s))}, nil
		},
	}
}
