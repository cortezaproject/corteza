package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

// Addition to list of defined resouces
// we need that because messagebus resource is not compliant
func (svc accessControl) list() []map[string]string {
	return []map[string]string{
		{"type": messagebus.QueueResourceType, "any": messagebus.QueueRbacResource(0), "op": "read"},
		{"type": messagebus.QueueResourceType, "any": messagebus.QueueRbacResource(0), "op": "update"},
		{"type": messagebus.QueueResourceType, "any": messagebus.QueueRbacResource(0), "op": "delete"},
		{"type": messagebus.QueueResourceType, "any": messagebus.QueueRbacResource(0), "op": "queue.read"},
		{"type": messagebus.QueueResourceType, "any": messagebus.QueueRbacResource(0), "op": "queue.write"},
	}
}

// CanReadQueue checks if current user can read queue
//
// This function is auto-generated
func (svc accessControl) CanReadQueue(ctx context.Context, r *messagebus.QueueSettings) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateQueue checks if current user can update queue
//
// This function is auto-generated
func (svc accessControl) CanUpdateQueue(ctx context.Context, r *messagebus.QueueSettings) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteQueue checks if current user can delete queue
//
// This function is auto-generated
func (svc accessControl) CanDeleteQueue(ctx context.Context, r *messagebus.QueueSettings) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadMessageOnQueue checks if current user can read from queue
//
// This function is auto-generated
func (svc accessControl) CanReadMessageOnQueue(ctx context.Context, r *messagebus.QueueSettings) bool {
	return svc.can(ctx, "queue.read", r)
}

// CanWriteMessageOnQueue checks if current user can write to queue
//
// This function is auto-generated
func (svc accessControl) CanWriteMessageOnQueue(ctx context.Context, r *messagebus.QueueSettings) bool {
	return svc.can(ctx, "queue.write", r)
}
