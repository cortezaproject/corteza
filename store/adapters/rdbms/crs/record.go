package crs

const (
	sysID          = "ID"
	sysNamespaceID = "namespaceID"
	sysModuleID    = "moduleID"
	sysCreatedAt   = "createdAt"
	sysCreatedBy   = "createdBy"
	sysUpdatedAt   = "updatedAt"
	sysUpdatedBy   = "updatedBy"
	sysDeletedAt   = "deletedAt"
	sysDeletedBy   = "deletedBy"
	sysOwnedBy     = "ownedBy"
)

// is the
func isSystemField(f string) bool {
	switch f {
	// handle system fields
	case sysID,
		sysNamespaceID,
		sysModuleID,
		sysCreatedAt,
		sysCreatedBy,
		sysUpdatedAt,
		sysUpdatedBy,
		sysDeletedAt,
		sysDeletedBy,
		sysOwnedBy:
		return true
	default:
		return false
	}
}
