package service

import (
	"context"

	"github.com/crusttech/crust/messaging/internal/service"
)

// Provision orchestrates various tasks after deployment
//
func Provision(ctx context.Context) (err error) {
	if err = resetDefaultPermissionRules(ctx); err != nil {
		return
	}

	// @todo move migration here

	return
}

// Resets default permission rules for compose resources
func resetDefaultPermissionRules(ctx context.Context) error {
	var ac = service.DefaultAccessControl

	return ac.Grant(ctx, ac.DefaultRules()...)
}
