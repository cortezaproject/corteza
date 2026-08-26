package types

type ProjectGraphNode struct {
	ID          uint64 `json:"id,string"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Sensitivity string `json:"sensitivity,omitempty"`

	// External marks nodes referenced from the project but living outside
	// its scope (tenant-level connections, LLM providers, ...)
	External bool `json:"external,omitempty"`
}

type ProjectGraphEdge struct {
	SourceID uint64 `json:"sourceID,string"`
	TargetID uint64 `json:"targetID,string"`
	Reason   string `json:"reason"`
}

// ProjectGraphMissingRef records a configured reference whose target could
// not be resolved (deleted resource, dangling handle, unknown kind)
type ProjectGraphMissingRef struct {
	SourceID    uint64 `json:"sourceID,string"`
	Kind        string `json:"kind"`
	TargetID    uint64 `json:"targetID,string,omitempty"`
	TargetIdent string `json:"targetIdent,omitempty"`
	Reason      string `json:"reason"`
}

// ProjectGraphWarning flags a configured reference which can only be resolved
// at runtime (computed step arguments, scope variables); shown on the source node
type ProjectGraphWarning struct {
	SourceID uint64 `json:"sourceID,string"`

	// Kind of the referenced resource, when determinable from config
	Kind   string `json:"kind,omitempty"`
	Reason string `json:"reason"`
	Path   string `json:"path,omitempty"`
}

type ProjectGraph struct {
	Nodes    []*ProjectGraphNode       `json:"nodes"`
	Edges    []*ProjectGraphEdge       `json:"edges"`
	Missing  []*ProjectGraphMissingRef `json:"missing,omitempty"`
	Warnings []*ProjectGraphWarning    `json:"warnings,omitempty"`
}
