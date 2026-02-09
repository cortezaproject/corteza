package types

type (
	// Trigger definitions
	ConstructTrigger struct {
		ResourceType string                       `json:"resourceType"`
		EventType    string                       `json:"eventType"`
		Properties   []ConstructTriggerProperty   `json:"properties"`
		Constraints  []ConstructTriggerConstraint `json:"constraints"`
	}

	ConstructTriggerProperty struct {
		Name string                       `json:"name"`
		Type string                       `json:"type"`
		Meta ConstructTriggerPropertyMeta `json:"meta,omitempty"`
	}

	ConstructTriggerPropertyMeta struct {
		Short       string `json:"short,omitempty"`
		Description string `json:"description,omitempty"`
	}

	ConstructTriggerConstraint struct {
		Name string                         `json:"name"`
		Type string                         `json:"type"`
		Meta ConstructTriggerConstraintMeta `json:"meta,omitempty"`
	}

	ConstructTriggerConstraintMeta struct {
		Short       string `json:"short,omitempty"`
		Description string `json:"description,omitempty"`
	}

	// Function definitions
	ConstructFunction struct {
		Ref  string                 `json:"ref,omitempty"`
		Kind string                 `json:"kind,omitempty"`
		Meta *ConstructFunctionMeta `json:"meta,omitempty"`

		Parameters ParamSet `json:"parameters,omitempty"`
		Results    ParamSet `json:"results,omitempty"`

		Segments []ConstructSegment `json:"segments"`

		ArgsMerger   FunctionMerger  `json:"-"`
		Handler  FunctionHandler `json:"-"`
		Iterator IteratorHandler `json:"-"`

		Labels   map[string]string `json:"labels,omitempty"`
		Disabled bool              `json:"disabled,omitempty"`
	}

	ConstructFunctionMeta struct {
		Short       string `json:"short,omitempty"`
		Description string `json:"description,omitempty"`
	}

	// UI structure
	ConstructSegment struct {
		Meta     ConstructSegmentMeta `json:"meta,omitempty"`
		Sections []ConstructSection   `json:"sections,omitempty"`
	}

	ConstructSegmentMeta struct {
		Short       string `json:"short,omitempty"`
		Description string `json:"description,omitempty"`
	}

	ConstructSection struct {
		Meta     ConstructSectionMeta `json:"meta,omitempty"`
		Elements []SectionElement     `json:"elements,omitempty"`
	}

	ConstructSectionMeta struct {
		Short       string `json:"short,omitempty"`
		Description string `json:"description,omitempty"`
	}

	SectionElement struct {
		Action struct{}            `json:"action,omitempty"`
		Input  SectionElementInput `json:"input,omitempty"`
	}

	SectionElementInput struct {
		Type        string `json:"type,omitempty"`
		Label       string `json:"label,omitempty"`
		Placeholder string `json:"placeholder,omitempty"`
		Argument    string `json:"argument,omitempty"`

		// Context defines data dependencies between inputs
		Context SectionElementInputContext `json:"context,omitempty"`
		// Visual is for styling/display options
		Visual struct{} `json:"visual,omitempty"`
	}

	SectionElementInputContext struct {
		// DependsOn maps prop names to source argument names
		// e.g., {"namespaceID": "namespace"} means this input needs
		// the value of "namespace" argument passed as "namespaceID" prop
		DependsOn map[string]string `json:"dependsOn,omitempty"`
	}
)
