package store

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
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
