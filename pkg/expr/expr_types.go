package expr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/spf13/cast"
)

type (
	resolvableType interface {
		Type
		ResolveTypes(func(string) Type) error
	}

	merger interface {
		Merge(...Iterator) (TypedValue, error)
	}

	filterer interface {
		Filter(...string) (TypedValue, error)
	}

	deleter interface {
		Delete(...string) (TypedValue, error)
	}
)

func KvFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("set", set),
		gval.Function("merge", merge),
		gval.Function("filter", filter),
		gval.Function("omit", omit),
	}
}

func EmptyKV() *KV {
	return &KV{value: make(map[string]string)}
}

func EmptyKVV() *KVV {
	return &KVV{value: make(map[string][]string)}
}

func EmptyVars() *Vars {
	return &Vars{value: make(map[string]TypedValue)}
}

func ResolveTypes(rt resolvableType, resolver func(typ string) Type) error {
	return rt.ResolveTypes(resolver)
}

func set(m merger, key string, val TypedValue) (out TypedValue, err error) {
	out, err = m.Merge()
	if err != nil {
		return
	}

	err = Assign(out, key, val)
	if err != nil {
		return
	}

	return
}

func merge(m merger, nn ...Iterator) (TypedValue, error) {
	return m.Merge(nn...)
}

func filter(f filterer, ss ...string) (TypedValue, error) {
	return f.Filter(ss...)
}

func omit(d deleter, ss ...string) (TypedValue, error) {
	return d.Delete(ss...)
}

// Typify detects input type and wraps it with expression type
func Typify(in interface{}) (tv TypedValue, err error) {
	var is bool
	if tv, is = in.(TypedValue); is {
		return
	}

	switch c := in.(type) {
	case []TypedValue:
		return &Array{value: c}, nil
	case bool:
		return &Boolean{value: c}, nil
	case uint8:
		return &UnsignedInteger{value: uint64(c)}, nil
	case uint16:
		return &UnsignedInteger{value: uint64(c)}, nil
	case uint32:
		return &UnsignedInteger{value: uint64(c)}, nil
	case uint64:
		return &UnsignedInteger{value: c}, nil
	case int8:
		return &Integer{value: int64(c)}, nil
	case int16:
		return &Integer{value: int64(c)}, nil
	case int32:
		return &Integer{value: int64(c)}, nil
	case int64:
		return &Integer{value: c}, nil
	case float32:
		return &Float{value: float64(c)}, nil
	case float64:
		return &Float{value: c}, nil
	case string:
		return &String{value: c}, nil
	case []byte:
		return &String{value: string(c)}, nil
	case *time.Time:
		return &DateTime{value: c}, nil
	case time.Time:
		return &DateTime{value: &c}, nil
	case *time.Duration:
		return &Duration{value: *c}, nil
	case time.Duration:
		return &Duration{value: c}, nil
	case map[string]interface{}:
		if v, err := CastToVars(c); err != nil {
			return nil, err
		} else {
			return &Vars{value: v}, nil
		}
	case map[string]TypedValue:
		return &Vars{value: c}, nil
	case map[string]string:
		return &KV{value: c}, nil
	case map[string][]string:
		return &KVV{value: c}, nil
	case io.Reader, io.ReadCloser, io.ReadSeeker, io.ReadSeekCloser, io.ReadWriteSeeker:
		return &Reader{value: c.(io.Reader)}, nil
	default:
		return &Any{value: c}, nil
	}
}

// Unresolved is a special type that holds value + type it needs to be resolved to
//
// This solves problem with typed value serialization
type Unresolved struct {
	typ   string
	value interface{}
}

// NewUnresolved creates new instance of Unresolved expression type
func NewUnresolved(typ string, val interface{}) (TypedValue, error) {
	return &Unresolved{
		typ:   typ,
		value: UntypedValue(val),
	}, nil
}

// Returns underlying value on Unresolved
func (t Unresolved) Get() interface{} { return t.value }

// Returns type name
func (t Unresolved) Type() string { return t.typ }

// Casts value to interface{}
func (Unresolved) Cast(interface{}) (TypedValue, error) {
	return nil, fmt.Errorf("cannot cast to unresolved type")
}

