package types

type (
	ProjectChangeRisk string

	// ProjectChangeOp is what happened to a resource between two revisions.
	// Without it a change is ambiguous: an added and a removed module both
	// describe the same path, and a removed field and a retyped one both read
	// as "dangerous".
	ProjectChangeOp string

	// ProjectChange is one difference between a draft revision and its parent.
	// Op/Kind/Name are what a UI renders; Path is the stable machine identity.
	ProjectChange struct {
		Op   ProjectChangeOp `json:"op"`
		Kind string          `json:"kind"`
		Name string          `json:"name"`

		// Module owning a field-level change.
		Module string `json:"module,omitempty"`

		// Detail spells out what Op alone cannot, e.g. a field's "String → Email".
		Detail string `json:"detail,omitempty"`

		// Risk is about DATA, not importance: dangerous means records are lost
		// unless the change is mapped. Removing a page loses nothing and stays
		// safe; removing a field with records behind it does not.
		Risk ProjectChangeRisk `json:"risk"`

		// Records held by the module this change affects — what is at stake if
		// it goes unmapped. Best-effort: 0 means "unknown", never "none".
		Records uint64 `json:"records"`

		Path string `json:"path"`
	}

	ModuleFieldMapping struct {
		SourceField string `json:"sourceField"`
		TargetField string `json:"targetField"`
		Op          string `json:"op"`
		Value       string `json:"value,omitempty"`
		OnError     string `json:"onError,omitempty"`
	}

	ModuleMapping struct {
		Module string               `json:"module"`
		Fields []ModuleFieldMapping `json:"fields"`
	}

	ProjectDeploymentPlan struct {
		// Risk is the worst risk across Changes.
		Risk              ProjectChangeRisk `json:"risk"`
		Changes           []ProjectChange   `json:"changes"`
		SuggestedMappings []ModuleMapping   `json:"suggestedMappings"`
	}

	PublishRequest struct {
		Confirm bool `json:"confirm"`
		// Mappings the caller is explicit about. Anything absent is filled in
		// from the server's own deployment plan -- publishing used to accept an
		// empty set and migrate nothing, which quietly emptied every module in
		// the newly live revision.
		Mappings []ModuleMapping `json:"mappings"`
		// DiscardRecords is the way to actually publish without carrying
		// records over. It has to be asked for; it is not what silence means.
		DiscardRecords bool `json:"discardRecords"`
	}
)

// HasSource reports whether this mapping has anything to migrate FROM. A
// module the revision added has a mapping (so it shows up in the plan) whose
// fields either are absent or only carry defaults -- there is no source table
// behind it, and asking the importer to read one fails the whole publish.
func (m ModuleMapping) HasSource() bool {
	for _, f := range m.Fields {
		if f.SourceField != "" {
			return true
		}
	}

	return false
}

const (
	ProjectChangeRiskSafe      ProjectChangeRisk = "safe"
	ProjectChangeRiskDangerous ProjectChangeRisk = "dangerous"

	ProjectChangeOpAdded   ProjectChangeOp = "added"
	ProjectChangeOpRemoved ProjectChangeOp = "removed"
	ProjectChangeOpChanged ProjectChangeOp = "changed"
)
