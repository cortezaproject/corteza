package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/slice"
	"github.com/spf13/cast"
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
func (t *ComposeRecord) AssignFieldValue(p expr.Pather, val expr.TypedValue) (err error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	switch p.Get() {
	case "values":
		err = p.Next()
		if err != nil {
			return
		}
		return assignToComposeRecordValues(t.value, p, val)
	// case "labels":
	// @todo deep setting labels
	default:
		return assignToComposeRecord(t.value, p.Get(), val)
	}
}

var _ gval.Selector = &ComposeRecord{}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Record's underlying value (*types.Record)
// and it's fields
func (t *ComposeRecord) SelectGVal(_ context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()

	if t.value != nil && k == "values" {
		if t.value.Values == nil {
			t.value.Values = types.RecordValueSet{}
		}

		return &ComposeRecordValues{t.value}, nil
	}

	return composeRecordGValSelector(t.value, k)
}

// Select is field accessor for *types.ComposeRecord
//
// Similar to SelectGVal but returns typed values
func (t *ComposeRecord) Select(k string) (expr.TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()

	if t.value != nil && k == "values" {
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
	case *types.RecordValue:
		return types.RecordValueSet{val}, nil
	case types.RecordValueSet:
		return val, nil
	case map[string]string:
		for _, k := range slice.Keys(val) {
			out = out.Set(&types.RecordValue{Name: k, Value: val[k]})
		}
		return
	case map[string][]string:
		for _, k := range slice.Keys(val) {
			for i, v := range val[k] {
				out = out.Set(&types.RecordValue{Name: k, Value: v, Place: uint(i)})
			}
		}
		return
	case map[string]interface{}:
		for _, k := range slice.Keys(val) {
			if isNil(val[k]) {
				continue
			}

			var vv []interface{}
			// covering a couple of typed slices
			ref := reflect.ValueOf(val[k])

			if ref.Kind() == reflect.Slice {
				for i := 0; i < ref.Len(); i++ {
					vv = append(vv, ref.Index(i).Interface())
				}
			} else {
				vv = []interface{}{val[k]}
			}

			i := 0
			for _, value := range vv {
				switch v := value.(type) {
				// a small exception for boolean fields, don't add false values and when its true cast to 1
				case bool:
					if !v {
						continue
					} else {
						value = "1"
					}
				}

				out = out.Set(&types.RecordValue{Name: k, Value: cast.ToString(value), Place: uint(i)})

				// explicitly counting because of how booleans are handled
				i++
			}
		}
		return
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func (t *ComposeRecordValues) AssignFieldValue(p expr.Pather, val expr.TypedValue) error {
	return assignToComposeRecordValues(t.value, p, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Record's underlying value (*types.RecordValues)
// and it's fields
func (t *ComposeRecordValues) SelectGVal(_ context.Context, k string) (interface{}, error) {
	return composeRecordValuesGValSelector(t.value, k)
}

// IsEmpty implements pkg/expr.empty requirements to be able to determine if
// the value is empty.
//
// This is needed cor cases when we are working with empty records, but are trying
// to access their values.
func (t *ComposeRecordValues) IsEmpty() bool {
	return t == nil || t.value == nil || len(t.value.Values) == 0
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
		field = mod.Fields.FindByName(k)

		if field == nil {
			return nil, fmt.Errorf("field '%s' does not exist on module %s", k, mod.Name)
		}

		multiValueField = field.Multi
	}

	switch {
	case len(vv) == 0:
		if field != nil && field.IsBoolean() {
			return false, nil
		}

		// return nil time.Time for date and time fields
		if field != nil && (field.IsDateOnly() || field.IsDateTime() || field.IsTimeOnly()) {
			return expr.NewDateTime(nil)
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
func assignToComposeRecordValues(res *types.Record, p expr.Pather, val interface{}) (err error) {
	if p == nil || !p.More() {
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
		k  = p.Get()
		rv = &types.RecordValue{Name: k}

		setSliceOfValues = func(vv []interface{}) error {
			// Handle situation where array of values is assigned to a single (multi-value) field
			// @todo this should use field context (when available) to determinate if we're actually
			//       setting array to a multi-value field

			if !p.IsLast() {
				// Tying to assign an array of values to a single value; that will not work
				return fmt.Errorf("can not assign array of values to a single value in a record value set")
			}

			if ss, err := cast.ToStringSliceE(vv); err != nil {
				return err
			} else {
				res.Values = res.Values.Replace(k, ss...)
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
		if utval == nil {
			rv.Value = ""
			break
		}
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

	if !p.IsLast() {
		err = p.Next()
		if err != nil {
			return
		}

		if rv.Place, err = cast.ToUintE(expr.UntypedValue(p.Get())); err != nil {
			return fmt.Errorf("failed to decode record value place from '%s': %w", p.String(), err)
		}
	}

	res.Values = res.Values.Set(rv)

	return nil
}

// NewComposeRecordValues creates new instance of ComposeRecordValues expression type
func NewComposeRecordValues(val interface{}) (*ComposeRecordValues, error) {
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
	if v, err := rv.Cast(field); err != nil {
		return nil, err
	} else {
		return expr.Typify(v)
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

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

func CastToAttachment(val interface{}) (out *types.Attachment, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *types.Attachment:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func EmptyComposeRecordValues() *ComposeRecordValues {
	return &ComposeRecordValues{value: &types.Record{Values: types.RecordValueSet{}}}
}

func (t *ComposeRecordValues) Each(fn func(k string, v expr.TypedValue) error) (err error) {
	if t.IsEmpty() {
		return
	}

	for _, vv := range t.GetValue() {
		key := vv.Name
		if len(key) == 0 {
			continue
		}

		var val expr.TypedValue
		val, err = expr.Typify(vv)
		if err = fn(key, val); err != nil {
			return
		}
	}

	return
}

// Merge combines the given ComposeRecordValues into ComposeRecordValues
// NOTE: It will return CLONE of the original ComposeRecordValues, if it's called without any parameters
func (t *ComposeRecordValues) Merge(nn ...expr.Iterator) (out expr.TypedValue, err error) {
	rv := EmptyComposeRecordValues()

	nn = append([]expr.Iterator{t}, nn...)

	for _, i := range nn {
		err = i.Each(func(k string, v expr.TypedValue) error {
			var (
				rVal types.RecordValue
				bb   []byte
			)

			// @todo implement casting of RecordValue from TypedValue
			bb, err = json.Marshal(v.Get())
			if err != nil {
				return err
			}

			err = json.Unmarshal(bb, &rVal)
			if err != nil {
				return err
			}

			rr := rv.GetValue().Set(&rVal)
			err = rv.Assign(rr)
			if err != nil {
				return err
			}

			return err
		})
		if err != nil {
			return
		}
	}

	return rv, nil
}

func (t *ComposeRecordValues) Filter(keys ...string) (out expr.TypedValue, err error) {
	if t.IsEmpty() {
		return
	}

	// get cloned ComposeRecordValues
	out, err = t.Merge()
	if err != nil {
		return
	}

	rv := out.(*ComposeRecordValues)
	if !rv.IsEmpty() {
		rv.value.Values = types.RecordValueSet{}
	}

	keyMap := make(map[string]string)
	for _, k := range keys {
		keyMap[k] = k
	}

	// Push the only with values with matching Name
	for _, val := range t.GetValue() {
		_, ok := keyMap[val.Name]
		if ok {
			rr := rv.GetValue().Set(val.Clone())
			err = rv.Assign(rr)
			if err != nil {
				return out, err
			}
		}
	}

	return rv, nil
}

func (t *ComposeRecordValues) Delete(keys ...string) (out expr.TypedValue, err error) {
	if t.IsEmpty() {
		return
	}

	// get cloned ComposeRecordValues
	out, err = t.Merge()
	if err != nil {
		return
	}

	rv := out.(*ComposeRecordValues)
	if !rv.IsEmpty() {
		rv.value.Values = types.RecordValueSet{}
	}

	keyMap := make(map[string]string)
	for _, k := range keys {
		keyMap[k] = k
	}

	// Push the only with values with non-matching Name
	for _, val := range t.GetValue() {
		_, ok := keyMap[val.Name]
		if !ok {
			rr := rv.GetValue().Set(val.Clone())
			err = rv.Assign(rr)
			if err != nil {
				return out, err
			}
		}
	}

	return rv, nil
}

func (v *Attachment) Clone() (expr.TypedValue, error) {
	att := *v.value
	aux, err := NewAttachment(&att)
	return aux, err
}

func (v *ComposeModule) Clone() (expr.TypedValue, error) {
	aux, err := NewComposeModule(v.value.Clone())
	return aux, err
}

func (v *ComposeNamespace) Clone() (expr.TypedValue, error) {
	aux, err := NewComposeNamespace(v.value.Clone())
	return aux, err
}

func (v *ComposeRecord) Clone() (expr.TypedValue, error) {
	aux, err := NewComposeRecord(v.value.Clone())
	return aux, err
}

func (v *ComposeRecordValues) Clone() (expr.TypedValue, error) {
	aux, err := NewComposeRecordValues(v.value.Clone())
	return aux, err
}

func (v *ComposeRecordValueErrorSet) Clone() (expr.TypedValue, error) {
	if v.value == nil {
		return nil, nil
	}

	errs := types.RecordValueErrorSet{
		Set: make([]types.RecordValueError, len(v.value.Set)),
	}

	for i, v := range v.value.Set {
		errs.Set[i] = v
	}

	aux, err := NewComposeRecordValueErrorSet(&errs)
	return aux, err
}
