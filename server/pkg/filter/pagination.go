package filter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/modern-go/reflect2"
)

type (
	// Paging is a helper struct that should be embedded in filter types
	// to help with the paging
	Paging struct {
		// How many items per fetch do we want
		Limit uint `json:"limit,omitempty"`

		// PeakNextPage bool `json:"peakNextPage,omitempty"`

		// Opaque pointer to 1st item on page
		// read-only
		PageCursor *PagingCursor `json:"cursor,omitempty"`

		// Opaque cursor that points to the first item on the next page
		// read-only
		NextPage *PagingCursor `json:"nextPage,omitempty"`

		// Opaque cursor that points to the first item on the previous page
		// value of {cursor} will be copied here
		// read-only
		PrevPage *PagingCursor `json:"prevPage,omitempty"`

		IncPageNavigation bool `json:"incPageNavigation,omitempty"`
		IncTotal          bool `json:"incTotal,omitempty"`

		PageNavigation []*Page `json:"pageNavigation,omitempty"`
		Total          uint    `json:"total,omitempty"`
	}

	Page struct {
		Page   uint          `json:"page"`
		Count  uint          `json:"items"`
		Cursor *PagingCursor `json:"cursor"`
	}

	PagingCursor struct {
		keys     []string
		kk       [][]string
		values   []interface{}
		modifier []string
		desc     []bool

		// sort in desc order
		ROrder bool

		// use < op instead of >
		LThen bool
	}

	pagingCursorValue struct {
		v interface{}
	}

	// copy of the DAL valueGetter so we don't need a reference to pkg/dal
	valueGetter interface {
		GetValue(string, uint) (any, error)
		CountValues() map[string]uint
	}
)

var (
	ErrIncompatibleSort = fmt.Errorf("sort incompatible with cursor; send empty sort")
)

func NewPaging(limit uint, cursor string) (p Paging, err error) {
	p = Paging{Limit: limit}
	if p.PageCursor, err = parseCursor(cursor); err != nil {
		return
	}

	return
}

func (p *Paging) GetLimit() uint {
	return p.Limit
}

func (p *Paging) Clone() *Paging {
	if p == nil {
		return nil
	}

	return &Paging{
		Limit:             p.Limit,
		PageCursor:        p.PageCursor,
		NextPage:          p.NextPage,
		PrevPage:          p.PrevPage,
		IncPageNavigation: p.IncPageNavigation,
		IncTotal:          p.IncTotal,
		PageNavigation:    p.PageNavigation,
		Total:             p.Total,
	}
}

func (p *PagingCursor) Walk(fn func(string, interface{}, bool)) {
	for i, key := range p.keys {
		fn(key, p.values[i], p.desc[i])
	}
}

func (p *PagingCursor) Set(k string, v interface{}, d bool) {
	p.SetModifier(k, v, d, "", k)
}

func (p *PagingCursor) SetModifier(k string, v interface{}, d bool, m string, kk ...string) {
	for i, key := range p.keys {
		if key == k {
			p.values[i] = v
			return
		}
	}

	p.keys = append(p.keys, k)
	p.values = append(p.values, v)
	p.kk = append(p.kk, kk)
	p.modifier = append(p.modifier, m)
	p.desc = append(p.desc, d)
}

func (p *PagingCursor) Keys() []string {
	return p.keys
}

func (p *PagingCursor) Values() []interface{} {
	return p.values
}

func (p *PagingCursor) KK() [][]string {
	return p.kk
}

func (p *PagingCursor) Modifiers() []string {
	return p.modifier
}

func (p *PagingCursor) IsLThen() bool {
	return p.LThen
}

func (p *PagingCursor) Desc() []bool {
	return p.desc
}

func (p *PagingCursor) IsROrder() bool {
	return p.ROrder
}

// Stirng to implement Stringer and to get human-readable representation of the cursor
func (p *PagingCursor) String() string {

	var o = "<"

	if p == nil {
		o += "nil"
	} else {
		for i, key := range p.keys {
			o += fmt.Sprintf("%s: %v", key, p.values[i])
			if p.desc[i] {
				o += " DESC"
			}
			o += ", "
		}

		if p.ROrder {
			o += "[REV"
		} else {
			o += "[FWD"
		}

		if p.LThen {
			o += ",<]"
		} else {
			o += ",>]"
		}
	}

	return o + ">"
}

