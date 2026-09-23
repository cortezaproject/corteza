package rest

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/stretchr/testify/require"
)

// Listing the catalogue asks for its own permission; the bundle and the trigger
// ask to be someone this instance knows
func TestAutomation_anonymousIsRefused(t *testing.T) {
	var (
		req  = require.New(t)
		ctrl = &Automation{ac: refusingAccessControl{}}
		ctx  = auth.SetIdentityToContext(context.Background(), auth.Anonymous())
	)

	_, err := ctrl.List(ctx, &request.AutomationList{})
	req.True(errors.IsUnauthorized(err), "list: %v", err)

	_, err = ctrl.Bundle(ctx, &request.AutomationBundle{Bundle: "admin", Type: "client-scripts", Ext: "js"})
	req.True(errors.IsUnauthorized(err), "bundle: %v", err)

	_, err = ctrl.TriggerScript(ctx, &request.AutomationTriggerScript{Script: "/server-scripts/x.js:default"})
	req.True(errors.IsUnauthorized(err), "trigger: %v", err)
}

type refusingAccessControl struct{}

func (refusingAccessControl) CanSearchCorredorScripts(context.Context) bool { return false }
