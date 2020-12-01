package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestModuleExpr(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)
	)
	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	ni := uint64(0)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite string) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := dd(ctx, suite)
		req.NoError(err)

		crs := resource.ComposeRecordShaper()
		nn, err = resource.Shape(nn, crs)
		req.NoError(err)

		return req, encode(ctx, s, nn)
	}
	// Prepare
	s, err = initStore(ctx)
	err = ce(
		err,

		s.TruncateRbacRules(ctx),
		s.TruncateRoles(ctx),
		s.TruncateActionlogs(ctx),
		s.TruncateApplications(ctx),
		s.TruncateAttachments(ctx),
		s.TruncateComposeAttachments(ctx),
		s.TruncateComposeCharts(ctx),
		s.TruncateComposeNamespaces(ctx),
		s.TruncateComposeModules(ctx),
		s.TruncateComposeModuleFields(ctx),
		s.TruncateComposePages(ctx),

		storeRole(ctx, s, 1, "everyone"),
		storeRole(ctx, s, 2, "admins"),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := prepare(ctx, s, t, "mod_expr")
	req.NoError(err)

	t.Run("module field expressions", func(t *testing.T) {
		ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "crm")
		req.NoError(err)
		req.NotNil(ns)

		mod, err := fullModLoad(ctx, s, req, ns.ID, "Account")
		req.NotNil(mod)
		req.NoError(err)
		req.Len(mod.Fields, 2)

		// Check the full thing
		mfF := mod.Fields[0]
		req.Equal("a > b", mfF.Expressions.ValueExpr)
		req.Subset(mfF.Expressions.Sanitizers, []string{"trim(value)"})
		v := mfF.Expressions.Validators[0]
		req.Equal("a == \"\"", v.Test)
		req.Equal("Value should not be empty", v.Error)

		// Check the other validator form
		mfV := mod.Fields[1]
		v = mfV.Expressions.Validators[0]
		req.Equal("value == \"\"", v.Test)
		req.Equal("Value should be filled", v.Error)
	})
}
