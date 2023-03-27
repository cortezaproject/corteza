package envoyx

import (
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/spf13/cast"
)

type (
	// Node is a wrapper around a Corteza resource for use within Envoy
	Node struct {
		Resource   resource
		Datasource Datasource

		ResourceType string
		Identifiers  Identifiers
		References   map[string]Ref
		Scope        Scope

		// Placeholders are resources which were added to help resolve missing deps
		Placeholder bool
		Config      EnvoyConfig
		Evaluated   Evaluated
	}

	Evaluated struct {
		Skip bool
	}

	EnvoyConfig struct {
		MergeAlg   mergeAlg
		SkipIf     string
		SkipIfEval expr.Evaluable
	}

	NodeSet []*Node

	resource interface {
		SetValue(name string, pos uint, value any) error
		GetValue(name string, pos uint) (any, error)
		GetID() uint64
	}

	Identifiers struct {
		Slice []string
		Index map[string]bool
	}

	// Scope lets us group nodes based on some common context
	//
	// Scope is primarily used to scope low-code applications to denote to what
	// namespace a specific module reference belongs to.
	// In the previous version this was referred to as reference constraints;
	// This is the same but different.
	//
	// When constructing dependency graphs, nodes with the same scope are grouped together.
	// Nodes from the same scope can reference each other.
	// Nodes from a defined scope can reference nodes from an undefined scope,
	// but not the other way around.
	Scope struct {
		ResourceType string
		Identifiers  Identifiers
	}

	// Ref defines a reference to a different resource
	//
	// The reference only holds if all three parts match -- the resource type,
	// there is an intersection between the identifiers, and the scope matches.
	Ref struct {
		ResourceType string
		Identifiers  Identifiers
		Scope        Scope
		// @todo consider replacing with something that indicates
		//       it can't be fetched from the DB
		Optional bool
	}

	prunner interface {
		Prune(rt string)
	}
)

// MakeIdentifiers initializes an Identifiers instance from the given slice
func MakeIdentifiers(ii ...any) (out Identifiers) {
	return Identifiers{}.Add(ii...)
}

// NodesByResourceType returns Nodes grouped by their resource type
func NodesByResourceType(nn ...*Node) (out map[string]NodeSet) {
	out = make(map[string]NodeSet, len(nn)/2)
	for _, n := range nn {
		out[n.ResourceType] = append(out[n.ResourceType], n)
	}

	return
}

// NodesForResourceType returns which belong to the given resource type
func NodesForResourceType(rt string, nn ...*Node) (out NodeSet) {
	return NodesByResourceType(nn...)[rt]
}

// NodeForRef returns the Node that matches the given ref
func NodeForRef(ref Ref, nn ...*Node) (out *Node) {
	for _, n := range nn {
		if !n.Scope.Equals(ref.Scope) {
			continue
		}

		if n.ResourceType != ref.ResourceType {
			continue
		}

		if n.Identifiers.HasIntersection(ref.Identifiers) {
			return n
		}
	}

	return
}

func OmitPlaceholderNodes(nn ...*Node) (out NodeSet) {
	out = make(NodeSet, 0, len(nn))
	for _, n := range nn {
		if n.Placeholder {
			continue
		}
		out = append(out, n)
	}
	return
}

func MergeRefs(a map[string]Ref, bb ...map[string]Ref) (c map[string]Ref) {
	c = make(map[string]Ref)

	for k, v := range a {
		c[k] = v
	}

	for _, b := range bb {
		for k, v := range b {
			c[k] = v
		}
	}

	return
}

// MergeIdents merges the two identifiers and returns a new one
func MergeIdents(a, b Identifiers) (cc Identifiers) {
	cc = Identifiers{
		Index: make(map[string]bool, 2),
	}

	for a := range a.Index {
		cc.Index[a] = true
	}
	for b := range b.Index {
		cc.Index[b] = true
	}

	for c := range cc.Index {
		cc.Slice = append(cc.Slice, c)
	}

	return
}

func (n Node) ToRef() Ref {
	return Ref{
		ResourceType: n.ResourceType,
		Identifiers:  n.Identifiers,
		Scope:        n.Scope,
	}
}

