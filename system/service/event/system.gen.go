package event

// This file is auto-generated.
//
// YAML event definitions:
//   system/service/event/events.yaml
//
// Regenerate with:
//   go run ./codegen/v2/events --service system
//

import (
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// systemBase
	//
	// This type is auto-generated.
	systemBase struct {
		immutable bool
		invoker   auth.Identifiable
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
		systemBase: &systemBase{
			immutable: false,
		},
	}
}

// SystemOnManualImmutable creates onManual for system resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SystemOnManualImmutable() *systemOnManual {
	return &systemOnManual{
		systemBase: &systemBase{
			immutable: true,
		},
	}
}

// SystemOnInterval creates onInterval for system resource
//
// This function is auto-generated.
func SystemOnInterval() *systemOnInterval {
	return &systemOnInterval{
		systemBase: &systemBase{
			immutable: false,
		},
	}
}

// SystemOnIntervalImmutable creates onInterval for system resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SystemOnIntervalImmutable() *systemOnInterval {
	return &systemOnInterval{
		systemBase: &systemBase{
			immutable: true,
		},
	}
}

// SystemOnTimestamp creates onTimestamp for system resource
//
// This function is auto-generated.
func SystemOnTimestamp() *systemOnTimestamp {
	return &systemOnTimestamp{
		systemBase: &systemBase{
			immutable: false,
		},
	}
}

// SystemOnTimestampImmutable creates onTimestamp for system resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SystemOnTimestampImmutable() *systemOnTimestamp {
	return &systemOnTimestamp{
		systemBase: &systemBase{
			immutable: true,
		},
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res systemBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *systemBase) Decode(results map[string][]byte) (err error) {
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
