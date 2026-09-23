package rest

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/stretchr/testify/require"
)

// The catalogue names every script and the bundle is the client scripts' own
// source, so a caller this instance does not know gets neither
func TestAutomation_anonymousIsRefused(t *testing.T) {
	var (
		req  = require.New(t)
		ctrl = &Automation{}
		ctx  = auth.SetIdentityToContext(context.Background(), auth.Anonymous())
	)

	_, err := ctrl.List(ctx, &request.AutomationList{})
	req.True(errors.IsUnauthorized(err), "list: %v", err)

	_, err = ctrl.Bundle(ctx, &request.AutomationBundle{Bundle: "compose", Type: "client-scripts", Ext: "js"})
	req.True(errors.IsUnauthorized(err), "bundle: %v", err)

	_, err = ctrl.TriggerScript(ctx, &request.AutomationTriggerScript{Script: "/server-scripts/x.js:default"})
	req.True(errors.IsUnauthorized(err), "trigger: %v", err)
}
