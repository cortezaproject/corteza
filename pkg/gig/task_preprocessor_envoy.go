package gig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	intStore "github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

func (w *workerEnvoy) noop(_ context.Context, _ preprocessorNoop) error {
	return nil
}

func (w *workerEnvoy) experimentalExport(ctx context.Context, params preprocessorExperimentalExport) error {
	df := store.NewDecodeFilter()

	if params.id != 0 {
		df = df.ComposeNamespace(&types.NamespaceFilter{
			NamespaceID: []uint64{params.id},
		})
	} else {
		df = df.ComposeNamespace(&types.NamespaceFilter{
			Slug: params.handle,
		})
	}

	df = df.
		ComposeModule(&types.ModuleFilter{}).
		ComposePage(&types.PageFilter{}).
		ComposeChart(&types.ChartFilter{})

	res, err := w.getStoreDecoders().Decode(ctx, w.store, df)
	res = w.pruneResTr(res)
	if err != nil {
		return err
	}

	if params.inclRBAC {
		res, err = w.loadAccessControl(ctx, params, res)
		if err != nil {
			return err
		}
	}

	if params.inclTranslations {
		res, err = w.loadTranslations(ctx, params, res)
		if err != nil {
			return err
		}
	} else {
	}

	w.resources = append(w.resources, res...)
	// w.resources = append(w.resources, rules...)
	return nil
}

func (w *workerEnvoy) loadAccessControl(ctx context.Context, params preprocessorExperimentalExport, base resource.InterfaceSet) (resource.InterfaceSet, error) {
	var (
		roles     resource.InterfaceSet
		out       resource.InterfaceSet
		roleIndex map[uint64]bool
		err       error
	)

	// Prepare the roles we need to have
	if len(params.inclRoles) > 0 {
		roles, roleIndex, err = w.loadInclRoles(ctx, params.inclRoles)
		if err != nil {
			return nil, err
		}
	} else {
		roles, roleIndex, err = w.loadExclRoles(ctx, params.exclRoles)
		if err != nil {
			return nil, err
		}
	}
	out = append(base, roles...)

	// Load the RBAC rules
	rules, err := w.loadRBACRules(ctx)
	if err != nil {
		return nil, err
	}

	// Connect rules to resources; omit duplicates; preserve only desired roles
	dupIndex := make(map[string]bool)
	for _, res := range append(out, resource.InterfaceSet{}...) {
		c, ok := (res.Resource()).(rbacResource)
		if !ok {
			continue
		}

		for _, rule := range rules {
			if !roleIndex[rule.RoleID] {
				continue
			}

			if !matchResource(rule.Resource, c.RbacResource()) {
				continue
			}

			k := fmt.Sprintf("%s, %s, %d; %d", rule.Resource, rule.Operation, rule.Access, rule.RoleID)
			if _, ok := dupIndex[k]; ok {
				continue
			}

			_, ref, pp, err := resource.ParseRule(rule.Resource)
			if err != nil {
				return nil, err
			}

			dupIndex[k] = true
			out = append(out, resource.NewRbacRule(
				rule,
				strconv.FormatUint(rule.RoleID, 10),
				ref,
				rule.Resource,
				pp...,
			))
		}
	}

	return out, nil
}

func (w *workerEnvoy) loadTranslations(ctx context.Context, params preprocessorExperimentalExport, base resource.InterfaceSet) (resource.InterfaceSet, error) {
	var (
		translations systemTypes.ResourceTranslationSet
		out          resource.InterfaceSet
		err          error
	)

	// Prepare the roles we need to have
	if len(params.inclRoles) > 0 {
		translations, err = w.loadInclTranslations(ctx, params.inclLanguage)
		if err != nil {
			return nil, err
		}
	} else {
		translations, err = w.loadExclTranslations(ctx, params.exclLanguage)
		if err != nil {
			return nil, err
		}
	}

	out = base

	// Connect translations to resources
	for _, res := range append(out, resource.InterfaceSet{}...) {
		c, ok := (res.Resource()).(translatableResource)
		if !ok {
			continue
		}

		for _, trans := range translations {
			if !matchResource(trans.Resource, c.ResourceTranslation()) {
				continue
			}

			_, ref, pp, err := resource.ParseResourceTranslation(trans.Resource)
			if err != nil {
				return nil, err
			}

			out = append(out, resource.NewResourceTranslation(
				systemTypes.ResourceTranslationSet{trans},
				ref.Identifiers.First(),
				ref,
				pp...,
			))
		}
	}

	return out, nil
}

