package store

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

func newTemplate(t *types.Template) *template {
	return &template{
		t: t,
	}
}

// MarshalEnvoy converts the template struct to a resource
func (u *template) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewTemplate(u.t),
	)
}
