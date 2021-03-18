package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	composeRecord struct {
		values map[string]string
		ts     *resource.Timestamps
		us     *resource.Userstamps
		config *resource.EnvoyConfig

		cfg *EncoderConfig

		refModule    string
		refNamespace string
	}
	composeRecordSet []*composeRecord

	composeRecordValues struct {
		rvs types.RecordValueSet
	}
)

func (nn composeRecordSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.cfg = cfg
	}
}