// @todo can we make this better?
func (w *workerEnvoy) loadRBACRules(ctx context.Context) ([]*rbac.Rule, error) {
	out, _, err := intStore.SearchRbacRules(ctx, w.store, rbac.RuleFilter{})
	return out, err
}

func (w *workerEnvoy) loadInclRoles(ctx context.Context, identifiers []string) (resource.InterfaceSet, map[uint64]bool, error) {
	df := store.NewDecodeFilter()

	for _, r := range identifiers {
		if id, err := strconv.ParseUint(r, 10, 64); err == nil {
			df = df.Roles(&systemTypes.RoleFilter{
				RoleID: []uint64{id},
			})
		} else {
			df = df.Roles(&systemTypes.RoleFilter{
				Handle: r,
			})
		}
	}

	res, err := w.getStoreDecoders().Decode(ctx, w.store, df)
	if err != nil {
		return nil, nil, err
	}

	index := make(map[uint64]bool)
	for _, r := range res {
		index[r.(resource.IdentifiableInterface).SysID()] = true
	}

	return res, index, nil
}

func (w *workerEnvoy) loadExclRoles(ctx context.Context, identifiers []string) (resource.InterfaceSet, map[uint64]bool, error) {
	df := store.NewDecodeFilter()

	res, err := w.getStoreDecoders().Decode(ctx, w.store, df)
	if err != nil {
		return nil, nil, err
	}

	out := make([]resource.Interface, 0, len(res))
	roleIndex := make(map[uint64]bool)
	identifierIndex := make(map[string]bool)
	push := func(res resource.Interface, r *systemTypes.Role) {
		out = append(out, res)
		roleIndex[r.ID] = true
	}

	for _, id := range identifiers {
		identifierIndex[id] = true
	}

	for _, r := range res {
		role := r.Resource().(*systemTypes.Role)
		if identifierIndex[role.Handle] || identifierIndex[strconv.FormatUint(role.ID, 10)] {
			continue
		}

		push(r, role)
	}

	return out, roleIndex, nil
}

// @todo can we make this better?
func (w *workerEnvoy) loadResourceTranslations(ctx context.Context) (systemTypes.ResourceTranslationSet, error) {
	out, _, err := intStore.SearchResourceTranslations(ctx, w.store, systemTypes.ResourceTranslationFilter{})
	return out, err
}

func (w *workerEnvoy) loadInclTranslations(ctx context.Context, lang []string) (systemTypes.ResourceTranslationSet, error) {
	translations, err := w.loadResourceTranslations(ctx)
	if err != nil {
		return nil, err
	}

	langIndex := make(map[string]bool)
	for _, l := range lang {
		langIndex[l] = true
	}

	out := make(systemTypes.ResourceTranslationSet, 0, 16)
	for _, rt := range translations {
		if langIndex[rt.Lang.String()] {
			out = append(out, rt)
		}
	}

	return out, nil
}

func (w *workerEnvoy) loadExclTranslations(ctx context.Context, lang []string) (systemTypes.ResourceTranslationSet, error) {
	translations, err := w.loadResourceTranslations(ctx)
	if err != nil {
		return nil, err
	}

	langIndex := make(map[string]bool)
	for _, l := range lang {
		langIndex[l] = true
	}

	out := make(systemTypes.ResourceTranslationSet, 0, 16)
	for _, rt := range translations {
		if !langIndex[rt.Lang.String()] {
			out = append(out, rt)
		}
	}

	return out, nil
}

func matchResource(matcher, resource string) (m bool) {
	return rbac.MatchResource(matcher, resource)
}

func (w *workerEnvoy) pruneResTr(nn []resource.Interface) (mm []resource.Interface) {
	mm = make([]resource.Interface, 0, len(nn))
	for _, n := range nn {
		if n.ResourceType() != resource.ResourceTranslationType {
			mm = append(mm, n)
			continue
		}

		r := n.(*resource.ResourceTranslation)
		if !r.IsDefault() {
			mm = append(mm, n)
			continue
		}
	}
	return mm
}
