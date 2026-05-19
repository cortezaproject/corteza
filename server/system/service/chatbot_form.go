package service

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/system/types"
)

// ChatbotScenarioConfig returns the raw scenario config decoded as a generic
// JSON value (object) so it can be ferried into SSE payloads. Returns nil if
// the config is empty or malformed.
func ChatbotScenarioConfig(s *types.ChatbotScenario) any {
	if len(s.Config) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(s.Config, &v); err != nil {
		return nil
	}
	return v
}

// ValidateChatbotFormFields enforces the "required" attribute on form-type
// scenarios. Returns a map of fieldName→error. An empty map means valid.
func ValidateChatbotFormFields(rawConfig json.RawMessage, values map[string]string) map[string]string {
	errs := map[string]string{}
	if len(rawConfig) == 0 {
		return errs
	}
	var cfg struct {
		Fields []struct {
			Name     string `json:"name"`
			Label    string `json:"label"`
			Type     string `json:"type"`
			Required bool   `json:"required"`
			Pattern  string `json:"pattern,omitempty"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(rawConfig, &cfg); err != nil {
		return errs
	}
	for _, f := range cfg.Fields {
		v := values[f.Name]
		if f.Required && v == "" {
			errs[f.Name] = "required"
		}
	}
	return errs
}

// EncodeChatbotFormSubmission packs a form submission into a deterministic
// string suitable for a chat history message.
func EncodeChatbotFormSubmission(values map[string]string) string {
	if len(values) == 0 {
		return "[form submitted]"
	}
	b, err := json.Marshal(values)
	if err != nil {
		return "[form submitted]"
	}
	return "[form submitted] " + string(b)
}

// ChatbotVarsFromMap converts a flat string map into expr.Vars usable as
// automation input. nil-safe; an empty map yields an empty vars set.
func ChatbotVarsFromMap(m map[string]string) *expr.Vars {
	v := &expr.Vars{}
	for k, val := range m {
		_ = v.Set(k, val)
	}
	return v
}
