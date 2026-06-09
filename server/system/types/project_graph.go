package types

type ProjectGraphNode struct {
	ID          uint64 `json:"id,string"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Sensitivity string `json:"sensitivity,omitempty"`
}

type ProjectGraphEdge struct {
	SourceID uint64 `json:"sourceID,string"`
	TargetID uint64 `json:"targetID,string"`
	Reason   string `json:"reason"`
}

type ProjectGraph struct {
	Nodes []*ProjectGraphNode `json:"nodes"`
	Edges []*ProjectGraphEdge `json:"edges"`
}
