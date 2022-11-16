package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza/server/pkg/envoy/store"
	"github.com/cortezaproject/corteza/server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestYamlStore_rbac(t *testing.T) {
	var (
		ctx       = context.Background()
		namespace = "base"
		s         = initServices(ctx, t)
		f         = "rbac_rules.yaml"
		req       = require.New(t)
	)

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	truncateStore(ctx, s, t)

	// util to check rule's operation
	checkPfx := func(req *require.Assertions, rr rbac.RuleSet, pfx string) {
		for _, r := range rr {
			req.Contains(r.Operation, pfx)
		}
	}
	// util to check the current store state if it's as expected
	checkStore := func() {
		// Preload all the resource
		ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
		req.NoError(err)

		mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "mod1")
		req.NoError(err)

		fld, err := store.LookupComposeModuleFieldByModuleIDName(ctx, s, mod.ID, "f1")
		req.NoError(err)

		chr, err := store.LookupComposeChartByNamespaceIDHandle(ctx, s, ns.ID, "chr1")
		req.NoError(err)

		// Preload RBAC rules
		rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		req.NoError(err)

		// Check
		// - component
		composeRR := rbacRuleFilterByResource(rr, types.ComponentRbacResource())
		req.Len(composeRR, 2)
		req.Len(composeRR.FilterAccess(rbac.Allow), 1)
		req.Len(composeRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, composeRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, composeRR.FilterAccess(rbac.Deny), "deny")

		// - namespace
		nsRR := rbacRuleFilterByResource(rr, types.NamespaceRbacResource(0))
		req.Len(nsRR, 2)
		req.Len(nsRR.FilterAccess(rbac.Allow), 1)
		req.Len(nsRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, nsRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, nsRR.FilterAccess(rbac.Deny), "deny")

		nsRR = rbacRuleFilterByResource(rr, ns.RbacResource())
		req.Len(nsRR, 2)
		req.Len(nsRR.FilterAccess(rbac.Allow), 1)
		req.Len(nsRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, nsRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, nsRR.FilterAccess(rbac.Deny), "deny")

		// - module
		modRR := rbacRuleFilterByResource(rr, types.ModuleRbacResource(0, 0))
		req.Len(modRR, 2)
		req.Len(modRR.FilterAccess(rbac.Allow), 1)
		req.Len(modRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, modRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, modRR.FilterAccess(rbac.Deny), "deny")

		modRR = rbacRuleFilterByResource(rr, types.ModuleRbacResource(ns.ID, 0))
		req.Len(modRR, 2)
		req.Len(modRR.FilterAccess(rbac.Allow), 1)
		req.Len(modRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, modRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, modRR.FilterAccess(rbac.Deny), "deny")

		modRR = rbacRuleFilterByResource(rr, types.ModuleRbacResource(ns.ID, mod.ID))
		req.Len(modRR, 2)
		req.Len(modRR.FilterAccess(rbac.Allow), 1)
		req.Len(modRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, modRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, modRR.FilterAccess(rbac.Deny), "deny")

		// - module field
		mfRR := rbacRuleFilterByResource(rr, types.ModuleFieldRbacResource(0, 0, 0))
		req.Len(mfRR, 2)
		req.Len(mfRR.FilterAccess(rbac.Allow), 1)
		req.Len(mfRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, mfRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, mfRR.FilterAccess(rbac.Deny), "deny")

		mfRR = rbacRuleFilterByResource(rr, types.ModuleFieldRbacResource(ns.ID, 0, 0))
		req.Len(mfRR, 2)
		req.Len(mfRR.FilterAccess(rbac.Allow), 1)
		req.Len(mfRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, mfRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, mfRR.FilterAccess(rbac.Deny), "deny")

		mfRR = rbacRuleFilterByResource(rr, types.ModuleFieldRbacResource(ns.ID, mod.ID, 0))
		req.Len(mfRR, 2)
		req.Len(mfRR.FilterAccess(rbac.Allow), 1)
		req.Len(mfRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, mfRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, mfRR.FilterAccess(rbac.Deny), "deny")

		mfRR = rbacRuleFilterByResource(rr, types.ModuleFieldRbacResource(ns.ID, mod.ID, fld.ID))
		req.Len(mfRR, 2)
		req.Len(mfRR.FilterAccess(rbac.Allow), 1)
		req.Len(mfRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, mfRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, mfRR.FilterAccess(rbac.Deny), "deny")

		// - chart
		chrRR := rbacRuleFilterByResource(rr, types.ChartRbacResource(0, 0))
		req.Len(chrRR, 2)
		req.Len(chrRR.FilterAccess(rbac.Allow), 1)
		req.Len(chrRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, chrRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, chrRR.FilterAccess(rbac.Deny), "deny")

		chrRR = rbacRuleFilterByResource(rr, types.ChartRbacResource(ns.ID, 0))
		req.Len(chrRR, 2)
		req.Len(chrRR.FilterAccess(rbac.Allow), 1)
		req.Len(chrRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, chrRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, chrRR.FilterAccess(rbac.Deny), "deny")

		chrRR = rbacRuleFilterByResource(rr, types.ChartRbacResource(ns.ID, chr.ID))
		req.Len(chrRR, 2)
		req.Len(chrRR.FilterAccess(rbac.Allow), 1)
		req.Len(chrRR.FilterAccess(rbac.Deny), 1)
		checkPfx(req, chrRR.FilterAccess(rbac.Allow), "allow")
		checkPfx(req, chrRR.FilterAccess(rbac.Deny), "deny")
	}

	// YAML -> Store
	nn, err := decodeYaml(ctx, namespace, f)
	req.NoError(err)
	err = encode(ctx, s, nn)
	req.NoError(err)

	checkStore()

	// Store -> YAML -> Store
	df := su.NewDecodeFilter().
		ComposeNamespace(&types.NamespaceFilter{
			Slug: "ns1",
		}).
		ComposeModule(&types.ModuleFilter{}).
		ComposeChart(&types.ChartFilter{}).
		Roles(&systemTypes.RoleFilter{Handle: "r1"}).
		Rbac(&rbac.RuleFilter{})

	// - decode store
	sd := su.Decoder()
	nn, err = sd.Decode(ctx, s, df)
	req.NoError(err)
	//
	// - into yaml
	ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{})
	bld := envoy.NewBuilder(ye)
	g, err := bld.Build(ctx, nn...)
	req.NoError(err)
	err = envoy.Encode(ctx, g, ye)
	req.NoError(err)
	ss := ye.Stream()
	//
	// - cleanup the store
	truncateStore(ctx, s, t)
	//
	// - back into store
	se := su.NewStoreEncoder(s, &su.EncoderConfig{})
	yd := yaml.Decoder()
	nn = make([]resource.Interface, 0, len(nn))
	for _, s := range ss {
		mm, err := yd.Decode(ctx, s.Source, nil)
		req.NoError(err)
		nn = append(nn, mm...)
	}
	bld = envoy.NewBuilder(se)
	g, err = bld.Build(ctx, nn...)
	req.NoError(err)

	err = envoy.Encode(ctx, g, se)
	req.NoError(err)
	//
	// - check store state
	checkStore()

	truncateStore(ctx, s, t)
}

func rbacRuleFilterByResource(rr rbac.RuleSet, res string) rbac.RuleSet {
	out := make(rbac.RuleSet, 0, len(rr)/2)

	for _, r := range rr {
		if r.Resource == res {
			out = append(out, r)
		}
	}

	return out
}
