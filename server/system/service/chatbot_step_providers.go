package service

import (
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/system/types"
)

// ChatbotStepProperty describes a source path a step type exposes for hook mappings.
// Path is relative to steps.<scenarioID> — FE prefixes with the actual scenarioID.
type ChatbotStepProperty struct {
	Path        string `json:"path"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type chatbotStepDescriptor struct {
	Properties []ChatbotStepProperty
	Resolve    func(s *types.ChatbotSessionStepState) *expr.Vars
}

var chatbotStepProviders = map[string]chatbotStepDescriptor{
	"form": {
		Properties: []ChatbotStepProperty{
			{Path: "form.submitted", Type: "Boolean", Description: "Whether form was submitted"},
			{Path: "form.fields", Type: "", Description: "Map of submitted field values, keyed by field name"},
		},
		Resolve: func(s *types.ChatbotSessionStepState) *expr.Vars {
			v := &expr.Vars{}
			if s.Form == nil {
				return v
			}
			fv := &expr.Vars{}
			_ = fv.Set("submitted", s.Form.Submitted)
			fieldsVars := &expr.Vars{}
			for k, val := range s.Form.Fields {
				_ = fieldsVars.Set(k, val)
			}
			_ = fv.Set("fields", fieldsVars)
			_ = v.Set("form", fv)
			return v
		},
	},
	"conversation": {
		Properties: []ChatbotStepProperty{
			{Path: "conversation.conversationID", Type: "ID"},
			{Path: "conversation.handoff.operatorID", Type: "ID", Description: "Operator who accepted handoff"},
			{Path: "conversation.handoff.operatorName", Type: "String", Description: "Operator display name"},
		},
		Resolve: func(s *types.ChatbotSessionStepState) *expr.Vars {
			v := &expr.Vars{}
			if s.Conversation == nil {
				return v
			}
			cv := &expr.Vars{}
			_ = cv.Set("conversationID", s.Conversation.ConversationID)
			if s.Conversation.Handoff != nil {
				hv := &expr.Vars{}
				_ = hv.Set("operatorID", s.Conversation.Handoff.OperatorID)
				_ = hv.Set("operatorName", s.Conversation.Handoff.OperatorName)
				_ = cv.Set("handoff", hv)
			}
			_ = v.Set("conversation", cv)
			return v
		},
	},
}

// StepTypeProperties returns the registered property list per step type.
// Consumed by the REST API to expose source paths to the FE mapping UI.
func StepTypeProperties() map[string][]ChatbotStepProperty {
	out := make(map[string][]ChatbotStepProperty, len(chatbotStepProviders))
	for t, d := range chatbotStepProviders {
		out[t] = d.Properties
	}
	return out
}

// resolveStepSource resolves a StateExpression Source path against session state.
// Format: "steps.<scenarioID>.<relativePath>" — provider for that step type is called lazily.
// phase="before" blocks resolution of currentScenarioID (step not yet complete).
func resolveStepSource(source string, session *types.ChatbotSession, currentScenarioID, phase string) (interface{}, error) {
	const prefix = "steps."
	if !strings.HasPrefix(source, prefix) {
		return nil, fmt.Errorf("unsupported source path %q: must start with %q", source, prefix)
	}
	rest := source[len(prefix):]
	dotIdx := strings.IndexByte(rest, '.')
	if dotIdx < 0 {
		return nil, fmt.Errorf("source %q: expected steps.<scenarioID>.<path>", source)
	}
	scenarioID := rest[:dotIdx]
	relPath := rest[dotIdx+1:]

	if phase == "before" && scenarioID == currentScenarioID {
		return nil, fmt.Errorf("source %q: step %q not yet complete at phase 'before'", source, scenarioID)
	}

	state := session.State.ForScenario(scenarioID)
	if state == nil {
		return nil, fmt.Errorf("source %q: no state found for scenario %q", source, scenarioID)
	}

	desc, ok := chatbotStepProviders[state.Type]
	if !ok {
		return nil, fmt.Errorf("source %q: no provider registered for step type %q", source, state.Type)
	}

	stepVars := desc.Resolve(state)
	return selectScopePath(stepVars, relPath)
}

// selectScopePath resolves a dot-separated path against a nested *expr.Vars.
func selectScopePath(scope *expr.Vars, path string) (interface{}, error) {
	parts := strings.SplitN(path, ".", 2)
	top, err := scope.Select(parts[0])
	if err != nil {
		return nil, fmt.Errorf("no value at %q", path)
	}
	if len(parts) == 1 {
		return top, nil
	}
	nested, ok := top.(*expr.Vars)
	if !ok {
		return nil, fmt.Errorf("cannot traverse %q: not a nested object", parts[0])
	}
	return selectScopePath(nested, parts[1])
}
