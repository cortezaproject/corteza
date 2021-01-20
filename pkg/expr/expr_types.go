package expr

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type (
	resolvableType interface {
		Type
		ResolveTypes(func(string) Type) error
	}
)

func ResolveTypes(rt resolvableType, resolver func(typ string) Type) error {
	return rt.ResolveTypes(resolver)
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
	return nil, fmt.Errorf("can not cast to unresolved type")
}

func (t *Unresolved) Assign(interface{}) (err error) {
	return fmt.Errorf("can not set on unresolved type")
}

func CastToAny(val interface{}) (interface{}, error) {
	return val, nil
}

func CastToArray(val interface{}) ([]TypedValue, error) {

	switch val := val.(type) {
	case *Array:
		return val.value, nil
	}

	ref := reflect.ValueOf(val)
	if ref.Kind() == reflect.Slice {
		out := make([]TypedValue, ref.Len())
		for i := 0; i < ref.Len(); i++ {
			item := ref.Index(i).Interface()
			if tVal, is := item.(TypedValue); is {
				out[i] = tVal
			} else {
				out[i] = &Any{value: item}
			}
		}
		return out, nil
	}

	return nil, fmt.Errorf("unable to cast %T to []TypedValue", val)
}

var _ TypeValueDecoder = &Array{}

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

// Select is field accessor for *types.Record
//
// Similar to SelectGVal but returns typed values
func (t Array) Select(k string) (TypedValue, error) {
	if i, err := cast.ToIntE(k); err != nil {
		return nil, err
	} else {
		return t.value[i], nil
	}
}

func CastToBoolean(val interface{}) (out bool, err error) {
	return cast.ToBoolE(UntypedValue(val))
}

func CastToString(val interface{}) (out string, err error) {
	return cast.ToStringE(UntypedValue(val))
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
	return cast.ToDurationE(UntypedValue(val))
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
	return cast.ToFloat64E(UntypedValue(val))
}

func CastToID(val interface{}) (out uint64, err error) {
	out, err = cast.ToUint64E(UntypedValue(val))
	if out == 0 {
		err = fmt.Errorf("invalid ID")
	}

	return
}

func CastToInteger(val interface{}) (out int64, err error) {
	return cast.ToInt64E(UntypedValue(val))
}

func CastToUnsignedInteger(val interface{}) (out uint64, err error) {
	return cast.ToUint64E(UntypedValue(val))
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

func (t *KV) AssignFieldValue(key string, val interface{}) error {
	return assignToKV(t, key, val)
}

func assignToKV(t *KV, key string, val interface{}) error {
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
		return cast.ToStringMapStringE(casted)
	}
}

func (t *KVV) AssignFieldValue(key string, val interface{}) error {
	return assignToKVV(t, key, val)
}

func assignToKVV(t *KVV, key string, val interface{}) error {
	if t.value == nil {
		t.value = make(map[string][]string)
	}

	str, err := cast.ToStringSliceE(val)
	t.value[key] = str
	return err
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
