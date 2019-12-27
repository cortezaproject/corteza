package event

// This file is auto-generated.
//
// YAML event definitions:
//   system/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service system
//

import (
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// systemBase
	//
	// This type is auto-generated.
	systemBase struct {
		invoker auth.Identifiable
	}

	// systemOnManual
	//
	// This type is auto-generated.
	systemOnManual struct {
		*systemBase
	}

	// systemOnInterval
	//
	// This type is auto-generated.
	systemOnInterval struct {
		*systemBase
	}

	// systemOnTimestamp
	//
	// This type is auto-generated.
	systemOnTimestamp struct {
		*systemBase
	}
)

// ResourceType returns "system"
//
// This function is auto-generated.
func (systemBase) ResourceType() string {
	return "system"
}

// EventType on systemOnManual returns "onManual"
//
// This function is auto-generated.
func (systemOnManual) EventType() string {
	return "onManual"
}

// EventType on systemOnInterval returns "onInterval"
//
// This function is auto-generated.
func (systemOnInterval) EventType() string {
	return "onInterval"
}

// EventType on systemOnTimestamp returns "onTimestamp"
//
// This function is auto-generated.
func (systemOnTimestamp) EventType() string {
	return "onTimestamp"
}

// SystemOnManual creates onManual for system resource
//
// This function is auto-generated.
func SystemOnManual() *systemOnManual {
	return &systemOnManual{
		systemBase: &systemBase{},
	}
}

// SystemOnInterval creates onInterval for system resource
//
// This function is auto-generated.
func SystemOnInterval() *systemOnInterval {
	return &systemOnInterval{
		systemBase: &systemBase{},
	}
}

// SystemOnTimestamp creates onTimestamp for system resource
//
// This function is auto-generated.
func SystemOnTimestamp() *systemOnTimestamp {
	return &systemOnTimestamp{
		systemBase: &systemBase{},
	}
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *systemBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res systemBase) Invoker() auth.Identifiable {
	return res.invoker
}