func (n *Node) Prune(ref Ref) {
	for k, nodeRef := range n.References {
		if nodeRef.Equals(ref) {
			delete(n.References, k)
		}
	}

	if pp, ok := n.Resource.(prunner); ok {
		// @todo improve when needed; for now it's ok
		pp.Prune(ref.ResourceType)
	}
}

func (r Ref) Idents() (ints []uint64, rest []string) {
	return r.Identifiers.Idents()
}

// ResourceFilter returns a filter which would match the referenced resource
func (r Ref) ResourceFilter() (out map[string]ResourceFilter) {
	out = make(map[string]ResourceFilter)
	out[r.ResourceType] = ResourceFilter{
		Identifiers: r.Identifiers,
		Scope:       r.Scope,

		// A ref would point to a single resource.
		// Don't set the limit so we can error out on ambiguity.
		// Limit: 1,
	}

	return
}

func (a Ref) Equals(b Ref) bool {
	if a.ResourceType != b.ResourceType {
		return false
	}
	if !a.Scope.Equals(b.Scope) {
		return false
	}
	if !a.Identifiers.HasIntersection(b.Identifiers) {
		return false
	}

	return true
}

func (a Scope) Equals(b Scope) bool {
	if a.IsEmpty() && b.IsEmpty() {
		return true
	}

	if a.ResourceType != b.ResourceType {
		return false
	}

	return a.Identifiers.HasIntersection(b.Identifiers)
}

func (s Scope) IsEmpty() bool {
	return s.ResourceType == "" && len(s.Identifiers.Slice) == 0
}

// Add adds the given values to the identifier
func (ii Identifiers) Add(vv ...any) (out Identifiers) {
	if ii.Index == nil {
		ii.Index = make(map[string]bool, len(vv))
		ii.Slice = make([]string, 0, len(vv))
	}

	for _, v := range vv {
		switch casted := v.(type) {
		case string:
			if casted == "" {
				continue
			}
		case uint64, uint, int, int64:
			if casted == 0 {
				continue
			}
		}

		if c, ok := v.(Identifiers); ok {
			ii = ii.Merge(c)
			continue
		}

		aux := cast.ToString(v)
		if aux == "" || aux == "0" {
			continue
		}

		ii.Slice = append(ii.Slice, aux)
		ii.Index[aux] = true
	}

	return ii
}

func (ii Identifiers) IdentsAsStrings() (ids, rest []string) {
	aux, rest := ii.Idents()

	for _, a := range aux {
		ids = append(ids, strconv.FormatUint(a, 10))
	}

	return
}

// Idents returns a slice of numeric and text identifiers
func (ii Identifiers) Idents() (ints []uint64, rest []string) {
	var aux uint64
	var err error

	for _, i := range ii.Slice {
		aux, err = cast.ToUint64E(i)
		if err != nil {
			rest = append(rest, i)
		} else {
			ints = append(ints, aux)
		}
	}

	return
}

// Intersection returns a slice of identifiers which are in an intersection
func (aa Identifiers) Intersection(bb Identifiers) (out []string) {
	for _, b := range bb.Slice {
		if aa.Index[b] {
			out = append(out, b)
		}
	}

	return
}

// HasIntersection returns true if the two identifiers define an intersection
func (aa Identifiers) HasIntersection(bb Identifiers) bool {
	if len(aa.Slice) == 0 && len(bb.Slice) == 0 {
		return true
	}

	// If the second identifies are empty we can assume that anything has an intersection with them.
	// This is useful when we want to remove any resource of some type.
	if len(bb.Slice) == 0 {
		return true
	}

	return len(aa.Intersection(bb)) > 0
}

// Merge merges the two identifiers and returns a new one
// @todo deprecate this; use MergeIdents instead
func (aa Identifiers) Merge(bb Identifiers) (cc Identifiers) {
	return MergeIdents(aa, bb)
}

// FriendlyIdentifier returns the best available identifier
//
// If any non-ID identifiers are available, it uses the first one.
// If no non-ID identifiers are available, it returns the first ID.
func (ii Identifiers) FriendlyIdentifier() (out string) {
	a, b := ii.Idents()
	if len(b) > 0 {
		return b[0]
	}

	if len(a) > 0 {
		return strconv.FormatUint(a[0], 10)
	}

	return
}
