package types

type (
	// ResponseMeta has all access restriction response
	ResponseMeta struct {
		Private   ResourceMeta `json:"private,omitempty"`
		Public    ResourceMeta `json:"public,omitempty"`
		Protected ResourceMeta `json:"protected,omitempty"`
	}

	// ResourceMeta have all resource related meta
	// 		for what to display in what order for each resource and its fields
	ResourceMeta struct {
		NamespaceMeta []NameMeta `json:"namespace_meta,omitempty"`
		ModuleMeta    []NameMeta `json:"module_meta,omitempty"`
		RecordMeta    []NameMeta `json:"record_meta,omitempty"`
	}

	// NameMeta is single row of discovery response fields with its weight
	NameMeta struct {
		Name   string `json:"name"`
		Title  string `json:"title"`
		Weight int    `json:"weight"`
	}
)
