package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// OrderedBuffer provides the buffered data in the specified order
	OrderedBuffer interface {
		Buffer

		// InOrder makes the buffer provide the stored data in the specified order
		InOrder(ss ...*filter.SortExpr) (err error)
	}

	// Buffer provides a place where you can buffer the data provided by DAL
	Buffer interface {
		// Seek moves the index pointer to the specified location
		// After the Seek call, a Next() call is required
		Seek(context.Context, int) error

		// Len returns the number of elements in the buffer
		Len() int

		Iterator
		Adder
	}

	Adder interface {
		// Add adds a new ValueGetter to the buffer
		Add(context.Context, ValueGetter) (err error)
	}

	Row struct {
		counters map[string]uint
		values   valueSet

		// ...

		// Metadata to make it easier to work with
		// @todo add when needed
	}

	valueSet map[string][]any
)

func (r Row) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	if r.values[k] == nil {
		return nil, nil
	}

	if len(r.values[k]) == 0 {
		return nil, nil
	}

	o := r.values[k][0]
	return o, nil
}

func (r *Row) Reset() {
	for k := range r.counters {
		r.counters[k] = 0
	}
}

func (r *Row) SetValue(name string, pos uint, v any) error {
	if r.values == nil {
		r.values = make(valueSet)
	}
	if r.counters == nil {
		r.counters = make(map[string]uint)
	}

	// Make sure there is space for it
	// @note benchmarking proves that the rest of the function introduces
	//       a lot of memory pressure.
	//       Investigate options on reworking this/reducing allocations.
	if int(pos)+1 > len(r.values[name]) {
		r.values[name] = append(r.values[name], make([]any, (int(pos)+1)-len(r.values[name]))...)
	}

	r.values[name][pos] = v
	if pos >= r.counters[name] {
		r.counters[name]++
	}

	return nil
}

// WithValue is a simple helper to construct rows with populated values
// The main use is for tests so restrain from using it in code.
func (r *Row) WithValue(name string, pos uint, v any) *Row {
	err := r.SetValue(name, pos, v)
	if err != nil {
		panic(err)
	}

	return r
}

func (r *Row) CountValues() map[string]uint {
	return r.counters
}

func (r *Row) GetValue(name string, pos uint) (any, error) {
	if r.values == nil {
		return nil, nil
	}
	if r.counters == nil {
		return nil, nil
	}
	if pos >= r.counters[name] {
		return nil, nil
	}

	return r.values[name][pos], nil
}

func (r *Row) String() string {
	out := make([]string, 0, 20)
	for k, vv := range r.values {
		for i, v := range vv {
			out = append(out, fmt.Sprintf("%s [%d] %v", k, i, v))
		}
	}

	return strings.Join(out, " | ")
}

func (r Row) Copy() *Row {
	out := &r

	out.values = out.values.Copy()

	return out
}

func (vv valueSet) Copy() valueSet {
	out := make(valueSet)

	for n, vv := range vv {
		out[n] = vv
	}

	return out
}

func mergeRows(mapping []AttributeMapping, dst *Row, ss ...*Row) (err error) {
	if len(mapping) == 0 {
		return mergeRowsFull(dst, ss...)
	}

	return mergeRowsMapped(mapping, dst, ss...)
}

func mergeRowsFull(dst *Row, rows ...*Row) (err error) {
	for _, r := range rows {
		for name, vv := range r.values {
			for i, values := range vv {
				if dst.values == nil {
					dst.values = make(valueSet)
					dst.counters = make(map[string]uint)
				}

				if i == 0 {
					dst.values[name] = make([]any, len(vv))
					dst.counters[name] = 0
				}

				err = dst.SetValue(name, uint(i), values)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func mergeRowsMapped(mapping []AttributeMapping, out *Row, rows ...*Row) (err error) {

	for _, mp := range mapping {
		name := mp.Source()
		for _, r := range rows {
			if r.values[name] != nil {
				if out.values == nil {
					out.values = make(valueSet)
					out.counters = make(map[string]uint)
				}

				out.values[mp.Identifier()] = r.values[name]
				out.counters[mp.Identifier()] = r.counters[name]
				break
			}
		}
	}

	return
}

// makeRowComparator is a utility for easily making a row comparator for
// the given sort expression
func makeRowComparator(ss ...*filter.SortExpr) func(a, b *Row) bool {
	return func(a, b *Row) bool {
		for _, s := range ss {
			cmp := compareGetters(a, b, a.counters, b.counters, s.Column)

			less, skip := evalCmpResult(cmp, s)
			if !skip {
				return less
			}
		}

		return false
	}
}

func evalCmpResult(cmp int, s *filter.SortExpr) (less, skip bool) {
	if cmp != 0 {
		if s.Descending {
			return cmp > 0, false
		}
		return cmp < 0, false
	}

	return false, true
}
