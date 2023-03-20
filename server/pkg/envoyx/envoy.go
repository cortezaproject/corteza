package envoyx

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	Service struct {
		decoders map[decodeType][]Decoder

		encoders  map[encodeType][]Encoder
		preparers map[encodeType][]Preparer
	}

	// Traverser provides a structure which can be used to traverse the node's deps
	Traverser interface {
		// ParentForRef returns the parent of the provided node which matches the ref
		//
		// If no parent is found, nil is returned.
		ParentForRef(*Node, Ref) *Node

		// ParentForRT returns a set of parent nodes matching the resource type
		ParentForRT(*Node, string) NodeSet

		// ChildrenForResourceType returns the children of the provided node which
		// match the provided resource type
		ChildrenForResourceType(*Node, string) NodeSet

		// Children returns all of the children of the provided node
		Children(*Node) NodeSet

		// NodeForRef returns the node which matches the provided ref
		NodeForRef(Ref) *Node
	}

	Preparer interface {
		// Prepare performs generic preprocessing on the provided nodes
		//
		// The function is called for every resource type where all of the nodes of
		// that resource type are passed as the argument.
		Prepare(context.Context, EncodeParams, string, NodeSet) error
	}

	Encoder interface {
		// Encode encodes the data
		//
		// The function receives a set of root-level nodes (with no parent dependencies)
		// and a Traverser it can use to handle all of the child nodes.
		Encode(context.Context, EncodeParams, string, NodeSet, Traverser) (err error)
	}

	PrepareEncoder interface {
		Preparer
		Encoder
	}

	Decoder interface {
		// Decode returns a set of Nodes extracted based on the provided definition
		Decode(ctx context.Context, p DecodeParams) (out NodeSet, err error)
	}

	canCheckFile interface {
		CanFile(f *os.File) bool
	}

	canCheckExt interface {
		CanExt(name string) bool
	}

	DecodeParams struct {
		Type   decodeType
		Params map[string]any
		Config DecoderConfig
		Filter map[string]ResourceFilter
	}
	DecoderConfig struct{}

	EncodeParams struct {
		Type    encodeType
		Params  map[string]any
		Envoy   EnvoyConfig
		Encoder EncoderConfig

		// @note these are only used by records since v1 did just that
		DeferOk  func()
		DeferNok func(error) error
		Defer    func()
	}
	EncoderConfig struct {
		DefaultUserID       uint64
		PreferredTimeLayout string
		PreferredTimezone   string
	}

	ResourceFilter struct {
		Identifiers Identifiers
		Refs        map[string]Ref

		Limit uint
		Scope Scope
	}

	decodeType string
	encodeType string
	mergeAlg   int
)

var (
	global *Service

	ex = expr.NewParser()
)

const (
	OnConflictDefault mergeAlg = iota
	OnConflictReplace
	OnConflictSkip
	OnConflictPanic
	// OnConflictMergeLeft  mergeAlg = "mergeLeft"
	// OnConflictMergeRight mergeAlg = "mergeRight"

	DecodeTypeURI   decodeType = "uri"
	DecodeTypeIO    decodeType = "io"
	DecodeTypeStore decodeType = "store"

	EncodeTypeURI   encodeType = "uri"
	EncodeTypeStore encodeType = "store"
	EncodeTypeIo    encodeType = "io"
)

// New initializes a new Envoy service
func New() *Service {
	return &Service{}
}

// SetGlobal sets the global envoy service
func SetGlobal(n *Service) {
	global = n
}

// Global gets the global envoy service
func Global() *Service {
	if global == nil {
		panic("global service not defined")
	}

	return global
}

func Initialized() bool {
	return global != nil
}

// Decode returns a set of envoy Nodes based on the given decode params
func (svc *Service) Decode(ctx context.Context, p DecodeParams) (nodes NodeSet, providers []Provider, err error) {
	err = p.validate()
	if err != nil {
		return
	}

	switch p.Type {
	case DecodeTypeURI:
		return svc.decodeUri(ctx, p)
	case DecodeTypeIO:
		return svc.decodeIo(ctx, p)
	case DecodeTypeStore:
		nodes, err = svc.decodeStore(ctx, p)
		return
	default:
		err = fmt.Errorf("unsupported decoder type %s", p.Type)
	}

	return
}

