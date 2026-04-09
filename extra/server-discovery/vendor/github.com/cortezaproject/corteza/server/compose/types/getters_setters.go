package types

func (r *Record) GetValues(name string) ([]any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return []any{r.CreatedAt}, nil
	case "createdBy", "CreatedBy", "created_by":
		return []any{r.CreatedBy}, nil
	case "deletedAt", "DeletedAt":
		return []any{r.DeletedAt}, nil
	case "deletedBy", "DeletedBy", "deleted_by":
		return []any{r.DeletedBy}, nil
	case "id", "ID":
		return []any{r.ID}, nil
	case "meta", "Meta":
		return []any{r.Meta}, nil
	case "moduleID", "ModuleID":
		return []any{r.ModuleID}, nil
	case "namespaceID", "NamespaceID":
		return []any{r.NamespaceID}, nil
	case "ownedBy", "OwnedBy", "owned_by":
		return []any{r.OwnedBy}, nil
	case "revision", "Revision":
		return []any{r.Revision}, nil
	case "updatedAt", "UpdatedAt":
		return []any{r.UpdatedAt}, nil
	case "updatedBy", "UpdatedBy", "updated_by":
		return []any{r.UpdatedBy}, nil

	default:
		return r.getValues(name)

	}
	return nil, nil
}
