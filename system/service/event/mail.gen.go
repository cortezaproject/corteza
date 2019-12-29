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
	"github.com/cortezaproject/corteza-server/system/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// mailBase
	//
	// This type is auto-generated.
	mailBase struct {
		request *types.MailMessage
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
	argRequest *types.MailMessage,
) *mailOnManual {
	return &mailOnManual{
		mailBase: &mailBase{
			request: argRequest,
		},
	}
}

// MailOnReceive creates onReceive for system:mail resource
//
// This function is auto-generated.
func MailOnReceive(
	argRequest *types.MailMessage,
) *mailOnReceive {
	return &mailOnReceive{
		mailBase: &mailBase{
			request: argRequest,
		},
	}
}

// MailOnSend creates onSend for system:mail resource
//
// This function is auto-generated.
func MailOnSend(
	argRequest *types.MailMessage,
) *mailOnSend {
	return &mailOnSend{
		mailBase: &mailBase{
			request: argRequest,
		},
	}
}

// SetRequest sets new request value
//
// This function is auto-generated.
func (res *mailBase) SetRequest(argRequest *types.MailMessage) {
	res.request = argRequest
}

// Request returns request
//
// This function is auto-generated.
func (res mailBase) Request() *types.MailMessage {
	return res.request
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
