package automation

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/spf13/cast"
	"strings"
	"time"
)

func CastToComposeNamespace(val interface{}) (out *types.Namespace, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.Namespace{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToComposeNamespace(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.Namespace:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToComposeModule(val interface{}) (out *types.Module, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.Module{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToComposeModule(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.Module:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToComposeRecord(val interface{}) (out *types.Record, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.Record{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToComposeRecord(out, k, v)
		})
	}
	switch val := expr.UntypedValue(val).(type) {
	case *types.Record:
		if val.Values == nil {
			val.Values = types.RecordValueSet{}
		}
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

var _ expr.DeepFieldAssigner = &ComposeRecord{}

// AssignFieldValue implements expr.DeepFieldAssigner
//
// We need to reroute value assigning for record-value-sets because
// we loose the reference to record-value slice
func (t *ComposeRecord) AssignFieldValue(kk []string, val interface{}) error {
	switch kk[0] {
	case "values":
		return assignToComposeRecordValues(&t.value.Values, kk[1:], val)
		// @todo deep setting labels
	default:
		return assignToComposeRecord(t.value, kk[0], val)
	}
}

var _ gval.Selector = &ComposeRecord{}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Record's underlying value (*types.Record)
// and it's fields
//
func (t ComposeRecord) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	if k == "values" {
		return t.value.Values.Dict(t.value.GetModule().Fields), nil
	}

	return composeRecordGValSelector(t.value, k)
}

func CastToComposeRecordValues(val interface{}) (out types.RecordValueSet, err error) {
	out = types.RecordValueSet{}
	switch val := val.(type) {
	case expr.Iterator:
		return out, val.Each(func(k string, v expr.TypedValue) error {
			// try with slice of strings first:
			if ss, err := cast.ToStringSliceE(expr.UntypedValue(v)); err == nil {
				for i, v := range ss {
					out = out.Set(&types.RecordValue{Name: k, Value: v, Place: uint(i)})
				}
				return nil
			}

			if str, err := expr.CastToString(v); err != nil {
				return err
			} else {
				out = out.Set(&types.RecordValue{Name: k, Value: str})
				return nil
			}
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.Record:
		return val.Values, nil
	case types.RecordValueSet:
		return val, nil
	case map[string]string:
		out = types.RecordValueSet{}
		for k, v := range val {
			out = out.Set(&types.RecordValue{Name: k, Value: v})
		}

		return

	case map[string][]string:
		out = types.RecordValueSet{}
		for k, vv := range val {
			for i, v := range vv {
				out = out.Set(&types.RecordValue{Name: k, Value: v, Place: uint(i)})
			}
		}

		return

	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func (t *ComposeRecordValues) AssignFieldValue(pp []string, val interface{}) error {
	return assignToComposeRecordValues(&t.value, pp, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Record's underlying value (*types.RecordValues)
// and it's fields
//
func (t ComposeRecordValues) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return composeRecordValuesGValSelector(t.value, k)
}

// Select is field accessor for *types.Record
//
// Similar to SelectGVal but returns typed values
func (t ComposeRecordValues) Select(k string) (expr.TypedValue, error) {
	return composeRecordValuesTypedValueSelector(t.value, k)
}

func (t ComposeRecordValues) Has(k string) bool {
	return t.value.Get(k, 0) != nil
}

// recordGValSelector is field accessor for *types.RecordValueSet
func composeRecordValuesGValSelector(res types.RecordValueSet, k string) (interface{}, error) {
	vv := res.FilterByName(k)

	switch len(vv) {
	case 0:
		return nil, nil
	case 1:
		return vv[0].Value, nil
	default:
		out := make([]string, 0, len(vv))
		return out, vv.Walk(func(v *types.RecordValue) error {
			out = append(out, v.Value)
			return nil
		})
	}
}

// recordValuesTypedValueSelector is field accessor for *types.RecordValueSet
//
// @todo return appropriate types (atm all values are returned as String)
func composeRecordValuesTypedValueSelector(res types.RecordValueSet, k string) (expr.TypedValue, error) {
	vv := res.FilterByName(k)

	switch {
	case len(vv) == 0:
		return nil, nil
	case len(vv) == 1:
		return expr.NewString(vv[0].Value)
	default:
		mval := make([]expr.TypedValue, 0, len(vv))
		_ = vv.Walk(func(v *types.RecordValue) error {
			mval = append(mval, expr.Must(expr.NewString(v.Value)))
			return nil
		})

		return expr.NewArray(mval)
	}
}

// assignToRecordValuesSet is field value setter for *types.Record
func assignToComposeRecordValues(res *types.RecordValueSet, pp []string, val interface{}) (err error) {
	if len(pp) < 1 {
		switch val := expr.UntypedValue(val).(type) {
		case types.RecordValueSet:
			*res = val
			return
		case *types.Record:
			*res = val.Values
			return
		}

		return fmt.Errorf("empty path used for assigning record values")
	}

	k := pp[0]
	rv := &types.RecordValue{Name: k}

	// @todo this needs to be implemented properly
	//       we're just guessing here and puting out fires
	switch utval := expr.UntypedValue(val).(type) {
	case time.Time:
		rv.Value = utval.Format(time.RFC3339)
	case *time.Time:
		rv.Value = utval.Format(time.RFC3339)
	case []expr.TypedValue: // expr.Array
		// Handle situation where array of values is assigned to a single (multi-value) field
		// @todo this should use field context (when available) to determinate if we're actually
		//       setting array to a multi-value field

		if len(pp) == 2 {
			// Tying to assign an array of values to a single value; that will not work
			return fmt.Errorf("can not assign array of values to a single value in a record value set")
		}

		for p, v := range utval {
			rv = &types.RecordValue{Name: k, Place: uint(p)}
			rv.Value, err = cast.ToStringE(v.Get())
			if err != nil {
				return
			}

			*res = res.Set(rv)
		}

		return nil
	case []interface{}: // expr.Any
		// Handle situation where array of values is assigned to a single (multi-value) field
		// @todo this should use field context (when available) to determinate if we're actually
		//       setting array to a multi-value field

		if len(pp) == 2 {
			// Tying to assign an array of values to a single value; that will not work
			return fmt.Errorf("can not assign array of values to a single value in a record value set")
		}

		for p, v := range utval {
			rv = &types.RecordValue{Name: k, Place: uint(p)}
			rv.Value, err = cast.ToStringE(v)
			if err != nil {
				return
			}

			*res = res.Set(rv)
		}

		return nil
	default:
		rv.Value, err = cast.ToStringE(utval)
	}

	if err != nil {
		return
	}

	if len(pp) == 2 {
		if rv.Place, err = cast.ToUintE(expr.UntypedValue(pp[1])); err != nil {
			return fmt.Errorf("failed to decode record value place from '%s': %w", strings.Join(pp, "."), err)
		}
	}

	*res = res.Set(rv)

	return nil
}

func CastToComposeRecordValueErrorSet(val interface{}) (out *types.RecordValueErrorSet, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *types.RecordValueErrorSet:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}
