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
	"sync"
	"time"
)

var _ = context.Background
var _ = fmt.Errorf

// Any is an expression type, wrapper for interface{} type
type Any struct {
	value interface{}
	mux   sync.RWMutex
}

// NewAny creates new instance of Any expression type
func NewAny(val interface{}) (*Any, error) {
	if c, err := CastToAny(val); err != nil {
		return nil, fmt.Errorf("unable to create Any: %w", err)
	} else {
		return &Any{value: c}, nil
	}
}

// Get return underlying value on Any
func (t *Any) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Any
func (t *Any) GetValue() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Any) Type() string { return "Any" }

// Cast converts value to interface{}
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
type Array struct {
	value []TypedValue
	mux   sync.RWMutex
}

// NewArray creates new instance of Array expression type
func NewArray(val interface{}) (*Array, error) {
	if c, err := CastToArray(val); err != nil {
		return nil, fmt.Errorf("unable to create Array: %w", err)
	} else {
		return &Array{value: c}, nil
	}
}

// Get return underlying value on Array
func (t *Array) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Array
func (t *Array) GetValue() []TypedValue {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Array) Type() string { return "Array" }

// Cast converts value to []TypedValue
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
type Boolean struct {
	value bool
	mux   sync.RWMutex
}

// NewBoolean creates new instance of Boolean expression type
func NewBoolean(val interface{}) (*Boolean, error) {
	if c, err := CastToBoolean(val); err != nil {
		return nil, fmt.Errorf("unable to create Boolean: %w", err)
	} else {
		return &Boolean{value: c}, nil
	}
}

// Get return underlying value on Boolean
func (t *Boolean) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Boolean
func (t *Boolean) GetValue() bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Boolean) Type() string { return "Boolean" }

// Cast converts value to bool
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

// Compare the two Boolean values
func (t Boolean) Compare(to TypedValue) (int, error) {
	return compareToBoolean(t, to)
}

// Bytes is an expression type, wrapper for []byte type
type Bytes struct {
	value []byte
	mux   sync.RWMutex
}

// NewBytes creates new instance of Bytes expression type
func NewBytes(val interface{}) (*Bytes, error) {
	if c, err := CastToBytes(val); err != nil {
		return nil, fmt.Errorf("unable to create Bytes: %w", err)
	} else {
		return &Bytes{value: c}, nil
	}
}

// Get return underlying value on Bytes
func (t *Bytes) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Bytes
func (t *Bytes) GetValue() []byte {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Bytes) Type() string { return "Bytes" }

// Cast converts value to []byte
func (Bytes) Cast(val interface{}) (TypedValue, error) {
	return NewBytes(val)
}