func (svc *Service) Bake(ctx context.Context, p EncodeParams, providers []Provider, nodes ...*Node) (gg *DepGraph, err error) {
	err = svc.bakeEnvoyConfig(p.Envoy, nodes...)
	if err != nil {
		return
	}

	err = svc.bakeExpressions(nodes...)
	if err != nil {
		return
	}

	svc.bakeProviders(providers, nodes...)

	gg = BuildDepGraph(nodes...)
	return
}

// Encode encodes Corteza resources bases on the provided encode params
//
// use the BuildDepGraph function to build the default dependency graph.
func (svc *Service) Encode(ctx context.Context, p EncodeParams, dg *DepGraph) (err error) {
	err = p.validate()
	if err != nil {
		return
	}

	switch p.Type {
	case EncodeTypeStore:
		return svc.encodeStore(ctx, dg, p)
	case EncodeTypeIo:
		return svc.encodeIo(ctx, dg, p)
	default:
		err = fmt.Errorf("unsupported encoder type %s", p.Type)
	}
	return
}

func (svc *Service) AddDecoder(t decodeType, dd ...Decoder) {
	if svc.decoders == nil {
		svc.decoders = make(map[decodeType][]Decoder)
	}
	svc.decoders[t] = append(svc.decoders[t], dd...)
}

func (svc *Service) AddEncoder(t encodeType, ee ...Encoder) {
	if svc.encoders == nil {
		svc.encoders = make(map[encodeType][]Encoder)
	}
	if svc.preparers == nil {
		svc.preparers = make(map[encodeType][]Preparer)
	}
	svc.encoders[t] = append(svc.encoders[t], ee...)

	for _, e := range ee {
		p, ok := e.(Preparer)
		if !ok {
			continue
		}

		svc.preparers[t] = append(svc.preparers[t], p)
	}
}

func (svc *Service) AddPreparer(t encodeType, pp ...Preparer) {
	if svc.preparers == nil {
		svc.preparers = make(map[encodeType][]Preparer)
	}
	svc.preparers[t] = append(svc.preparers[t], pp...)
}

// Utility functions

func (p DecodeParams) validate() (err error) {
	switch p.Type {
	case DecodeTypeURI:
		_, ok := p.Params["uri"]
		if !ok {
			return fmt.Errorf("URI decoder expects a uri location parameter")
		}

	case DecodeTypeStore:

	}

	// @todo...

	return
}

func (p EncodeParams) validate() (err error) {
	switch p.Type {
	// @todo...
	}

	// @todo...

	return
}

func CastMergeAlg(v string) (mergeAlg mergeAlg) {
	switch strings.ToLower(v) {
	case "replace", "mergeleft":
		mergeAlg = OnConflictReplace
	case "skip", "mergeright":
		mergeAlg = OnConflictSkip
	case "panic", "error":
		mergeAlg = OnConflictPanic
	}

	return
}

func (svc *Service) bakeEnvoyConfig(dft EnvoyConfig, nodes ...*Node) (err error) {
	for _, n := range nodes {
		n.Config = svc.mergeEnvoyConfigs(n.Config, dft)
	}

	return
}

func (svc *Service) bakeExpressions(nodes ...*Node) (err error) {
	for _, n := range nodes {
		if n.Config.SkipIf == "" {
			continue
		}

		n.Config.SkipIfEval, err = ex.Parse(n.Config.SkipIf)
		if err != nil {
			return
		}
	}

	return
}

func (svc *Service) bakeProviders(providers []Provider, nodes ...*Node) {
	SetDecoderSources(nodes, providers...)
}

func (svc *Service) mergeEnvoyConfigs(a, b EnvoyConfig) (c EnvoyConfig) {
	c = a
	if c.MergeAlg == OnConflictDefault {
		c.MergeAlg = b.MergeAlg
	}
	if c.MergeAlg == OnConflictDefault {
		c.MergeAlg = OnConflictReplace
	}

	// @todo pre eval this?
	if c.SkipIf == "" {
		c.SkipIf = b.SkipIf
	}

	return
}
