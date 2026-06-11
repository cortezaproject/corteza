package resourceref

import (
	"strconv"
)

type (
	// Ref describes a configuration-level reference from one resource to another
	//
	// Extractors living next to resource type definitions emit Refs; consumers
	// (dependency graph assembly, envoy decoders, ...) resolve and interpret them.
	Ref struct {
		// Kind of the referenced resource (see Kind* constants)
		Kind string

		// ID of the referenced resource when the configuration holds a concrete ID
		ID uint64

		// Ident holds a textual identifier (handle) when the configuration
		// references the resource by something other than an ID
		Ident string

		// Reason classifies the dependency (see Reason* constants)
		Reason string

		// Path points to the location inside the source resource's
		// configuration where the reference was found
		Path string

		// Dynamic marks a reference whose target is only known at runtime
		// (computed step arguments, scope variables); no ID/Ident is set and
		// consumers should surface these as warnings, not edges
		Dynamic bool
	}
)

// Resource kinds; values mirror the generated *ResourceType constants
// (see {compose,system,automation}/types/resources.gen.go).
//
// Duplicated here so resource type packages can emit cross-component refs
// without importing each other (would cause import cycles; e.g. system/types
// can not import automation/types).
const (
	KindComposeNamespace = "corteza::compose:namespace"
	KindComposeModule    = "corteza::compose:module"
	KindComposeChart     = "corteza::compose:chart"
	KindComposePage      = "corteza::compose:page"

	KindAutomationWorkflow = "corteza::automation:workflow"
	KindNgAutomation       = "corteza::automation:ng-automation"

	KindDalConnection        = "corteza::system:dal-connection"
	KindConfiguredConnection = "corteza::system:configured-connection"
	KindLlmProvider          = "corteza::system:llm-provider"
	KindKnowledgeBase        = "corteza::system:knowledge-base"
	KindAgent                = "corteza::system:agent"
	KindRole                 = "corteza::system:role"
	KindTemplate             = "corteza::system:template"
)

// Dependency reasons
const (
	ReasonModuleFieldRef      = "module-field-ref"
	ReasonModuleConnection    = "module-connection"
	ReasonPageModule          = "page-module"
	ReasonPageChart           = "page-chart"
	ReasonPageWorkflow        = "page-workflow"
	ReasonChartModule         = "chart-module"
	ReasonTriggerModule       = "trigger-module"
	ReasonTriggerWorkflow     = "trigger-workflow"
	ReasonAgentKnowledgeBase  = "agent-knowledge-base"
	ReasonAgentLlmProvider    = "agent-llm-provider"
	ReasonAgentAutomation     = "agent-automation"
	ReasonAgentModule         = "agent-module"
	ReasonChatbotAgent        = "chatbot-agent"
	ReasonChatbotAutomation   = "chatbot-automation"
	ReasonKnowledgeBaseModule = "knowledge-base-module"
	ReasonStepArgument        = "step-argument"
	ReasonStepConnection      = "step-connection"
)

// Make returns an ID-based Ref; zero ID yields an empty Ref (dropped by Append)
func Make(kind string, id uint64, reason, path string) Ref {
	if id == 0 {
		return Ref{}
	}

	return Ref{Kind: kind, ID: id, Reason: reason, Path: path}
}

// MakeIdent returns a Ref from a textual identifier; numeric identifiers are
// parsed into the ID field, anything else is kept as Ident
func MakeIdent(kind, ident, reason, path string) Ref {
	if ident == "" || ident == "0" {
		return Ref{}
	}

	if id, err := strconv.ParseUint(ident, 10, 64); err == nil {
		return Ref{Kind: kind, ID: id, Reason: reason, Path: path}
	}

	return Ref{Kind: kind, Ident: ident, Reason: reason, Path: path}
}

// MakeDynamic returns a Ref for a runtime-resolved reference; kind may be
// empty when even the target kind can not be determined from the config
func MakeDynamic(kind, reason, path string) Ref {
	return Ref{Kind: kind, Reason: reason, Path: path, Dynamic: true}
}

// Append appends the given refs to out, dropping empty ones
func Append(out []Ref, rr ...Ref) []Ref {
	for _, r := range rr {
		if r.IsEmpty() {
			continue
		}

		out = append(out, r)
	}

	return out
}

func (r Ref) IsEmpty() bool {
	return r.ID == 0 && r.Ident == "" && !r.Dynamic
}
