package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	decoder struct {
		ux *userIndex
	}

	DecodeFilter struct {

		// Compose stuff
		composeNamespace []*composeNamespaceFilter
		composeModule    []*composeModuleFilter
		composeRecord    []*composeRecordFilter
		composePage      []*composePageFilter
		composeChart     []*composeChartFilter

		// System stuff
		roles        []*roleFilter
		users        []*userFilter
		templates    []*templateFilter
		applications []*applicationFilter
		apiGwRoutes  []*apiGwRouteFilter
		reports      []*reportFilter

		settings             []*settingFilter
		rbac                 []*rbacFilter
		resourceTranslations []*resourceTranslationFilter

		// Automation stuff
		automationWorkflow []*automationWorkflowFilter
	}

	// These two aux things let us simplify the code a little bit
	auxMarshaller []envoy.Marshaller
	auxRsp        struct {
		mm  []envoy.Marshaller
		err error
	}

	// We'll use the userIndex to hold required users
	userIndex struct {
		users map[uint64]*types.User
		s     store.Users
	}
)

func Decoder() *decoder {
	return &decoder{}
}

func NewDecodeFilter() *DecodeFilter {
	return &DecodeFilter{}
}

func (df *DecodeFilter) FromResource(rr ...string) *DecodeFilter {
	df = df.systemFromResource(rr...)
	df = df.automationFromResource(rr...)
	df = df.composeFromResource(rr...)

	return df
}

func (df *DecodeFilter) FromRef(rr ...*resource.Ref) *DecodeFilter {
	df = df.systemFromRef(rr...)
	df = df.automationFromRef(rr...)
	df = df.composeFromRef(rr...)

	return df
}

func (aum auxMarshaller) MarshalEnvoy() ([]resource.Interface, error) {
	ii := make([]resource.Interface, 0, len(aum))
	for _, m := range aum {
		tmp, err := m.MarshalEnvoy()
		if err != nil {
			return nil, err
		}

		// Default locales for applicable resources
		for _, t := range tmp {
			if li, ok := t.(resource.LocaleInterface); ok {
				ll, err := li.EncodeTranslations()
				if err != nil {
					return nil, err
				}

				for _, l := range ll {
					l.MarkDefault()
					ii = append(ii, l)
				}
			}
		}

		ii = append(ii, tmp...)
	}

	return ii, nil
}

// Decode decodes all of the things in the provided store
func (d *decoder) Decode(ctx context.Context, s store.Storer, dal dalService, f *DecodeFilter) ([]resource.Interface, error) {
	mm := make(auxMarshaller, 0, 100)

	if d.ux == nil {
		d.ux = &userIndex{
			users: make(map[uint64]*types.User),
		}
	}
	d.ux.s = s

	pof := func(rr ...*auxRsp) (auxMarshaller, error) {
		mm := make(auxMarshaller, 0, 200)

		for _, r := range rr {
			if r == nil {
				continue
			}
			if r.err != nil {
				return nil, r.err
			}

			mm = append(mm, r.mm...)
		}

		return mm, nil
	}

	compose := newComposeDecoder(dal)
	system := newSystemDecoder(d.ux)
	automation := newAutomationDecoder(d.ux)

	mm, err := pof(
		compose.decodeComposeNamespace(ctx, s, f.composeNamespace),
		compose.decodeComposeModule(ctx, s, f.composeModule),
		compose.decodeComposeRecord(ctx, s, f.composeRecord),
		compose.decodeComposePage(ctx, s, f.composePage),
		compose.decodeComposeChart(ctx, s, f.composeChart),

		system.decodeRoles(ctx, s, f.roles),
		system.decodeUsers(ctx, s, f.users),
		system.decodeTemplates(ctx, s, f.templates),
		system.decodeApplications(ctx, s, f.applications),
		system.decodeAPIGWRoutes(ctx, s, f.apiGwRoutes),
		system.decodeReports(ctx, s, f.reports),
		system.decodeSettings(ctx, s, f.settings),
		system.decodeResourceTranslation(ctx, s, f.resourceTranslations),

		automation.decodeWorkflows(ctx, s, f.automationWorkflow),
	)
	if err != nil {
		return nil, err
	}

	f.allowRbacResource(compose.resourceID...)
	f.allowRbacResource(system.resourceID...)
	f.allowRbacResource(automation.resourceID...)
	rr, err := pof(
		system.decodeRbac(ctx, s, f.rbac),
	)
	if err != nil {
		return nil, err
	}

	return append(rr, mm...).MarshalEnvoy()
}

func (ux *userIndex) add(ctx context.Context, uu ...uint64) error {
	// List of filtered users
	filtered := make([]uint64, 0, len(uu))

	for _, u := range uu {
		// not defined, we don't need
		if u == 0 {
			continue
		}

		// we have, no need
		if ux.users[u] != nil {
			continue
		}

		filtered = append(filtered, u)
	}

	if len(filtered) == 0 {
		return nil
	}

	users, _, err := store.SearchUsers(ctx, ux.s, types.UserFilter{
		UserID: filtered,
	})
	if err != nil {
		return err
	}

	for _, u := range users {
		ux.users[u.ID] = u
	}

	return nil
}
