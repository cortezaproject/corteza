package expr

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"reflect"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/spf13/cast"
)

func (t *Vars) Len() (out int) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	return len(t.value)
}

func (t *Vars) Select(k string) (out TypedValue, err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if v, is := t.value[k]; is {
		return v, nil
	} else {
		return nil, errors.NotFound("no such key '%s'", k)
	}
}

func (t *Vars) AssignFieldValue(key string, val TypedValue) (err error) {
	if t == nil {
		return
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	if t.value == nil {
		t.value = make(map[string]TypedValue)
	}

	t.value[key] = val

	return
}

func (t *Vars) ResolveTypes(res func(typ string) Type) (err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	for k, v := range t.value {
		if u, is := v.(*Unresolved); is {
			if res(u.Type()) == nil {
				return errors.NotFound("failed to resolve unknown or unregistered type %q on %q", u.Type(), k)
			}

			t.value[k], err = res(u.Type()).Cast(t.value[k])
			if err != nil {
				return fmt.Errorf("failed to resolve: %w", err)
			}
		}

		if r, is := t.value[k].(resolvableType); is {
			if err = r.ResolveTypes(res); err != nil {
				return
			}
		}
	}

	return
}

// Merge combines the given Vars(es) into Vars
// NOTE: It will return CLONE of the original Vars, if it's called without any parameters
func (t *Vars) Merge(nn ...Iterator) (out TypedValue, err error) {
	return t.MustMerge(nn...), nil
}

// MustMerge returns Vars after merging the given Vars(es) into it
func (t *Vars) MustMerge(nn ...Iterator) *Vars {
	if t != nil {
		t.mux.RLock()
		defer t.mux.RUnlock()

		nn = append([]Iterator{t}, nn...)
	}

	var (
		out = &Vars{value: make(map[string]TypedValue)}
	)

	for _, i := range nn {
		_ = i.Each(func(k string, v TypedValue) error {
			out.value[k] = v
			return nil
		})
	}

	return out
}

// Copy takes base variables and copies all to dst
func (t *Vars) Copy(dst *Vars, kk ...string) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if t == nil {
		return
	}

	if dst.value == nil {
		dst.value = make(map[string]TypedValue)
	}

	for _, k := range kk {
		dst.value[k] = t.value[k]
	}
}

// Has returns true key is present
func (t *Vars) Has(key string) bool {
	return t.HasAll(key)
}

// HasAll returns true if all keys are present
func (t *Vars) HasAll(key string, kk ...string) bool {
	if t == nil {
		return false
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if t == nil {
		return false
	}

	for _, key = range append([]string{key}, kk...) {
		if _, has := t.value[key]; !has {
			return false
		}
	}

	return true
}

// HasAny returns true if all keys are present
func (t *Vars) HasAny(key string, kk ...string) bool {
	if t == nil {
		return false
	}

	for _, key = range append([]string{key}, kk...) {
		if _, has := t.value[key]; has {
			return true
		}
	}

	return false
}

var _ gval.Selector = &Vars{}

func (t *Vars) Dict() (out map[string]interface{}) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	out = make(map[string]interface{})
	for k, v := range t.value {
		switch v := v.(type) {
		case Dict:
			out[k] = v.Dict()

		case Slice:
			out[k] = v.Slice()

		case TypedValue:
			tmp := v.Get()
			if d, is := tmp.(Dict); is {
				out[k] = d.Dict()
			} else {
				out[k] = tmp
			}

		default:
			out[k] = v
		}
	}

	return
}

