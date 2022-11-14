package dal

import "strings"

// NormalizeAttrNames normalizes system identifiers to attribute identifiers
// @todo this is temporary until we re-visit field referencing
//
// The normalization step must exist since legacy versions allowed multiple
// identifier variations for the same system field, such as created_at and createdAt.
func NormalizeAttrNames(name string) string {
	switch strings.ToLower(name) {
	case
		"recordid",
		"record_id",
		"id":
		return "ID"

	case "moduleid",
		"module_id",
		"module":
		return "moduleID"

	case "owned_by",
		"ownedby":
		return "ownedBy"

	case "created_by",
		"createdby":
		return "createdBy"

	case "updated_by",
		"updatedby":
		return "updatedBy"

	case "deleted_by",
		"deletedby":
		return "deletedBy"

	case "created_at",
		"createdat":
		return "createdAt"

	case "updated_at",
		"updatedat":
		return "updatedAt"

	case "deleted_at",
		"deletedat":
		return "deletedAt"
	}

	return name

}
