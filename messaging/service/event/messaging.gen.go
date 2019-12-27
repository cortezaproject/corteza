package event

// This file is auto-generated.
//
// YAML event definitions:
//   messaging/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service messaging
//

import (
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// messagingBase
	//
	// This type is auto-generated.
	messagingBase struct {
		invoker auth.Identifiable
	}

	// messagingOnManual
	//
	// This type is auto-generated.
	messagingOnManual struct {
		*messagingBase
	}

	// messagingOnInterval
	//
	// This type is auto-generated.
	messagingOnInterval struct {
		*messagingBase
	}

	// messagingOnTimestamp
	//
	// This type is auto-generated.
	messagingOnTimestamp struct {
		*messagingBase
	}
)

// ResourceType returns "messaging"
//
// This function is auto-generated.
func (messagingBase) ResourceType() string {
	return "messaging"
}

// EventType on messagingOnManual returns "onManual"
//
// This function is auto-generated.
func (messagingOnManual) EventType() string {
	return "onManual"
}

// EventType on messagingOnInterval returns "onInterval"
//
// This function is auto-generated.
func (messagingOnInterval) EventType() string {
	return "onInterval"
}

// EventType on messagingOnTimestamp returns "onTimestamp"
//
// This function is auto-generated.
func (messagingOnTimestamp) EventType() string {
	return "onTimestamp"
}

// MessagingOnManual creates onManual for messaging resource
//
// This function is auto-generated.
func MessagingOnManual() *messagingOnManual {
	return &messagingOnManual{
		messagingBase: &messagingBase{},
	}
}

// MessagingOnInterval creates onInterval for messaging resource
//
// This function is auto-generated.
func MessagingOnInterval() *messagingOnInterval {
	return &messagingOnInterval{
		messagingBase: &messagingBase{},
	}
}

// MessagingOnTimestamp creates onTimestamp for messaging resource
//
// This function is auto-generated.
func MessagingOnTimestamp() *messagingOnTimestamp {
	return &messagingOnTimestamp{
		messagingBase: &messagingBase{},
	}
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *messagingBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res messagingBase) Invoker() auth.Identifiable {
	return res.invoker
}