func (t *Vars) Decode(dst interface{}) (err error) {
	if t == nil {
		return nil
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	dstRef := reflect.ValueOf(dst)

	if dstRef.Kind() != reflect.Ptr {
		return fmt.Errorf("expecting a pointer, not a value")
	}

	if dstRef.IsNil() {
		return fmt.Errorf("nil pointer passed")
	}

	dstRef = dstRef.Elem()

	for i := 0; i < dstRef.NumField(); i++ {
		var (
			value TypedValue
			has   bool
			ftyp  = dstRef.Type().Field(i)
		)

		keyName := ftyp.Tag.Get("var")
		if keyName == "" {
			keyName = strings.ToLower(ftyp.Name[:1]) + ftyp.Name[1:]
		}

		value, has = t.value[keyName]
		if !has {
			continue
		}

		if tvd, is := value.(TypeValueDecoder); is {
			if err = tvd.Decode(dstRef.Field(i)); err != nil {
				return
			}
		} else if err = decode(dstRef.Field(i), value); err != nil {
			return fmt.Errorf("failed to decode value to field %s: %w", ftyp.Name, err)
		}
	}

	return
}

func (t *Vars) Scan(src any) error { return sql.ParseJSON(src, t) }

func (t *Vars) Value() (driver.Value, error) {
	if t != nil {
		t.mux.RLock()
		defer t.mux.RUnlock()
	}

	return json.Marshal(t)
}

func (t *Vars) SelectGVal(_ context.Context, k string) (out interface{}, err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	out, err = t.Select(k)
	switch c := out.(type) {
	case gval.Selector:
		return c, err
	default:
		return UntypedValue(out), err
	}
}

// UnmarshalJSON unmarshal JSON value into Vars
func (t *Vars) UnmarshalJSON(in []byte) (err error) {
	if t == nil {
		return
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	var (
		aux = make(map[string]*typedValueWrap)
	)

	if t.value == nil {
		t.value = make(map[string]TypedValue)
	}

	if len(in) == 0 {
		return nil
	}

	if err = json.Unmarshal(in, &aux); err != nil {
		return
	}

	for k, v := range aux {
		if t.value[k], err = NewUnresolved(v.Type, v.Value); err != nil {
			return
		}
	}

	return
}

func (t *Vars) Each(fn func(k string, v TypedValue) error) (err error) {
	if t == nil {
		return nil
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if t.value == nil {
		return
	}

	for k, v := range t.value {
		if err = fn(k, v); err != nil {
			return
		}
	}

	return
}

// Set or update the specific key value in Vars
func (t *Vars) Set(k string, v interface{}) (err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if t.value == nil {
		t.value = make(map[string]TypedValue)
	}

	t.value[k], err = Typify(v)
	return
}

// MarshalJSON returns JSON encoding of expression
func (t *Vars) MarshalJSON() (out []byte, err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	aux := make(map[string]*typedValueWrap)
	for k, v := range t.value {
		if v == nil {
			continue
		}

		aux[k] = &typedValueWrap{Type: v.Type()}

		if _, is := v.(json.Marshaler); is {
			aux[k].Value = v
		} else {
			aux[k].Value = v.Get()
		}
	}

	return json.Marshal(aux)
}

func decode(dst reflect.Value, src TypedValue) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		switch r := r.(type) {
		case error:
			err = r
		default:
			err = fmt.Errorf("%v", r)
		}
	}()

	if dst.Kind() == reflect.Interface && reflect.ValueOf(src).Type().Implements(dst.Type()) {
		dst.Set(reflect.ValueOf(src))
		return
	}

	if reflect.ValueOf(src).Type().ConvertibleTo(dst.Type()) {
		dst.Set(reflect.ValueOf(src))
		return
	}

	raw := UntypedValue(src)
	// Optimistically try to decode source to destination by comparing (internal) value type for destination
	if reflect.ValueOf(raw).Type().ConvertibleTo(dst.Type()) {
		dst.Set(reflect.ValueOf(raw))
		return
	}

	var (
		vBool    bool
		vInt64   int64
		vUint64  uint64
		vFloat64 float64
		vString  string
	)

	switch dst.Kind() {
	case reflect.Bool:
		if vBool, err = cast.ToBoolE(raw); err == nil {
			dst.SetBool(vBool)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if vInt64, err = cast.ToInt64E(raw); err == nil {
			dst.SetInt(vInt64)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if vUint64, err = cast.ToUint64E(raw); err == nil {
			dst.SetUint(vUint64)
		}

	case reflect.Float32, reflect.Float64:
		if vFloat64, err = cast.ToFloat64E(raw); err == nil {
			dst.SetFloat(vFloat64)
		}

	case reflect.String:
		if vString, err = cast.ToStringE(raw); err == nil {
			dst.SetString(vString)
		}

	case reflect.Map:
		dst.Set(reflect.ValueOf(src.Get()))

	//case reflect.Interface:
	//	dst.Set(reflect.ValueOf(raw))

	default:
		return fmt.Errorf("failed to cast %T to %s", raw, dst.Kind())
	}

	if err != nil {
		return fmt.Errorf("failed to cast %T to %s: %w", raw, dst.Kind(), err)
	}

	return nil
}

func CastToMeta(val interface{}) (out map[string]any, err error) {
	val = UntypedValue(val)

	if val == nil {
		return make(map[string]any), nil
	}

	switch c := val.(type) {
	case *Vars:
		c.mux.RLock()
		defer c.mux.RUnlock()

		out = make(map[string]any)
		for k, v := range c.value {
			out[k] = v.Get()
		}

		return
	case map[string]string:
		out = make(map[string]any)
		for k, v := range c {
			out[k] = v
		}

		return
	case map[string][]string:
		out = make(map[string]any)
		for k, v := range c {
			out[k] = v
		}

		return
	case KV:
		out = make(map[string]any)
		for k, v := range c.value {
			out[k] = v
		}

		return
	case KVV:
		out = make(map[string]any)
		for k, v := range c.value {
			out[k] = v
		}

		return
	case map[string]any:
		return c, nil
	}

	return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
}
func CastToVars(val interface{}) (out map[string]TypedValue, err error) {
	val = UntypedValue(val)

	if val == nil {
		return make(map[string]TypedValue), nil
	}

	switch c := val.(type) {
	case *Vars:
		c.mux.RLock()
		defer c.mux.RUnlock()

		return c.value, nil
	case map[string]TypedValue:
		return c, nil
	case map[string]interface{}:
		out = make(map[string]TypedValue)
		for k, v := range c {
			out[k], err = Typify(v)
			if err != nil {
				return
			}
		}

		return
	}

	return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
}

// Filter take keys returns Vars with only those key value pair
func (t *Vars) Filter(keys ...string) (out TypedValue, err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if t.value == nil {
		return
	}
	vars := EmptyVars()

	for _, k := range keys {
		_, has := t.value[k]
		if has {
			vars.value[k] = t.value[k]
		}
	}

	return vars, nil
}

// Delete take keys returns Vars without those key value pair
func (t *Vars) Delete(keys ...string) (out TypedValue, err error) {
	if t == nil {
		return
	}

	t.mux.RLock()
	defer t.mux.RUnlock()

	if t.value == nil {
		return
	}

	// get cloned Vars
	out, err = t.Merge()
	if err != nil {
		return
	}

	vars := out.(*Vars)

	// Delete key from t.value if exist
	for _, k := range keys {
		delete(vars.value, k)
	}

	return vars, nil
}
