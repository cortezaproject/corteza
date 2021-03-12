package expr

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/expr/expr_types.yaml

import (
	"context"
	"fmt"
	"io"
	"time"
)

var _ = context.Background
var _ = fmt.Errorf

// Any is an expression type, wrapper for interface{} type
type Any struct{ value interface{} }

// NewAny creates new instance of Any expression type
func NewAny(val interface{}) (*Any, error) {
	if c, err := CastToAny(val); err != nil {
		return nil, fmt.Errorf("unable to create Any: %w", err)
	} else {
		return &Any{value: c}, nil
	}
}

// Return underlying value on Any
func (t Any) Get() interface{} { return t.value }

// Return underlying value on Any
func (t Any) GetValue() interface{} { return t.value }

// Return type name
func (Any) Type() string { return "Any" }

// Convert value to interface{}
func (Any) Cast(val interface{}) (TypedValue, error) {
	return NewAny(val)
}

// Assign new value to Any
//
// value is first passed through CastToAny
func (t *Any) Assign(val interface{}) error {
	if c, err := CastToAny(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Array is an expression type, wrapper for []TypedValue type
type Array struct{ value []TypedValue }

// NewArray creates new instance of Array expression type
func NewArray(val interface{}) (*Array, error) {
	if c, err := CastToArray(val); err != nil {
		return nil, fmt.Errorf("unable to create Array: %w", err)
	} else {
		return &Array{value: c}, nil
	}
}

// Return underlying value on Array
func (t Array) Get() interface{} { return t.value }

// Return underlying value on Array
func (t Array) GetValue() []TypedValue { return t.value }

// Return type name
func (Array) Type() string { return "Array" }

// Convert value to []TypedValue
func (Array) Cast(val interface{}) (TypedValue, error) {
	return NewArray(val)
}

// Assign new value to Array
//
// value is first passed through CastToArray
func (t *Array) Assign(val interface{}) error {
	if c, err := CastToArray(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Boolean is an expression type, wrapper for bool type
type Boolean struct{ value bool }

// NewBoolean creates new instance of Boolean expression type
func NewBoolean(val interface{}) (*Boolean, error) {
	if c, err := CastToBoolean(val); err != nil {
		return nil, fmt.Errorf("unable to create Boolean: %w", err)
	} else {
		return &Boolean{value: c}, nil
	}
}

// Return underlying value on Boolean
func (t Boolean) Get() interface{} { return t.value }

// Return underlying value on Boolean
func (t Boolean) GetValue() bool { return t.value }

// Return type name
func (Boolean) Type() string { return "Boolean" }

// Convert value to bool
func (Boolean) Cast(val interface{}) (TypedValue, error) {
	return NewBoolean(val)
}

// Assign new value to Boolean
//
// value is first passed through CastToBoolean
func (t *Boolean) Assign(val interface{}) error {
	if c, err := CastToBoolean(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// DateTime is an expression type, wrapper for *time.Time type
type DateTime struct{ value *time.Time }

// NewDateTime creates new instance of DateTime expression type
func NewDateTime(val interface{}) (*DateTime, error) {
	if c, err := CastToDateTime(val); err != nil {
		return nil, fmt.Errorf("unable to create DateTime: %w", err)
	} else {
		return &DateTime{value: c}, nil
	}
}

// Return underlying value on DateTime
func (t DateTime) Get() interface{} { return t.value }

// Return underlying value on DateTime
func (t DateTime) GetValue() *time.Time { return t.value }

// Return type name
func (DateTime) Type() string { return "DateTime" }

// Convert value to *time.Time
func (DateTime) Cast(val interface{}) (TypedValue, error) {
	return NewDateTime(val)
}

// Assign new value to DateTime
//
// value is first passed through CastToDateTime
func (t *DateTime) Assign(val interface{}) error {
	if c, err := CastToDateTime(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Duration is an expression type, wrapper for time.Duration type
type Duration struct{ value time.Duration }

// NewDuration creates new instance of Duration expression type
func NewDuration(val interface{}) (*Duration, error) {
	if c, err := CastToDuration(val); err != nil {
		return nil, fmt.Errorf("unable to create Duration: %w", err)
	} else {
		return &Duration{value: c}, nil
	}
}

// Return underlying value on Duration
func (t Duration) Get() interface{} { return t.value }

// Return underlying value on Duration
func (t Duration) GetValue() time.Duration { return t.value }

// Return type name
func (Duration) Type() string { return "Duration" }

// Convert value to time.Duration
func (Duration) Cast(val interface{}) (TypedValue, error) {
	return NewDuration(val)
}

// Assign new value to Duration
//
// value is first passed through CastToDuration
func (t *Duration) Assign(val interface{}) error {
	if c, err := CastToDuration(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Float is an expression type, wrapper for float64 type
type Float struct{ value float64 }

// NewFloat creates new instance of Float expression type
func NewFloat(val interface{}) (*Float, error) {
	if c, err := CastToFloat(val); err != nil {
		return nil, fmt.Errorf("unable to create Float: %w", err)
	} else {
		return &Float{value: c}, nil
	}
}

// Return underlying value on Float
func (t Float) Get() interface{} { return t.value }

// Return underlying value on Float
func (t Float) GetValue() float64 { return t.value }

// Return type name
func (Float) Type() string { return "Float" }

// Convert value to float64
func (Float) Cast(val interface{}) (TypedValue, error) {
	return NewFloat(val)
}

// Assign new value to Float
//
// value is first passed through CastToFloat
func (t *Float) Assign(val interface{}) error {
	if c, err := CastToFloat(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Handle is an expression type, wrapper for string type
type Handle struct{ value string }

// NewHandle creates new instance of Handle expression type
func NewHandle(val interface{}) (*Handle, error) {
	if c, err := CastToHandle(val); err != nil {
		return nil, fmt.Errorf("unable to create Handle: %w", err)
	} else {
		return &Handle{value: c}, nil
	}
}

// Return underlying value on Handle
func (t Handle) Get() interface{} { return t.value }

// Return underlying value on Handle
func (t Handle) GetValue() string { return t.value }

// Return type name
func (Handle) Type() string { return "Handle" }

// Convert value to string
func (Handle) Cast(val interface{}) (TypedValue, error) {
	return NewHandle(val)
}

// Assign new value to Handle
//
// value is first passed through CastToHandle
func (t *Handle) Assign(val interface{}) error {
	if c, err := CastToHandle(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// ID is an expression type, wrapper for uint64 type
type ID struct{ value uint64 }

// NewID creates new instance of ID expression type
func NewID(val interface{}) (*ID, error) {
	if c, err := CastToID(val); err != nil {
		return nil, fmt.Errorf("unable to create ID: %w", err)
	} else {
		return &ID{value: c}, nil
	}
}

// Return underlying value on ID
func (t ID) Get() interface{} { return t.value }

// Return underlying value on ID
func (t ID) GetValue() uint64 { return t.value }

// Return type name
func (ID) Type() string { return "ID" }

// Convert value to uint64
func (ID) Cast(val interface{}) (TypedValue, error) {
	return NewID(val)
}

// Assign new value to ID
//
// value is first passed through CastToID
func (t *ID) Assign(val interface{}) error {
	if c, err := CastToID(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Integer is an expression type, wrapper for int64 type
type Integer struct{ value int64 }

// NewInteger creates new instance of Integer expression type
func NewInteger(val interface{}) (*Integer, error) {
	if c, err := CastToInteger(val); err != nil {
		return nil, fmt.Errorf("unable to create Integer: %w", err)
	} else {
		return &Integer{value: c}, nil
	}
}

// Return underlying value on Integer
func (t Integer) Get() interface{} { return t.value }

// Return underlying value on Integer
func (t Integer) GetValue() int64 { return t.value }

// Return type name
func (Integer) Type() string { return "Integer" }

// Convert value to int64
func (Integer) Cast(val interface{}) (TypedValue, error) {
	return NewInteger(val)
}

// Assign new value to Integer
//
// value is first passed through CastToInteger
func (t *Integer) Assign(val interface{}) error {
	if c, err := CastToInteger(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// KV is an expression type, wrapper for map[string]string type
type KV struct{ value map[string]string }

// NewKV creates new instance of KV expression type
func NewKV(val interface{}) (*KV, error) {
	if c, err := CastToKV(val); err != nil {
		return nil, fmt.Errorf("unable to create KV: %w", err)
	} else {
		return &KV{value: c}, nil
	}
}

// Return underlying value on KV
func (t KV) Get() interface{} { return t.value }

// Return underlying value on KV
func (t KV) GetValue() map[string]string { return t.value }

// Return type name
func (KV) Type() string { return "KV" }

// Convert value to map[string]string
func (KV) Cast(val interface{}) (TypedValue, error) {
	return NewKV(val)
}

// Assign new value to KV
//
// value is first passed through CastToKV
func (t *KV) Assign(val interface{}) error {
	if c, err := CastToKV(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// KVV is an expression type, wrapper for map[string][]string type
type KVV struct{ value map[string][]string }

// NewKVV creates new instance of KVV expression type
func NewKVV(val interface{}) (*KVV, error) {
	if c, err := CastToKVV(val); err != nil {
		return nil, fmt.Errorf("unable to create KVV: %w", err)
	} else {
		return &KVV{value: c}, nil
	}
}

// Return underlying value on KVV
func (t KVV) Get() interface{} { return t.value }

// Return underlying value on KVV
func (t KVV) GetValue() map[string][]string { return t.value }

// Return type name
func (KVV) Type() string { return "KVV" }

// Convert value to map[string][]string
func (KVV) Cast(val interface{}) (TypedValue, error) {
	return NewKVV(val)
}

// Assign new value to KVV
//
// value is first passed through CastToKVV
func (t *KVV) Assign(val interface{}) error {
	if c, err := CastToKVV(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Reader is an expression type, wrapper for io.Reader type
type Reader struct{ value io.Reader }

// NewReader creates new instance of Reader expression type
func NewReader(val interface{}) (*Reader, error) {
	if c, err := CastToReader(val); err != nil {
		return nil, fmt.Errorf("unable to create Reader: %w", err)
	} else {
		return &Reader{value: c}, nil
	}
}

// Return underlying value on Reader
func (t Reader) Get() interface{} { return t.value }

// Return underlying value on Reader
func (t Reader) GetValue() io.Reader { return t.value }

// Return type name
func (Reader) Type() string { return "Reader" }

// Convert value to io.Reader
func (Reader) Cast(val interface{}) (TypedValue, error) {
	return NewReader(val)
}

// Assign new value to Reader
//
// value is first passed through CastToReader
func (t *Reader) Assign(val interface{}) error {
	if c, err := CastToReader(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// String is an expression type, wrapper for string type
type String struct{ value string }

// NewString creates new instance of String expression type
func NewString(val interface{}) (*String, error) {
	if c, err := CastToString(val); err != nil {
		return nil, fmt.Errorf("unable to create String: %w", err)
	} else {
		return &String{value: c}, nil
	}
}

// Return underlying value on String
func (t String) Get() interface{} { return t.value }

// Return underlying value on String
func (t String) GetValue() string { return t.value }

// Return type name
func (String) Type() string { return "String" }

// Convert value to string
func (String) Cast(val interface{}) (TypedValue, error) {
	return NewString(val)
}

// Assign new value to String
//
// value is first passed through CastToString
func (t *String) Assign(val interface{}) error {
	if c, err := CastToString(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// UnsignedInteger is an expression type, wrapper for uint64 type
type UnsignedInteger struct{ value uint64 }

// NewUnsignedInteger creates new instance of UnsignedInteger expression type
func NewUnsignedInteger(val interface{}) (*UnsignedInteger, error) {
	if c, err := CastToUnsignedInteger(val); err != nil {
		return nil, fmt.Errorf("unable to create UnsignedInteger: %w", err)
	} else {
		return &UnsignedInteger{value: c}, nil
	}
}

// Return underlying value on UnsignedInteger
func (t UnsignedInteger) Get() interface{} { return t.value }

// Return underlying value on UnsignedInteger
func (t UnsignedInteger) GetValue() uint64 { return t.value }

// Return type name
func (UnsignedInteger) Type() string { return "UnsignedInteger" }

// Convert value to uint64
func (UnsignedInteger) Cast(val interface{}) (TypedValue, error) {
	return NewUnsignedInteger(val)
}

// Assign new value to UnsignedInteger
//
// value is first passed through CastToUnsignedInteger
func (t *UnsignedInteger) Assign(val interface{}) error {
	if c, err := CastToUnsignedInteger(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// Vars is an expression type, wrapper for RVars type
type Vars struct{ value RVars }

// NewVars creates new instance of Vars expression type
func NewVars(val interface{}) (*Vars, error) {
	if c, err := CastToVars(val); err != nil {
		return nil, fmt.Errorf("unable to create Vars: %w", err)
	} else {
		return &Vars{value: c}, nil
	}
}

// Return underlying value on Vars
func (t Vars) Get() interface{} { return t.value }

// Return underlying value on Vars
func (t Vars) GetValue() RVars { return t.value }

// Return type name
func (Vars) Type() string { return "Vars" }

// Convert value to RVars
func (Vars) Cast(val interface{}) (TypedValue, error) {
	return NewVars(val)
}

// Assign new value to Vars
//
// value is first passed through CastToVars
func (t *Vars) Assign(val interface{}) error {
	if c, err := CastToVars(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}
