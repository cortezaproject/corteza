package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	decoder struct{}

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
		settings     []*settingFilter
		rbac         []*rbacFilter
	}

	auxMarshaller []envoy.Marshaller
	auxRsp        struct {
		mm  []envoy.Marshaller
		err error
	}
)

func Decoder() *decoder {
	return &decoder{}
}
func NewDecodeFilter() *DecodeFilter {
	return &DecodeFilter{}
}

func (aum auxMarshaller) MarshalEnvoy() ([]resource.Interface, error) {
	ii := make([]resource.Interface, 0, len(aum))
	for _, m := range aum {
		tmp, err := m.MarshalEnvoy()
		if err != nil {
			return nil, err
		}
		ii = append(ii, tmp...)
	}

	return ii, nil
}

// Decode decodes all of the things in the provided store
func (d *decoder) Decode(ctx context.Context, s store.Storer, f *DecodeFilter) ([]resource.Interface, error) {
	mm := make(auxMarshaller, 0, 100)

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

	compose := newComposeDecoder()
	system := newSystemDecoder()

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
		system.decodeSettings(ctx, s, f.settings),
	)
	if err != nil {
		return nil, err
	}

	f.allowRbacResource(compose.resourceID...)
	f.allowRbacResource(system.resourceID...)
	rr, err := pof(
		system.decodeRbac(ctx, s, f.rbac),
	)
	if err != nil {
		return nil, err
	}

	return append(rr, mm...).MarshalEnvoy()
}
