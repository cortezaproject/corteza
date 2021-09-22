package store

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

func newResourceTranslation(l types.ResourceTranslationSet) (*resourceTranslation, error) {
	res := l[0].Resource
	_, ref, pp, err := resource.ParseResourceTranslation(res)

	return &resourceTranslation{
		locales: l,

		refResourceTranslation: res,
		refLocaleRes:           ref,

		refPathRes: pp,
	}, err
}

func (lr *resourceTranslation) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewResourceTranslation(lr.locales, lr.refResourceTranslation, lr.refLocaleRes, lr.refPathRes...),
	)
}
