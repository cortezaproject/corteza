package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/spf13/cast"
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

func rbacQueueResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperationsQueue(r)

	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system Queue resource", o)
		}
	}

	if !strings.HasPrefix(r, messagebus.QueueResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(messagebus.QueueResourceType):], sep), sep)
		prc = []string{
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid resource path wildcard level (%d) for Queue", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

func rbacResourceOperationsQueue(r string) map[string]bool {
	return map[string]bool{
		"read":        true,
		"update":      true,
		"delete":      true,
		"queue.read":  true,
		"queue.write": true,
	}
}
