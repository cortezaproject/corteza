package store

import (
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
)

func newComposeNamespace(ns *types.Namespace) *composeNamespace {
	return &composeNamespace{
		ns: ns,
	}
}

func (ns *composeNamespace) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewComposeNamespace(ns.ns),
	)
}