func (t *Unresolved) Assign(interface{}) (err error) {
	return fmt.Errorf("cannot set on unresolved type")
}

func CastToAny(val interface{}) (interface{}, error) {
	return val, nil
}

func CastToArray(val interface{}) (out []TypedValue, err error) {

	switch val := val.(type) {
	case nil:
		return make([]TypedValue, 0), nil
	case *Array:
		return val.value, nil
	}

	cast := func(val interface{}) (out []TypedValue, err error) {
		ref := reflect.ValueOf(val)
		if ref.Kind() == reflect.Slice {
			out = make([]TypedValue, ref.Len())
			for i := 0; i < ref.Len(); i++ {
				item := ref.Index(i).Interface()
				out[i], err = Typify(item)
				if err != nil {
					return
				}
			}
		}
		return
	}

	if c, ok := val.(TypedValue); ok {
		out, err = cast(c.Get())
	} else {
		out, err = cast(val)
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("unable to cast %T to []TypedValue", val)
	}

	return
}

var _ TypeValueDecoder = &Array{}

func (t Array) MarshalJSON() ([]byte, error) {
	var (
		aux = make([]*typedValueWrap, len(t.value))
	)

	for i, v := range t.value {
		aux[i] = &typedValueWrap{Type: v.Type()}

		if _, is := v.(json.Marshaler); is {
			aux[i].Value = v
		} else {
			aux[i].Value = v.Get()
		}
	}

	return json.Marshal(aux)
}

func (t *Array) Decode(dst reflect.Value) error {
	if dst.Kind() != reflect.Slice {
		return fmt.Errorf("failed to decode Array to non-slice")
	}

	if reflect.ValueOf(make([]TypedValue, 0)).Type() == dst.Type() {
		dst.Set(reflect.ValueOf(t.value))
		return nil
	}

	out := reflect.MakeSlice(dst.Type(), len(t.value), len(t.value))
	for i := range t.value {
		out.Index(i).Set(reflect.ValueOf(UntypedValue(t.value[i])))
	}

	dst.Set(out)

	return nil
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Record's underlying value (*types.Array)
// and it's fields
//
func (t Array) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	if s, err := t.Select(k); err != nil {
		return nil, err
	} else {
		return UntypedValue(s), nil
	}
}

var _ FieldSelector = &Array{}

// Select is field accessor for *types.Record
//
// Similar to SelectGVal but returns typed values
func (t Array) Has(k string) bool {
	if i, err := cast.ToIntE(k); err != nil {
		return false
	} else {
		return i >= 0 && i < len(t.value)
	}
}

// Push appends value to array
func (t *Array) Push(v TypedValue) {
	t.value = append(t.value, v)
}

func (t *Array) Slice() []interface{} {
	rr := make([]interface{}, len(t.GetValue()))
	for i, v := range t.GetValue() {
		switch v := v.(type) {
		case Dict:
			rr[i] = v.Dict()

		case Slice:
			rr[i] = v.Slice()

		case TypedValue:
			tmp := v.Get()
			if d, is := tmp.(Dict); is {
				rr[i] = d.Dict()
			} else {
				rr[i] = tmp
			}

		default:
			rr[i] = v
		}
	}

	return rr
}

// Select is field accessor for *types.Array
//
// Similar to SelectGVal but returns typed values
func (t Array) Select(k string) (TypedValue, error) {
	if i, err := cast.ToIntE(k); err != nil {
		return nil, err
	} else {
		return t.value[i], nil
	}
}

// emptyStringFailsafe returns 0 on empty strings
func emptyStringFailsafe(val interface{}) interface{} {
	val = UntypedValue(val)
	if aux, is := val.(string); is && len(strings.TrimSpace(aux)) == 0 {
		return 0
	} else {
		return val
	}
}

func (t ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%d", t.value))
}

func CastToBoolean(val interface{}) (out bool, err error) {
	val = UntypedValue(val)
	if aux, is := val.(string); is && len(strings.TrimSpace(aux)) == 0 {
		return false, nil
	}

	return cast.ToBoolE(val)
}

