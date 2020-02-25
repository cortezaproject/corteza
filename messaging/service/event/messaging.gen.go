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
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// messagingBase
	//
	// This type is auto-generated.
	messagingBase struct {
		immutable bool
		invoker   auth.Identifiable
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
		messagingBase: &messagingBase{
			immutable: false,
		},
	}
}

// MessagingOnManualImmutable creates onManual for messaging resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessagingOnManualImmutable() *messagingOnManual {
	return &messagingOnManual{
		messagingBase: &messagingBase{
			immutable: true,
		},
	}
}

// MessagingOnInterval creates onInterval for messaging resource
//
// This function is auto-generated.
func MessagingOnInterval() *messagingOnInterval {
	return &messagingOnInterval{
		messagingBase: &messagingBase{
			immutable: false,
		},
	}
}

// MessagingOnIntervalImmutable creates onInterval for messaging resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessagingOnIntervalImmutable() *messagingOnInterval {
	return &messagingOnInterval{
		messagingBase: &messagingBase{
			immutable: true,
		},
	}
}

// MessagingOnTimestamp creates onTimestamp for messaging resource
//
// This function is auto-generated.
func MessagingOnTimestamp() *messagingOnTimestamp {
	return &messagingOnTimestamp{
		messagingBase: &messagingBase{
			immutable: false,
		},
	}
}

// MessagingOnTimestampImmutable creates onTimestamp for messaging resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessagingOnTimestampImmutable() *messagingOnTimestamp {
	return &messagingOnTimestamp{
		messagingBase: &messagingBase{
			immutable: true,
		},
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res messagingBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *messagingBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}
