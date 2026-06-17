package service

import (
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/dml"
)

// Exported type aliases so rest/dml.go can refer to these without importing
// the dml sub-package directly.
type (
	DmlConnectionSvc = dml.Connection
	DmlMappingSvc    = dml.Mapping
	DmlApplierSvc    = dml.Applier
	DmlImporterSvc   = dml.Importer

	// unexported aliases keep the var declarations tidy
	dmlConnectionSvc = dml.Connection
	dmlMappingSvc    = dml.Mapping
	dmlApplierSvc    = dml.Applier
	dmlImporterSvc   = dml.Importer

	// DmlComposeDeps carries compose-service deps injected from boot_levels.go.
	DmlComposeDeps struct {
		ModuleSvc    dml.ComposeModuleSvc
		NamespaceSvc dml.ComposeNamespaceSvc
		RecordSvc    dml.ComposeRecordSvc
	}
)

// InitDmlServices initialises the core DML services that only depend on the
// DAL service and the store. Applier and Importer are wired later via
// SetDmlComposeDeps, called from app/boot_levels.go.
func InitDmlServices(s store.Storer, d dal.FullService) {
	conn := dml.NewConnection(s, d)
	mp := dml.NewMapping(conn)

	DefaultDmlConnection = conn
	DefaultDmlMapping = mp
}

// SetDmlComposeDeps wires compose-service dependencies into the DML Applier
// and Importer. Call this from boot_levels.go after compose services init.
func SetDmlComposeDeps(deps DmlComposeDeps, d dal.FullService) {
	if DefaultDmlMapping == nil {
		return
	}
	DefaultDmlApplier = dml.NewApplier(DefaultDmlMapping, deps.ModuleSvc, deps.NamespaceSvc)
	DefaultDmlImporter = dml.NewImporter(DefaultDmlMapping, DefaultDmlApplier, deps.ModuleSvc, d, deps.RecordSvc, deps.NamespaceSvc)
}