func CastToString(val interface{}) (out string, err error) {
	switch v := val.(type) {
	case io.Reader:
		bb, err := ioutil.ReadAll(v)
		if err != nil {
			return "", err
		}
		return string(bb), nil
	default:
		return cast.ToStringE(UntypedValue(val))
	}
}

func CastToBytes(val interface{}) (out []byte, err error) {
	s, err := CastToString(val)

	if err != nil {
		return
	}

	out = []byte(s)
	return
}

func CastToStringSlice(val interface{}) (out []string, err error) {
	return cast.ToStringSliceE(UntypedValue(val))
}

func CastToHandle(val interface{}) (string, error) {
	val = UntypedValue(val)

	h, err := cast.ToStringE(val)

	if !handle.IsValid(h) {
		return "", fmt.Errorf("invalid handle format: '%s'", h)
	}

	return h, err
}

func CastToDuration(val interface{}) (out time.Duration, err error) {
	return cast.ToDurationE(emptyStringFailsafe(val))
}

func CastToDateTime(val interface{}) (out *time.Time, err error) {
	val = UntypedValue(val)
	switch casted := val.(type) {
	case *time.Time:
		return casted, nil
	case time.Time:
		return &casted, nil
	default:
		var c time.Time
		if c, err = cast.ToTimeE(casted); err != nil {
			return nil, err
		}

		return &c, nil
	}
}

func CastToFloat(val interface{}) (out float64, err error) {
	val = UntypedValue(val)
	switch val := val.(type) {
	case nil:
		return 0, nil
	case *Float:
		return val.value, nil
	}

	return cast.ToFloat64E(emptyStringFailsafe(val))
}

func CastToID(val interface{}) (out uint64, err error) {
	return cast.ToUint64E(emptyStringFailsafe(val))
}

func CastToInteger(val interface{}) (out int64, err error) {
	return cast.ToInt64E(emptyStringFailsafe(val))
}

func CastToUnsignedInteger(val interface{}) (out uint64, err error) {
	return cast.ToUint64E(emptyStringFailsafe(val))
}

func (t *KV) Has(k string) bool {
	_, has := t.value[k]
	return has
}

func (t *KV) Select(k string) (TypedValue, error) {
	if v, has := t.value[k]; has {
		return Must(NewString(v)), nil
	} else {
		return nil, errors.NotFound("no such key '%s'", k)
	}
}

func (t *KV) AssignFieldValue(key string, val TypedValue) error {
	return assignToKV(t, key, val)
}

func assignToKV(t *KV, key string, val TypedValue) error {
	if t.value == nil {
		t.value = make(map[string]string)
	}

	str, err := cast.ToStringE(UntypedValue(val))
	t.value[key] = str
	return err
}

func CastToKV(val interface{}) (out map[string]string, err error) {
	val = UntypedValue(val)

	if val == nil {
		return make(map[string]string), nil
	}

	switch casted := val.(type) {
	case map[string]string:
		return casted, nil
	default:
		return cast.ToStringMapStringE(UntypedValue(casted))
	}
}

func (t *KV) Each(fn func(k string, v TypedValue) error) (err error) {
	if t == nil || t.value == nil {
		return
	}

	for k, v := range t.value {
		var val TypedValue
		val, err = Typify(v)
		if err != nil {
			return err
		}

		if err = fn(k, val); err != nil {
			return
		}
	}

	return
}

// Merge combines the given KVs into KV
// NOTE: It will return CLONE of the original KV, if its called without any parameters
func (t *KV) Merge(nn ...Iterator) (out TypedValue, err error) {
	kv := EmptyKV()

	nn = append([]Iterator{t}, nn...)

	for _, i := range nn {
		err = i.Each(func(k string, v TypedValue) error {
			str, err := cast.ToStringE(v.Get())
			kv.value[k] = str
			return err
		})
		if err != nil {
			return
		}
	}

	return kv, nil
}

// Filter take keys returns KV with only those key value pair
func (t *KV) Filter(keys ...string) (out TypedValue, err error) {
	if t.value == nil {
		return
	}
	kv := EmptyKV()

	for _, k := range keys {
		_, has := t.value[k]
		if has {
			kv.value[k] = t.value[k]
		}
	}

	return kv, nil
}

