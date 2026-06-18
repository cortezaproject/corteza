package types

type (
	DmlMapping struct {
		ID              uint64         `json:"mappingID,string"`
		ConnectionID    uint64         `json:"connectionID,string"`
		NamespaceHandle string         `json:"namespaceHandle"`
		Tables          []*DmlTableMap `json:"tables"`
	}

	DmlTableMap struct {
		SourceIdent  string          `json:"sourceIdent"`
		ModuleHandle string          `json:"moduleHandle"`
		ModuleName   string          `json:"moduleName"`
		Skip         bool            `json:"skip"`
		Identifier   string          `json:"identifier"`
		Columns      []*DmlColumnMap `json:"columns"`
	}

	DmlColumnMap struct {
		SourceIdent string `json:"sourceIdent"`
		FieldName   string `json:"fieldName"`
		FieldKind   string `json:"fieldKind"`
		Skip        bool   `json:"skip"`
	}

	DmlImportRun struct {
		ID           uint64            `json:"runID,string"`
		ConnectionID uint64            `json:"connectionID,string"`
		MappingID    uint64            `json:"mappingID,string"`
		Status       string            `json:"status"`
		Processed    uint64            `json:"processed"`
		Failed       uint64            `json:"failed"`
		Error        string            `json:"error,omitempty"`
		Cursor       map[string]string `json:"cursor,omitempty"`
	}
)
