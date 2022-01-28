package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/expr_types.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"sync"
)

var _ = context.Background
var _ = fmt.Errorf

// EmailMessage is an expression type, wrapper for *emailMessage type
type EmailMessage struct {
	value *emailMessage
	mux   sync.RWMutex
}

// NewEmailMessage creates new instance of EmailMessage expression type
func NewEmailMessage(val interface{}) (*EmailMessage, error) {
	if c, err := CastToEmailMessage(val); err != nil {
		return nil, fmt.Errorf("unable to create EmailMessage: %w", err)
	} else {
		return &EmailMessage{value: c}, nil
	}
}

// Get return underlying value on EmailMessage
func (t *EmailMessage) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on EmailMessage
func (t *EmailMessage) GetValue() *emailMessage {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (EmailMessage) Type() string { return "EmailMessage" }

// Cast converts value to *emailMessage
func (EmailMessage) Cast(val interface{}) (TypedValue, error) {
	return NewEmailMessage(val)
}

// Assign new value to EmailMessage
//
// value is first passed through CastToEmailMessage
func (t *EmailMessage) Assign(val interface{}) error {
	if c, err := CastToEmailMessage(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// HttpRequest is an expression type, wrapper for *types.HttpRequest type
type HttpRequest struct {
	value *types.HttpRequest
	mux   sync.RWMutex
}

// NewHttpRequest creates new instance of HttpRequest expression type
func NewHttpRequest(val interface{}) (*HttpRequest, error) {
	if c, err := CastToHttpRequest(val); err != nil {
		return nil, fmt.Errorf("unable to create HttpRequest: %w", err)
	} else {
		return &HttpRequest{value: c}, nil
	}
}

// Get return underlying value on HttpRequest
func (t *HttpRequest) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on HttpRequest
func (t *HttpRequest) GetValue() *types.HttpRequest {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (HttpRequest) Type() string { return "HttpRequest" }

// Cast converts value to *types.HttpRequest
func (HttpRequest) Cast(val interface{}) (TypedValue, error) {
	return NewHttpRequest(val)
}

// Assign new value to HttpRequest
//
// value is first passed through CastToHttpRequest
func (t *HttpRequest) Assign(val interface{}) error {
	if c, err := CastToHttpRequest(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *HttpRequest) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToHttpRequest(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access HttpRequest's underlying value (*types.HttpRequest)
// and it's fields
//
func (t *HttpRequest) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return httpRequestGValSelector(t.value, k)
}

// Select is field accessor for *types.HttpRequest
//
// Similar to SelectGVal but returns typed values
func (t *HttpRequest) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return httpRequestTypedValueSelector(t.value, k)
}

func (t *HttpRequest) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "Method":
		return true
	case "URL":
		return true
	case "Header":
		return true
	case "Body":
		return true
	case "Form":
		return true
	case "PostForm":
		return true
	}
	return false
}

// httpRequestGValSelector is field accessor for *types.HttpRequest
func httpRequestGValSelector(res *types.HttpRequest, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Method":
		return res.Method, nil
	case "URL":
		return res.URL, nil
	case "Header":
		return res.Header, nil
	case "Body":
		return res.Body, nil
	case "Form":
		return res.Form, nil
	case "PostForm":
		return res.PostForm, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// httpRequestTypedValueSelector is field accessor for *types.HttpRequest
func httpRequestTypedValueSelector(res *types.HttpRequest, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Method":
		return NewString(res.Method)
	case "URL":
		return NewUrl(res.URL)
	case "Header":
		return NewKVV(res.Header)
	case "Body":
		return NewHttpRequestBody(res.Body)
	case "Form":
		return NewKVV(res.Form)
	case "PostForm":
		return NewKVV(res.PostForm)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToHttpRequest is field value setter for *types.HttpRequest
func assignToHttpRequest(res *types.HttpRequest, k string, val interface{}) error {
	switch k {
	case "Method":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Method = aux
		return nil
	case "URL":
		aux, err := CastToUrl(val)
		if err != nil {
			return err
		}

		res.URL = aux
		return nil
	case "Header":
		aux, err := CastToKVV(val)
		if err != nil {
			return err
		}

		res.Header = aux
		return nil
	case "Body":
		aux, err := CastToHttpRequestBody(val)
		if err != nil {
			return err
		}

		res.Body = aux
		return nil
	case "Form":
		aux, err := CastToKVV(val)
		if err != nil {
			return err
		}

		res.Form = aux
		return nil
	case "PostForm":
		aux, err := CastToKVV(val)
		if err != nil {
			return err
		}

		res.PostForm = aux
		return nil
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// HttpRequestBody is an expression type, wrapper for *types.HttpRequestBody type
type HttpRequestBody struct {
	value *types.HttpRequestBody
	mux   sync.RWMutex
}

// NewHttpRequestBody creates new instance of HttpRequestBody expression type
func NewHttpRequestBody(val interface{}) (*HttpRequestBody, error) {
	if c, err := CastToHttpRequestBody(val); err != nil {
		return nil, fmt.Errorf("unable to create HttpRequestBody: %w", err)
	} else {
		return &HttpRequestBody{value: c}, nil
	}
}

// Get return underlying value on HttpRequestBody
func (t *HttpRequestBody) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on HttpRequestBody
func (t *HttpRequestBody) GetValue() *types.HttpRequestBody {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (HttpRequestBody) Type() string { return "HttpRequestBody" }

// Cast converts value to *types.HttpRequestBody
func (HttpRequestBody) Cast(val interface{}) (TypedValue, error) {
	return NewHttpRequestBody(val)
}

// Assign new value to HttpRequestBody
//
// value is first passed through CastToHttpRequestBody
func (t *HttpRequestBody) Assign(val interface{}) error {
	if c, err := CastToHttpRequestBody(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *HttpRequestBody) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToHttpRequestBody(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access HttpRequestBody's underlying value (*types.HttpRequestBody)
// and it's fields
//
func (t *HttpRequestBody) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return httpRequestBodyGValSelector(t.value, k)
}

// Select is field accessor for *types.HttpRequestBody
//
// Similar to SelectGVal but returns typed values
func (t *HttpRequestBody) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return httpRequestBodyTypedValueSelector(t.value, k)
}

func (t *HttpRequestBody) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "Body":
		return true
	case "Buffer":
		return true
	}
	return false
}

// httpRequestBodyGValSelector is field accessor for *types.HttpRequestBody
func httpRequestBodyGValSelector(res *types.HttpRequestBody, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Body":
		return res.Body, nil
	case "Buffer":
		return res.Buffer, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// httpRequestBodyTypedValueSelector is field accessor for *types.HttpRequestBody
func httpRequestBodyTypedValueSelector(res *types.HttpRequestBody, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Body":
		return NewReader(res.Body)
	case "Buffer":
		return NewBytes(res.Buffer)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToHttpRequestBody is field value setter for *types.HttpRequestBody
func assignToHttpRequestBody(res *types.HttpRequestBody, k string, val interface{}) error {
	switch k {
	case "Body":
		aux, err := CastToReader(val)
		if err != nil {
			return err
		}

		res.Body = aux
		return nil
	case "Buffer":
		aux, err := CastToBytes(val)
		if err != nil {
			return err
		}

		res.Buffer = aux
		return nil
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// Url is an expression type, wrapper for *types.Url type
type Url struct {
	value *types.Url
	mux   sync.RWMutex
}

// NewUrl creates new instance of Url expression type
func NewUrl(val interface{}) (*Url, error) {
	if c, err := CastToUrl(val); err != nil {
		return nil, fmt.Errorf("unable to create Url: %w", err)
	} else {
		return &Url{value: c}, nil
	}
}

// Get return underlying value on Url
func (t *Url) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Url
func (t *Url) GetValue() *types.Url {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Url) Type() string { return "Url" }

// Cast converts value to *types.Url
func (Url) Cast(val interface{}) (TypedValue, error) {
	return NewUrl(val)
}

// Assign new value to Url
//
// value is first passed through CastToUrl
func (t *Url) Assign(val interface{}) error {
	if c, err := CastToUrl(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *Url) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToUrl(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Url's underlying value (*types.Url)
// and it's fields
//
func (t *Url) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return urlGValSelector(t.value, k)
}

// Select is field accessor for *types.Url
//
// Similar to SelectGVal but returns typed values
func (t *Url) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return urlTypedValueSelector(t.value, k)
}

func (t *Url) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "Scheme":
		return true
	case "Opaque":
		return true
	case "Host":
		return true
	case "Path":
		return true
	case "RawPath":
		return true
	case "ForceQuery":
		return true
	case "RawQuery":
		return true
	case "Fragment":
		return true
	case "RawFragment":
		return true
	}
	return false
}

// urlGValSelector is field accessor for *types.Url
func urlGValSelector(res *types.Url, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Scheme":
		return res.Scheme, nil
	case "Opaque":
		return res.Opaque, nil
	case "Host":
		return res.Host, nil
	case "Path":
		return res.Path, nil
	case "RawPath":
		return res.RawPath, nil
	case "ForceQuery":
		return res.ForceQuery, nil
	case "RawQuery":
		return res.RawQuery, nil
	case "Fragment":
		return res.Fragment, nil
	case "RawFragment":
		return res.RawFragment, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// urlTypedValueSelector is field accessor for *types.Url
func urlTypedValueSelector(res *types.Url, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Scheme":
		return NewString(res.Scheme)
	case "Opaque":
		return NewString(res.Opaque)
	case "Host":
		return NewString(res.Host)
	case "Path":
		return NewString(res.Path)
	case "RawPath":
		return NewString(res.RawPath)
	case "ForceQuery":
		return NewBoolean(res.ForceQuery)
	case "RawQuery":
		return NewString(res.RawQuery)
	case "Fragment":
		return NewString(res.Fragment)
	case "RawFragment":
		return NewString(res.RawFragment)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToUrl is field value setter for *types.Url
func assignToUrl(res *types.Url, k string, val interface{}) error {
	switch k {
	case "Scheme":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Scheme = aux
		return nil
	case "Opaque":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Opaque = aux
		return nil
	case "Host":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Host = aux
		return nil
	case "Path":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Path = aux
		return nil
	case "RawPath":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.RawPath = aux
		return nil
	case "ForceQuery":
		aux, err := CastToBoolean(val)
		if err != nil {
			return err
		}

		res.ForceQuery = aux
		return nil
	case "RawQuery":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.RawQuery = aux
		return nil
	case "Fragment":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Fragment = aux
		return nil
	case "RawFragment":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.RawFragment = aux
		return nil
	}

	return fmt.Errorf("unknown field '%s'", k)
}
