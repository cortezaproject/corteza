package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	settings map[string]interface{}
)

func (wrap settings) MarshalEnvoy() (nn []resource.Interface, err error) {
	n := resource.NewSettings(wrap)

	return []resource.Interface{n}, nil
}
