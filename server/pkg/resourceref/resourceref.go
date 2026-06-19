package resourceref

import (
	"strconv"
	"strings"
)

type Ref struct {
	// "corteza::component:kind/id" for ID refs, "corteza::component:kind" for label/unresolved
	Resource string

	// textual handle when the config references by something other than an ID
	Label string

	Reason string

	// target only known at runtime; surface as warning, not edge
	Unresolved bool

	// Wildcard ref: connects to every resource of Resource (kind). Used by
	// RBAC-derived role edges where a rule grants/denies on a whole kind
	// (e.g. "corteza::compose:module/*/*"). Resolved by fanning out to all
	// in-scope nodes of the kind, not a single target.
	Wildcard bool
}

// Kind constants mirror *ResourceType constants in {compose,system,automation}/types/resources.gen.go.
// Duplicated here to avoid import cycles between component type packages.
const (
	KindComposeNamespace  = "corteza::compose:namespace"
	KindComposeModule     = "corteza::compose:module"
	KindComposeChart      = "corteza::compose:chart"
	KindComposePage       = "corteza::compose:page"
	KindComposePageLayout = "corteza::compose:page-layout"

	KindAutomationWorkflow = "corteza::automation:workflow"
	KindNgAutomation       = "corteza::automation:ng-automation"

	KindDalConnection        = "corteza::system:dal-connection"
	KindConfiguredConnection = "corteza::system:configured-connection"
	KindLlmProvider          = "corteza::system:llm-provider"
	KindKnowledgeBase        = "corteza::system:knowledge-base"
	KindAgent                = "corteza::system:agent"
	KindChatbot              = "corteza::system:chatbot"
	KindRole                 = "corteza::system:role"
	KindTemplate             = "corteza::system:template"
)

const (
	ReasonModuleFieldRef      = "module-field-ref"
	ReasonModuleConnection    = "module-connection"
	ReasonPageModule          = "page-module"
	ReasonPageChart           = "page-chart"
	ReasonPageWorkflow        = "page-workflow"
	ReasonPageAutomation      = "page-automation"
	ReasonPageAgent           = "page-agent"
	ReasonPageChatbot         = "page-chatbot"
	ReasonPageNavigation      = "page-navigation"
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
	ReasonRoleRbac            = "role-rbac"
)

func Make(kind string, id uint64, reason string) Ref {
	if id == 0 {
		return Ref{}
	}
	return Ref{Resource: kind + "/" + strconv.FormatUint(id, 10), Reason: reason}
}

func MakeIdent(kind, ident, reason string) Ref {
	if ident == "" || ident == "0" {
		return Ref{}
	}
	if id, err := strconv.ParseUint(ident, 10, 64); err == nil {
		return Ref{Resource: kind + "/" + strconv.FormatUint(id, 10), Reason: reason}
	}
	return Ref{Resource: kind, Label: ident, Reason: reason}
}

func MakeDynamic(kind, reason string) Ref {
	return Ref{Resource: kind, Reason: reason, Unresolved: true}
}

// MakeWildcard builds a ref that targets every resource of the given kind.
func MakeWildcard(kind, reason string) Ref {
	return Ref{Resource: kind, Reason: reason, Wildcard: true}
}

func Append(out []Ref, rr ...Ref) []Ref {
	for _, r := range rr {
		if r.IsEmpty() {
			continue
		}
		out = append(out, r)
	}
	return out
}

func (r Ref) Kind() string {
	if i := strings.Index(r.Resource, "/"); i >= 0 {
		return r.Resource[:i]
	}
	return r.Resource
}

func (r Ref) ID() uint64 {
	i := strings.LastIndex(r.Resource, "/")
	if i < 0 || i == len(r.Resource)-1 {
		return 0
	}
	id, _ := strconv.ParseUint(r.Resource[i+1:], 10, 64)
	return id
}

func (r Ref) IsEmpty() bool {
	return r.Resource == "" && r.Label == "" && !r.Unresolved && !r.Wildcard
}