// MarshalJSON serializes cursor struct as JSON and encodes it as base64 + adds quotes to be treated as JSON string
func (p *PagingCursor) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(struct {
		K  []string
		KK [][]string
		V  []interface{}
		M  []string
		D  []bool
		R  bool
		LT bool
	}{
		p.keys,
		p.kk,
		p.values,
		p.modifier,
		p.desc,
		p.ROrder,
		p.LThen,
	})

	if err != nil {
		return nil, err
	}

	std := base64.StdEncoding
	dbuf := make([]byte, std.EncodedLen(len(buf)))
	std.Encode(dbuf, buf)

	return append([]byte{'"'}, append(dbuf, '"')...), nil
}

func (p *PagingCursor) Encode() string {
	b, _ := p.MarshalJSON()
	return string(b)
}

func (p *PagingCursor) UnmarshalJSON(in []byte) error {
	var (
		aux struct {
			K  []string
			KK [][]string
			V  []pagingCursorValue
			M  []string
			D  []bool
			R  bool
			LT bool
		}

		err error
	)

	// Decode the the string if it's not a JSON
	if in[0] != byte('{') {
		s := string(in)
		if s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		in, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(in, &aux); err != nil {
		return err
	}

	p.keys = aux.K
	p.kk = aux.KK
	p.modifier = aux.M
	p.desc = aux.D
	p.ROrder = aux.R
	p.LThen = aux.LT

	// json.Unmarshal treats uint64 in values ([]interface{}) as float64 and we don't like that.
	p.values = make([]interface{}, len(aux.V))
	for i, v := range aux.V {
		p.values[i] = v.v
	}

	return nil
}

func (p *PagingCursor) Decode(cursor string) error {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	err = p.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	return err
}

// Sort returns:
//  - sort if cursor is nil
//  - sort from cursor when sort is empty
//  - sort from cursor when sort is compatible with cursor
//  - error if sort & cursor are incompatible
func (p *PagingCursor) Sort(sort SortExprSet) (SortExprSet, error) {
	if p == nil {
		return sort, nil
	}

	if len(sort) == 0 {
		// sort empty, create it from cursor
		for k := range p.keys {
			sort = append(sort, &SortExpr{
				Column:     p.keys[k],
				columns:    p.kk[k],
				modifier:   p.modifier[k],
				Descending: p.desc[k],
			})
		}
		return sort, nil
	}

	// make sure there are at least as many keys in the sort as they
	// are in the cursor
	if len(p.keys) < len(sort) {
		return nil, ErrIncompatibleSort
	}

	// make sure keys from sort match the ones from cursor
	for k := range sort {
		if !(p.keys[k] == sort[k].Column && p.desc[k] == sort[k].Descending) {
			return nil, ErrIncompatibleSort
		}
	}

	return sort, nil
}

// ToAST converts the given PagingCursor to a corresponding AST tree
//
// The method should be used as the base for generating filtering expressions
// when working with databases, DAL reports, ...
//
// @todo discuss this one
func (cur *PagingCursor) ToAST(identLookup func(i string) (string, error), castFn func(i string, val any) (expr.TypedValue, error)) (out *ql.ASTNode, err error) {
	var (
		cc = cur.Keys()
		vv = cur.Values()

		value expr.TypedValue

		ident string

		ltOp = map[bool]string{
			true:  "lt",
			false: "gt",
		}

		// Modifying this function to use expressions instead of constant boolean
		// values because MSSQL doesn't have those.
		//
		// @todo rethink and redo the whole/all of the filtering logic surrounding paging
		// cursors to make them consistent/reusable
		isValueNull = func(i int, neg bool) *ql.ASTNode {
			out := &ql.ASTNode{
				Ref: "eq",
				Args: ql.ASTNodeSet{{
					Value: ql.MakeValueOf("Integer", 1),
				}},
			}
			if (reflect2.IsNil(vv[i]) && !neg) || (!reflect2.IsNil(vv[i]) && neg) {
				// Makes the expr. true
				out.Args = append(out.Args, &ql.ASTNode{Value: ql.MakeValueOf("Integer", 1)})
			}
			// Makes the expr false
			out.Args = append(out.Args, &ql.ASTNode{Value: ql.MakeValueOf("Integer", 0)})

			return out
		}
	)

	_ = value

	if len(cc) == 0 {
		return
	}

	// going from the last key/column to the 1st one
	for i := len(cc) - 1; i >= 0; i-- {
		if identLookup != nil {
			// Get the key context so we know how to format fields and format typecasts
			ident, err = identLookup(cc[i])
			if err != nil {
				return
			}
		} else {
			ident = cc[i]
		}

		if castFn == nil {
			value, err = expr.Typify(vv[i])
		} else {
			value, err = castFn(cc[i], vv[i])
		}
		if err != nil {
			return
		}

		// We need to cut off the values that are before the cursor (when ascending)
		// and vice-versa for descending.
		lt := cur.Desc()[i]
		if cur.IsROrder() {
			lt = !lt
		}

		op := ltOp[lt]

		// // Typecast the value so comparison can work properly

		// Either BOTH (field and value) are NULL or field is grater-then value
		base := &ql.ASTNode{
			Ref: "group",
			Args: ql.ASTNodeSet{
				&ql.ASTNode{
					Ref: "or",
					Args: ql.ASTNodeSet{
						&ql.ASTNode{
							Ref: "group",
							Args: ql.ASTNodeSet{
								&ql.ASTNode{
									Ref: "and",
									Args: ql.ASTNodeSet{
										&ql.ASTNode{
											Ref: "nnull",
											Args: ql.ASTNodeSet{
												&ql.ASTNode{
													Symbol: ident,
												},
											},
										},
										isValueNull(i, false),
									},
								},
							},
						},
						&ql.ASTNode{
							Ref: op,
							Args: ql.ASTNodeSet{{
								Symbol: ident,
							}, {
								Value: ql.WrapValue(value),
							}},
						},
					},
				},
			},
		}

		if out == nil {
			out = base
		} else {
			out = &ql.ASTNode{
				Ref: "group",
				Args: ql.ASTNodeSet{
					&ql.ASTNode{
						Ref: "or",
						Args: ql.ASTNodeSet{
							base,
							&ql.ASTNode{
								Ref: "group",
								Args: ql.ASTNodeSet{
									&ql.ASTNode{
										Ref: "and",
										Args: ql.ASTNodeSet{
											&ql.ASTNode{
												Ref: "group",
												Args: ql.ASTNodeSet{
													&ql.ASTNode{
														Ref: "or",
														Args: ql.ASTNodeSet{
															&ql.ASTNode{
																Ref: "and",
																Args: ql.ASTNodeSet{
																	&ql.ASTNode{
																		Ref: "null",
																		Args: ql.ASTNodeSet{
																			&ql.ASTNode{
																				Symbol: ident,
																			},
																		},
																	},
																	isValueNull(i, false),
																},
															},

															&ql.ASTNode{
																Ref: "eq",
																Args: ql.ASTNodeSet{{
																	Symbol: ident,
																}, {
																	Value: ql.WrapValue(value),
																}},
															},
														},
													},
												},
											},
											out,
										},
									},
								},
							},
						},
					},
				},
			}

		}
	}

	return
}

// PagingCursorFrom constructs a new paging cursor for the given valueGetter
func PagingCursorFrom(ss SortExprSet, v valueGetter, primaries ...string) (_ *PagingCursor, err error) {
	var (
		cur     = &PagingCursor{LThen: ss.Reversed()}
		pkUsed  = make(map[string]bool)
		pkIndex = make(map[string]bool)

		value any
	)

	// Index all of the primary keys for easier code
	for _, a := range primaries {
		pkIndex[a] = true
	}

	if len(pkIndex) == 0 {
		err = fmt.Errorf("can not construct cursor without primary key attributes")
		return
	}

	for _, s := range ss {
		if ok := pkIndex[s.Column]; ok {
			pkUsed[s.Column] = true
		}

		// @todo multi values?
		value, err = v.GetValue(s.Column, 0)
		if err != nil {
			return
		}

		cur.Set(s.Column, value, s.Descending)
	}

	// Make sure the rest of the unused primary keys are applied
	if len(pkUsed) != len(pkIndex) {
		for _, ident := range primaries {
			if _, ok := pkUsed[ident]; ok {
				continue
			}

			value, err = v.GetValue(ident, 0)
			if err != nil {
				return
			}

			cur.Set(ident, value, false)
		}
	}

	return cur, nil
}

func parseCursor(in string) (p *PagingCursor, err error) {
	if len(in) == 0 {
		return nil, nil
	}

	var buf []byte
	if buf, err = base64.StdEncoding.DecodeString(in); err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	p = &PagingCursor{}
	if err = p.UnmarshalJSON(buf); err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	return p, nil
}

// Making sure uint64 other int* values are properly unmarshaled
func (v *pagingCursorValue) UnmarshalJSON(in []byte) (err error) {
	var (
		u uint64
		i int64
	)

	if string(in) == "null" {
		// if we do not do this we risk conversion to int(0)
		v.v = nil
		return
	}

	if err = json.Unmarshal(in, &u); err == nil {
		// handle big integers properly
		v.v = u
		return
	}

	if err = json.Unmarshal(in, &i); err == nil {
		v.v = i
		return
	}

	return json.Unmarshal(in, &v.v)
}
