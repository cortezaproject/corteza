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
	. "github.com/cortezaproject/corteza/server/pkg/expr"
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
