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
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// composeBase
	//
	// This type is auto-generated.
	composeBase struct {
		immutable bool
		invoker   auth.Identifiable
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
		composeBase: &composeBase{
			immutable: false,
		},
	}
}

// ComposeOnManualImmutable creates onManual for compose resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ComposeOnManualImmutable() *composeOnManual {
	return &composeOnManual{
		composeBase: &composeBase{
			immutable: true,
		},
	}
}

// ComposeOnInterval creates onInterval for compose resource
//
// This function is auto-generated.
func ComposeOnInterval() *composeOnInterval {
	return &composeOnInterval{
		composeBase: &composeBase{
			immutable: false,
		},
	}
}

// ComposeOnIntervalImmutable creates onInterval for compose resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ComposeOnIntervalImmutable() *composeOnInterval {
	return &composeOnInterval{
		composeBase: &composeBase{
			immutable: true,
		},
	}
}

// ComposeOnTimestamp creates onTimestamp for compose resource
//
// This function is auto-generated.
func ComposeOnTimestamp() *composeOnTimestamp {
	return &composeOnTimestamp{
		composeBase: &composeBase{
			immutable: false,
		},
	}
}

// ComposeOnTimestampImmutable creates onTimestamp for compose resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ComposeOnTimestampImmutable() *composeOnTimestamp {
	return &composeOnTimestamp{
		composeBase: &composeBase{
			immutable: true,
		},
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res composeBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *composeBase) Decode(results map[string][]byte) (err error) {
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
