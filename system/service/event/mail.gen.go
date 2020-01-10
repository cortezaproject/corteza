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
	"encoding/json"

	"github.com/cortezaproject/corteza-server/system/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// mailBase
	//
	// This type is auto-generated.
	mailBase struct {
		message *types.MailMessage
		invoker auth.Identifiable
	}

	// mailOnManual
	//
	// This type is auto-generated.
	mailOnManual struct {
		*mailBase
	}

	// mailOnReceive
	//
	// This type is auto-generated.
	mailOnReceive struct {
		*mailBase
	}

	// mailOnSend
	//
	// This type is auto-generated.
	mailOnSend struct {
		*mailBase
	}
)

// ResourceType returns "system:mail"
//
// This function is auto-generated.
func (mailBase) ResourceType() string {
	return "system:mail"
}

// EventType on mailOnManual returns "onManual"
//
// This function is auto-generated.
func (mailOnManual) EventType() string {
	return "onManual"
}

// EventType on mailOnReceive returns "onReceive"
//
// This function is auto-generated.
func (mailOnReceive) EventType() string {
	return "onReceive"
}

// EventType on mailOnSend returns "onSend"
//
// This function is auto-generated.
func (mailOnSend) EventType() string {
	return "onSend"
}

// MailOnManual creates onManual for system:mail resource
//
// This function is auto-generated.
func MailOnManual(
	argMessage *types.MailMessage,
) *mailOnManual {
	return &mailOnManual{
		mailBase: &mailBase{
			message: argMessage,
		},
	}
}

// MailOnReceive creates onReceive for system:mail resource
//
// This function is auto-generated.
func MailOnReceive(
	argMessage *types.MailMessage,
) *mailOnReceive {
	return &mailOnReceive{
		mailBase: &mailBase{
			message: argMessage,
		},
	}
}

// MailOnSend creates onSend for system:mail resource
//
// This function is auto-generated.
func MailOnSend(
	argMessage *types.MailMessage,
) *mailOnSend {
	return &mailOnSend{
		mailBase: &mailBase{
			message: argMessage,
		},
	}
}

// SetMessage sets new message value
//
// This function is auto-generated.
func (res *mailBase) SetMessage(argMessage *types.MailMessage) {
	res.message = argMessage
}

// Message returns message
//
// This function is auto-generated.
func (res mailBase) Message() *types.MailMessage {
	return res.message
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *mailBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res mailBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res mailBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["message"], err = json.Marshal(res.message); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *mailBase) Decode(results map[string][]byte) (err error) {
	if r, ok := results["result"]; ok && len(results) == 1 {
		if err = json.Unmarshal(r, res.message); err != nil {
			return
		}
	}

	if r, ok := results["message"]; ok && len(results) == 1 {
		if err = json.Unmarshal(r, res.message); err != nil {
			return
		}
	}

	if r, ok := results["invoker"]; ok && len(results) == 1 {
		if err = json.Unmarshal(r, res.invoker); err != nil {
			return
		}
	}
	return
}
