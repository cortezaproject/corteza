package service

import (
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/dml"
)

// DmlComposeDeps carries compose-service deps injected from boot_levels.go.
type DmlComposeDeps struct {
	ModuleSvc    dml.ComposeModuleSvc
	NamespaceSvc dml.ComposeNamespaceSvc
	RecordSvc    dml.ComposeRecordSvc
}

// InitDmlServices initialises the core DML services that only depend on the
// DAL service and the store. Applier and Importer are wired later via
// SetDmlComposeDeps, called from app/boot_levels.go.
func InitDmlServices(s store.Storer, d dal.FullService) {
	DefaultDmlConnection = NewDmlConnectionSvc(s, d, DefaultAccessControl)
	DefaultDmlMapping = NewDmlMappingSvc(s, DefaultDmlConnection, DefaultAccessControl)
}

// SetDmlComposeDeps wires compose-service dependencies into the DML Applier
// and Importer. Call this from boot_levels.go after compose services init.
func SetDmlComposeDeps(deps DmlComposeDeps, d dal.FullService) {
	if DefaultDmlMapping == nil {
		return
	}
	DefaultDmlApplier = dml.NewApplier(DefaultDmlMapping, deps.ModuleSvc, deps.NamespaceSvc)
	DefaultDmlImporter = dml.NewImporter(
		DefaultDmlMapping,
		DefaultDmlApplier,
		deps.ModuleSvc,
		d,
		deps.RecordSvc,
		deps.NamespaceSvc,
		DefaultStore,
		DefaultAccessControl,
	)
}