// Delete take keys returns KV without those key value pair
func (t *KV) Delete(keys ...string) (out TypedValue, err error) {
	if t.value == nil {
		return
	}

	// get cloned KV
	out, err = t.Merge()
	if err != nil {
		return
	}

	kv := out.(*KV)

	// Delete key from out.value if exist
	for _, k := range keys {
		delete(kv.value, k)
	}

	return kv, nil
}

func (t *KVV) AssignFieldValue(key []string, val TypedValue) error {
	return assignToKVV(t, key, val)
}

func assignToKVV(t *KVV, kk []string, val TypedValue) error {
	if t.value == nil {
		t.value = make(map[string][]string)
	}

	switch len(kk) {
	case 2:
		str, err := cast.ToStringE(val.Get())
		if err != nil {
			return err
		}

		key, ind := kk[0], kk[1]

		if len(ind) > 0 {
			// handles kvv.field[42] = "value"
			index, err := cast.ToIntE(ind)
			if err != nil {
				return err
			}

			if index >= 0 && index < len(t.value[key]) {
				// handles positive & in-range indexes
				t.value[key][index] = str
				return nil
			}

			//negative & out-of-range indexes are always appended
		}

		// handles kvv.field[] = "value"
		t.value[key] = append(t.value[key], str)

	case 1:
		str, err := cast.ToStringSliceE(val.Get())
		if err != nil {
			return err
		}

		t.value[kk[0]] = str

	default:
		return fmt.Errorf("cannot set value on %s with path '%s'", t.Type(), strings.Join(kk, "."))
	}

	return nil
}

func CastToKVV(val interface{}) (out map[string][]string, err error) {
	val = UntypedValue(val)

	if val == nil {
		return make(map[string][]string), nil
	}

	switch casted := val.(type) {
	case http.Header:
		return casted, nil
	case url.Values:
		return casted, nil
	default:
		return cast.ToStringMapStringSliceE(casted)
	}
}

func CastToReader(val interface{}) (out io.Reader, err error) {
	val = UntypedValue(val)

	switch casted := val.(type) {
	case []byte:
		return bytes.NewReader(casted), nil
	case string:
		return strings.NewReader(casted), nil
	case io.Reader:
		return casted, nil
	default:
		return nil, fmt.Errorf("unable to cast %T to io.Reader", val)
	}
}

func (t *KVV) Each(fn func(k string, v TypedValue) error) (err error) {
	if t == nil || t.value == nil {
		return
	}

	for k, v := range t.value {
		var val TypedValue
		val, err = Typify(v)
		if err != nil {
			return err
		}

		if err = fn(k, val); err != nil {
			return
		}
	}

	return
}

// Merge combines the given KVVs into KVV
// NOTE: It will return CLONE of the original KVV, if its called without any parameters
func (t *KVV) Merge(nn ...Iterator) (out TypedValue, err error) {
	kvv := EmptyKVV()

	nn = append([]Iterator{t}, nn...)

	for _, i := range nn {
		err = i.Each(func(k string, v TypedValue) error {
			ss, err := cast.ToStringSliceE(v.Get())
			if err != nil {
				return err
			}

			_, has := kvv.value[k]
			if has {
				kvv.value[k] = append(kvv.value[k], ss...)
			} else {
				kvv.value[k] = ss
			}
			return nil
		})
		if err != nil {
			return
		}
	}

	return kvv, nil
}

// Filter take keys returns KVV with only those key value pair
func (t *KVV) Filter(keys ...string) (out TypedValue, err error) {
	if t.value == nil {
		return
	}
	kvv := EmptyKVV()

	for _, k := range keys {
		_, has := t.value[k]
		if has {
			kvv.value[k] = t.value[k]
		}
	}

	return kvv, nil
}

// Delete take keys returns KVV without those key value pair
func (t *KVV) Delete(keys ...string) (out TypedValue, err error) {
	if t.value == nil {
		return
	}

	// get cloned KVV
	out, err = t.Merge()
	if err != nil {
		return
	}

	kvv := out.(*KVV)

	// Delete key from t.value if exist
	for _, k := range keys {
		delete(kvv.value, k)
	}

	return kvv, nil
}