// Assign new value to Bytes
//
// value is first passed through CastToBytes
func (t *Bytes) Assign(val interface{}) error {
	if c, err := CastToBytes(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// DateTime is an expression type, wrapper for *time.Time type
type DateTime struct {
	value *time.Time
	mux   sync.RWMutex
}

// NewDateTime creates new instance of DateTime expression type
func NewDateTime(val interface{}) (*DateTime, error) {
	if c, err := CastToDateTime(val); err != nil {
		return nil, fmt.Errorf("unable to create DateTime: %w", err)
	} else {
		return &DateTime{value: c}, nil
	}
}

// Get return underlying value on DateTime
func (t *DateTime) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on DateTime
func (t *DateTime) GetValue() *time.Time {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (DateTime) Type() string { return "DateTime" }

// Cast converts value to *time.Time
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

// Compare the two DateTime values
func (t DateTime) Compare(to TypedValue) (int, error) {
	return compareToDateTime(t, to)
}

// Duration is an expression type, wrapper for time.Duration type
type Duration struct {
	value time.Duration
	mux   sync.RWMutex
}

// NewDuration creates new instance of Duration expression type
func NewDuration(val interface{}) (*Duration, error) {
	if c, err := CastToDuration(val); err != nil {
		return nil, fmt.Errorf("unable to create Duration: %w", err)
	} else {
		return &Duration{value: c}, nil
	}
}

// Get return underlying value on Duration
func (t *Duration) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Duration
func (t *Duration) GetValue() time.Duration {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Duration) Type() string { return "Duration" }

// Cast converts value to time.Duration
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

// Compare the two Duration values
func (t Duration) Compare(to TypedValue) (int, error) {
	return compareToDuration(t, to)
}

// Float is an expression type, wrapper for float64 type
type Float struct {
	value float64
	mux   sync.RWMutex
}

// NewFloat creates new instance of Float expression type
func NewFloat(val interface{}) (*Float, error) {
	if c, err := CastToFloat(val); err != nil {
		return nil, fmt.Errorf("unable to create Float: %w", err)
	} else {
		return &Float{value: c}, nil
	}
}

// Get return underlying value on Float
func (t *Float) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Float
func (t *Float) GetValue() float64 {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Float) Type() string { return "Float" }

// Cast converts value to float64
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

// Compare the two Float values
func (t Float) Compare(to TypedValue) (int, error) {
	c, err := NewFloat(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}

// Handle is an expression type, wrapper for string type
type Handle struct {
	value string
	mux   sync.RWMutex
}

// NewHandle creates new instance of Handle expression type
func NewHandle(val interface{}) (*Handle, error) {
	if c, err := CastToHandle(val); err != nil {
		return nil, fmt.Errorf("unable to create Handle: %w", err)
	} else {
		return &Handle{value: c}, nil
	}
}

// Get return underlying value on Handle
func (t *Handle) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Handle
func (t *Handle) GetValue() string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Handle) Type() string { return "Handle" }

// Cast converts value to string
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

// Compare the two Handle values
func (t Handle) Compare(to TypedValue) (int, error) {
	c, err := NewHandle(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}

// ID is an expression type, wrapper for uint64 type
type ID struct {
	value uint64
	mux   sync.RWMutex
}

// NewID creates new instance of ID expression type
func NewID(val interface{}) (*ID, error) {
	if c, err := CastToID(val); err != nil {
		return nil, fmt.Errorf("unable to create ID: %w", err)
	} else {
		return &ID{value: c}, nil
	}
}

// Get return underlying value on ID
func (t *ID) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on ID
func (t *ID) GetValue() uint64 {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (ID) Type() string { return "ID" }

// Cast converts value to uint64
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

// Compare the two ID values
func (t ID) Compare(to TypedValue) (int, error) {
	c, err := NewID(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}

// Integer is an expression type, wrapper for int64 type
type Integer struct {
	value int64
	mux   sync.RWMutex
}

// NewInteger creates new instance of Integer expression type
func NewInteger(val interface{}) (*Integer, error) {
	if c, err := CastToInteger(val); err != nil {
		return nil, fmt.Errorf("unable to create Integer: %w", err)
	} else {
		return &Integer{value: c}, nil
	}
}

// Get return underlying value on Integer
func (t *Integer) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Integer
func (t *Integer) GetValue() int64 {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Integer) Type() string { return "Integer" }

// Cast converts value to int64
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

// Compare the two Integer values
func (t Integer) Compare(to TypedValue) (int, error) {
	c, err := NewInteger(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}

// KV is an expression type, wrapper for map[string]string type
type KV struct {
	value map[string]string
	mux   sync.RWMutex
}

// NewKV creates new instance of KV expression type
func NewKV(val interface{}) (*KV, error) {
	if c, err := CastToKV(val); err != nil {
		return nil, fmt.Errorf("unable to create KV: %w", err)
	} else {
		return &KV{value: c}, nil
	}
}

// Get return underlying value on KV
func (t *KV) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on KV
func (t *KV) GetValue() map[string]string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (KV) Type() string { return "KV" }

// Cast converts value to map[string]string
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
type KVV struct {
	value map[string][]string
	mux   sync.RWMutex
}

// NewKVV creates new instance of KVV expression type
func NewKVV(val interface{}) (*KVV, error) {
	if c, err := CastToKVV(val); err != nil {
		return nil, fmt.Errorf("unable to create KVV: %w", err)
	} else {
		return &KVV{value: c}, nil
	}
}

// Get return underlying value on KVV
func (t *KVV) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on KVV
func (t *KVV) GetValue() map[string][]string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (KVV) Type() string { return "KVV" }

// Cast converts value to map[string][]string
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
type Reader struct {
	value io.Reader
	mux   sync.RWMutex
}

// NewReader creates new instance of Reader expression type
func NewReader(val interface{}) (*Reader, error) {
	if c, err := CastToReader(val); err != nil {
		return nil, fmt.Errorf("unable to create Reader: %w", err)
	} else {
		return &Reader{value: c}, nil
	}
}

// Get return underlying value on Reader
func (t *Reader) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Reader
func (t *Reader) GetValue() io.Reader {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Reader) Type() string { return "Reader" }

// Cast converts value to io.Reader
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
type String struct {
	value string
	mux   sync.RWMutex
}

// NewString creates new instance of String expression type
func NewString(val interface{}) (*String, error) {
	if c, err := CastToString(val); err != nil {
		return nil, fmt.Errorf("unable to create String: %w", err)
	} else {
		return &String{value: c}, nil
	}
}

// Get return underlying value on String
func (t *String) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on String
func (t *String) GetValue() string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (String) Type() string { return "String" }

// Cast converts value to string
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

// Compare the two String values
func (t String) Compare(to TypedValue) (int, error) {
	c, err := NewString(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}

// UnsignedInteger is an expression type, wrapper for uint64 type
type UnsignedInteger struct {
	value uint64
	mux   sync.RWMutex
}

// NewUnsignedInteger creates new instance of UnsignedInteger expression type
func NewUnsignedInteger(val interface{}) (*UnsignedInteger, error) {
	if c, err := CastToUnsignedInteger(val); err != nil {
		return nil, fmt.Errorf("unable to create UnsignedInteger: %w", err)
	} else {
		return &UnsignedInteger{value: c}, nil
	}
}

// Get return underlying value on UnsignedInteger
func (t *UnsignedInteger) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on UnsignedInteger
func (t *UnsignedInteger) GetValue() uint64 {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (UnsignedInteger) Type() string { return "UnsignedInteger" }

// Cast converts value to uint64
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

// Compare the two UnsignedInteger values
func (t UnsignedInteger) Compare(to TypedValue) (int, error) {
	c, err := NewUnsignedInteger(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}

// Vars is an expression type, wrapper for map[string]TypedValue type
type Vars struct {
	value map[string]TypedValue
	mux   sync.RWMutex
}

// NewVars creates new instance of Vars expression type
func NewVars(val interface{}) (*Vars, error) {
	if c, err := CastToVars(val); err != nil {
		return nil, fmt.Errorf("unable to create Vars: %w", err)
	} else {
		return &Vars{value: c}, nil
	}
}

// Get return underlying value on Vars
func (t *Vars) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Vars
func (t *Vars) GetValue() map[string]TypedValue {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Vars) Type() string { return "Vars" }

// Cast converts value to map[string]TypedValue
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
