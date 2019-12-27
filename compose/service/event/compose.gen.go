package event

// This file is auto-generated.
//
// YAML event definitions:
//   compose/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service compose
//

import (
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// composeBase
	//
	// This type is auto-generated.
	composeBase struct {
		invoker auth.Identifiable
	}

	// composeOnManual
	//
	// This type is auto-generated.
	composeOnManual struct {
		*composeBase
	}

	// composeOnInterval
	//
	// This type is auto-generated.
	composeOnInterval struct {
		*composeBase
	}

	// composeOnTimestamp
	//
	// This type is auto-generated.
	composeOnTimestamp struct {
		*composeBase
	}
)

// ResourceType returns "compose"
//
// This function is auto-generated.
func (composeBase) ResourceType() string {
	return "compose"
}

// EventType on composeOnManual returns "onManual"
//
// This function is auto-generated.
func (composeOnManual) EventType() string {
	return "onManual"
}

// EventType on composeOnInterval returns "onInterval"
//
// This function is auto-generated.
func (composeOnInterval) EventType() string {
	return "onInterval"
}

// EventType on composeOnTimestamp returns "onTimestamp"
//
// This function is auto-generated.
func (composeOnTimestamp) EventType() string {
	return "onTimestamp"
}

// ComposeOnManual creates onManual for compose resource
//
// This function is auto-generated.
func ComposeOnManual() *composeOnManual {
	return &composeOnManual{
		composeBase: &composeBase{},
	}
}

// ComposeOnInterval creates onInterval for compose resource
//
// This function is auto-generated.
func ComposeOnInterval() *composeOnInterval {
	return &composeOnInterval{
		composeBase: &composeBase{},
	}
}

// ComposeOnTimestamp creates onTimestamp for compose resource
//
// This function is auto-generated.
func ComposeOnTimestamp() *composeOnTimestamp {
	return &composeOnTimestamp{
		composeBase: &composeBase{},
	}
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *composeBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res composeBase) Invoker() auth.Identifiable {
	return res.invoker
}
