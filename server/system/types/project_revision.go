package types

type (
	ProjectChangeRisk string

	ProjectChange struct {
		Path string            `json:"path"`
		Risk ProjectChangeRisk `json:"risk"`
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
		Path              string          `json:"path"`
		Changes           []ProjectChange `json:"changes"`
		SuggestedMappings []ModuleMapping `json:"suggestedMappings"`
	}

	PublishRequest struct {
		Confirm  bool            `json:"confirm"`
		Mappings []ModuleMapping `json:"mappings"`
	}
)

const (
	ProjectChangeRiskSafe      ProjectChangeRisk = "safe"
	ProjectChangeRiskDangerous ProjectChangeRisk = "dangerous"
)
