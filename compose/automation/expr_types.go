package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/spf13/cast"
	"strconv"
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
	case map[string]interface{}:
		out = &types.Namespace{}
		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, out)

		return
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
	case map[string]interface{}:
		out = &types.Module{}
		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, out)
		return

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
		if val == nil {
			val = &types.Record{}
		}

		if val.Values == nil {
			val.Values = types.RecordValueSet{}
		}
		return val, nil
	case map[string]interface{}:
		out = &types.Record{}
		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, out)

		return

	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

var _ expr.DeepFieldAssigner = &ComposeRecord{}

// AssignFieldValue implements expr.DeepFieldAssigner
//
// We need to reroute value assigning for record-value-sets because
// we loose the reference to record-value slice
func (t *ComposeRecord) AssignFieldValue(kk []string, val expr.TypedValue) error {
	switch kk[0] {
	case "values":
		return assignToComposeRecordValues(t.value, kk[1:], val)
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
func (t ComposeRecord) SelectGVal(_ context.Context, k string) (interface{}, error) {
	if k == "values" {
		if t.value.Values == nil {
			t.value.Values = types.RecordValueSet{}
		}

		return t.value.Values.Dict(t.value.GetModule().Fields), nil
	}

	return composeRecordGValSelector(t.value, k)
}

// Select is field accessor for *types.ComposeRecord
//
// Similar to SelectGVal but returns typed values
func (t ComposeRecord) Select(k string) (expr.TypedValue, error) {
	if k == "values" {
		if t.value.Values == nil {
			t.value.Values = types.RecordValueSet{}
		}

		return &ComposeRecordValues{value: t.value}, nil
	}

	return composeRecordTypedValueSelector(t.value, k)
}

type ComposeRecordValues struct{ value *types.Record }

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

func (t *ComposeRecordValues) AssignFieldValue(pp []string, val expr.TypedValue) error {
	return assignToComposeRecordValues(t.value, pp, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Record's underlying value (*types.RecordValues)
// and it's fields
//
func (t *ComposeRecordValues) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return composeRecordValuesGValSelector(t.value, k)
}

// Select is field accessor for *types.Record
//
// Similar to SelectGVal but returns typed values
func (t *ComposeRecordValues) Select(k string) (expr.TypedValue, error) {
	return composeRecordValuesTypedValueSelector(t.value, k)
}

func (t ComposeRecordValues) Has(k string) bool {
	return t.value.Values.Get(k, 0) != nil
}

// recordGValSelector is field accessor for *types.RecordValueSet
func composeRecordValuesGValSelector(res *types.Record, k string) (interface{}, error) {
	var (
		vv = res.Values.FilterByName(k)

		multiValueField bool
		field           *types.ModuleField
	)

	if mod := res.GetModule(); mod != nil {
		fld := mod.Fields.FindByName(k)

		if fld == nil {
			return nil, fmt.Errorf("field '%s' does not exist on module %s", k, mod.Name)
		}

		multiValueField = fld.Multi
	}

	switch {
	case len(vv) == 0:
		if field != nil && field.IsBoolean() {
			return false, nil
		}

		return nil, nil

	case len(vv) == 1 && !multiValueField:
		return vv[0].Cast(field)

	default:
		out := make([]interface{}, 0, len(vv))
		return out, vv.Walk(func(v *types.RecordValue) error {
			i, err := v.Cast(field)
			if err != nil {
				return err
			}

			out = append(out, i)
			return nil
		})
	}
}

// recordValuesTypedValueSelector is field accessor for *types.RecordValueSet
//
// @todo return appropriate types (atm all values are returned as String)
func composeRecordValuesTypedValueSelector(res *types.Record, k string) (expr.TypedValue, error) {
	var (
		vv = res.Values.FilterByName(k)

		multiValueField bool
		field           *types.ModuleField
	)

	if mod := res.GetModule(); mod != nil {
		field = mod.Fields.FindByName(k)

		if field == nil {
			return nil, fmt.Errorf("field '%s' does not exist on module %s", k, mod.Name)
		}

		multiValueField = field.Multi
	}

	switch {
	case len(vv) == 0:
		if field != nil && field.IsBoolean() {
			return &expr.Boolean{}, nil
		}

		return nil, nil
	case len(vv) == 1 && !multiValueField:
		return recordValueToExprTypedValue(field, vv[0])
	default:
		return recordValueSetToExprArray(field, vv...)
	}
}

// assignToRecordValuesSet is field value setter for *types.RecordValueSet
//
// We'll be using types.Record for the base (and not types.RecordValueSet)
func assignToComposeRecordValues(res *types.Record, pp []string, val interface{}) (err error) {
	if len(pp) < 1 {
		switch val := expr.UntypedValue(val).(type) {
		case types.RecordValueSet:
			res.Values = val
			return
		case *types.Record:
			*res = *val
			return
		}

		return fmt.Errorf("empty path used for assigning record values")
	}

	var (
		k  = pp[0]
		rv = &types.RecordValue{Name: k}

		setSliceOfValues = func(vv []interface{}) error {
			// Handle situation where array of values is assigned to a single (multi-value) field
			// @todo this should use field context (when available) to determinate if we're actually
			//       setting array to a multi-value field

			if len(pp) == 2 {
				// Tying to assign an array of values to a single value; that will not work
				return fmt.Errorf("can not assign array of values to a single value in a record value set")
			}

			for p, v := range vv {
				rv = &types.RecordValue{Name: k, Place: uint(p)}
				rv.Value, err = cast.ToStringE(v)
				if err != nil {
					return err
				}

				res.Values = res.Values.Set(rv)
			}

			return nil
		}
	)

	// @todo this needs to be implemented properly
	//       we're just guessing here and putting out fires
	switch utval := expr.UntypedValue(val).(type) {
	case time.Time:
		rv.Value = utval.Format(time.RFC3339)
	case *time.Time:
		rv.Value = utval.Format(time.RFC3339)
	case []string:
		aux := make([]interface{}, len(utval))
		for i := range utval {
			aux[i] = utval[i]
		}

		return setSliceOfValues(aux)
	case []expr.TypedValue: // expr.Array
		aux := make([]interface{}, len(utval))
		for i := range utval {
			aux[i] = utval[i].Get()
		}

		return setSliceOfValues(aux)
	case []interface{}: // expr.Any
		return setSliceOfValues(utval)
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

	res.Values = res.Values.Set(rv)

	return nil
}

// NewComposeRecordValues creates new instance of ComposeRecordValues expression type
func NewComposeRecordValues(val interface{}) (*ComposeRecordValues, error) {
	// Try to cast to ComposeRecord first
	if rec, err := CastToComposeRecord(val); err == nil {
		return &ComposeRecordValues{value: rec}, nil
	}

	if c, err := CastToComposeRecordValues(val); err != nil {
		return nil, fmt.Errorf("unable to create ComposeRecordValues: %w", err)
	} else {
		return &ComposeRecordValues{value: &types.Record{Values: c}}, nil
	}
}

// Return underlying value on ComposeRecordValues
func (t ComposeRecordValues) Get() interface{} { return t.value }

// Return underlying value on ComposeRecordValues
func (t ComposeRecordValues) GetValue() types.RecordValueSet { return t.value.Values }

// Return type name
func (ComposeRecordValues) Type() string { return "ComposeRecordValues" }

// Convert value to types.RecordValueSet
func (ComposeRecordValues) Cast(val interface{}) (expr.TypedValue, error) {
	return NewComposeRecordValues(val)
}

// Assign new value to ComposeRecordValues
//
// value is first passed through CastToComposeRecordValues
func (t *ComposeRecordValues) Assign(val interface{}) error {
	if c, err := CastToComposeRecordValues(val); err != nil {
		return err
	} else {
		t.value.Values = c
		return nil
	}
}

func CastToComposeRecordValueErrorSet(val interface{}) (out *types.RecordValueErrorSet, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *types.RecordValueErrorSet:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func recordValueToExprTypedValue(field *types.ModuleField, rv *types.RecordValue) (expr.TypedValue, error) {
	switch {
	case field == nil:
		return expr.NewString(rv.Value)

	case field.IsRef():
		ref := rv.Ref
		if ref == 0 && len(rv.Value) > 0 {
			// Cover cases when we Ref is not set but Value is
			// This happens when RVS is transferred as JSON
			ref, _ = strconv.ParseUint(rv.Value, 10, 64)
		}

		return expr.NewID(ref)

	default:
		if v, err := rv.Cast(field); err != nil {
			return nil, err
		} else {
			return expr.Typify(v)
		}
	}
}

func recordValueSetToExprArray(field *types.ModuleField, vv ...*types.RecordValue) (arr *expr.Array, err error) {
	var (
		tv expr.TypedValue
	)

	arr = &expr.Array{}

	for _, v := range vv {
		tv, err = recordValueToExprTypedValue(field, v)
		if err != nil {
			return
		}

		arr.Push(tv)
	}

	return
}
