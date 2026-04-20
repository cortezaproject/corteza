package store

import (
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
)

func newComposeRecordFromAux(rec *composeRecordAux) *composeRecord {
	return &composeRecord{
		rec:    rec,
		relMod: rec.relMod,
	}
}

func (rec *composeRecord) MarshalEnvoy() ([]resource.Interface, error) {
	rr := resource.NewComposeRecordSet(rec.rec.walker, rec.rec.refNs, rec.rec.refMod)
	rr.SetUserFlakes(rec.rec.relUsers)
	rr.RelMod = rec.relMod

	return envoy.CollectNodes(
		rr,
	)
}
