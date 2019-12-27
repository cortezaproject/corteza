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
	// sinkBase
	//
	// This type is auto-generated.
	sinkBase struct {
		request *types.SinkRequest
		invoker auth.Identifiable
	}

	// sinkOnRequest
	//
	// This type is auto-generated.
	sinkOnRequest struct {
		*sinkBase
	}
)

// ResourceType returns "system:sink"
//
// This function is auto-generated.
func (sinkBase) ResourceType() string {
	return "system:sink"
}

// EventType on sinkOnRequest returns "onRequest"
//
// This function is auto-generated.
func (sinkOnRequest) EventType() string {
	return "onRequest"
}

// SinkOnRequest creates onRequest for system:sink resource
//
// This function is auto-generated.
func SinkOnRequest(
	argRequest *types.SinkRequest,
) *sinkOnRequest {
	return &sinkOnRequest{
		sinkBase: &sinkBase{
			request: argRequest,
		},
	}
}

// SetRequest sets new request value
//
// This function is auto-generated.
func (res *sinkBase) SetRequest(argRequest *types.SinkRequest) {
	res.request = argRequest
}

// Request returns request
//
// This function is auto-generated.
func (res sinkBase) Request() *types.SinkRequest {
	return res.request
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *sinkBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res sinkBase) Invoker() auth.Identifiable {
	return res.invoker
}
