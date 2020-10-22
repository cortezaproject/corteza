package yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type (
	// decoder is a wrapper struct for yaml related methods
	decoder struct {
		loader loader
	}

	loader interface {
		LoadComposeNamespace()
	}

	nodeDecoder interface {
		DecodeNodes(ctx context.Context, l loader) ([]envoy.Node, error)
	}
)

func Decoder(l loader) *decoder {
	return &decoder{l}
}

// CanDecodeFile
func (y *decoder) CanDecodeFile(i os.FileInfo) bool {
	switch filepath.Ext(i.Name()) {
	case "yaml", "yml":
		return true
	}

	return false
}

func (y *decoder) Decode(ctx context.Context, r io.Reader, i os.FileInfo) ([]envoy.Node, error) {
	var (
		doc = &Document{}
	)

	if err := yaml.NewDecoder(r).Decode(doc); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", i.Name(), err)
	}

	return doc.Decode(ctx, y.loader)
}

//// convert converts the decoded document into a set of envoy nodes
//func (y *decoder) convert(c *Document) ([]envoy.Node, error) {
//	nn := make([]envoy.Node, 0, 100)
//
//	// In case of namespaces...
//	if c.Namespaces != nil {
//		nodes, err := y.convertNamespaces(c.Namespaces)
//		if err != nil {
//			return nil, err
//		}
//		nn = append(nn, nodes...)
//	}
//
//	ns := &types.Namespace{}
//	if c.Namespace != "" {
//		// In case of a namespace to provide dependencies
//		ns.Slug = c.Namespace
//		ns.Name = c.Namespace
//	} else if len(nn) > 0 {
//		// Try to fall back to a namespace node
//		ns = ((nn[0]).(*envoy.ComposeNamespaceNode)).Ns
//	} else {
//		// No good; we can't link with a namespace.
//		// @note This should be checked when converting Compose resources only.
//		//			 Some resources don't belong to a namespace.
//		spew.Dump(c)
//		return nil, errors.New("yaml decoder: cannot resolve namespace")
//	}
//
//	// In case of modules...
//	if c.Modules != nil {
//		nodes, err := y.convertModules(c.Modules, ns)
//		if err != nil {
//			return nil, err
//		}
//		nn = append(nn, nodes...)
//	}
//
//	if c.Records != nil {
//		for modRef, rr := range c.Records {
//			// We can define a basic module representation as it will be updated later
//			// during validation/runtime
//			mod := &types.module{}
//			mod.Handle = modRef
//			mod.Name = modRef
//
//			nodes, err := y.convertRecords(rr, mod)
//			if err != nil {
//				return nil, err
//			}
//			nn = append(nn, nodes...)
//		}
//	}
//
//	return nn, nil
//}

//func (y *decoder) convertNamespaces(nss ComposeNamespaceSet) ([]envoy.Node, error) {
//	nn := make([]envoy.Node, 0, 2)
//
//	for _, ns := range nss {
//		nn = append(nn, &envoy.ComposeNamespaceNode{Ns: ns.res})
//
//		// Nested modules
//		if ns.modules != nil {
//			mm, err := y.convertModules(ns.modules, ns.res)
//			if err != nil {
//				return nil, err
//			}
//			nn = append(nn, mm...)
//		}
//
//		// @todo nested RBAC
//	}
//
//	return nn, nil
//}
//
//func (y *decoder) convertModules(mm ComposeModuleSet, ns *types.Namespace) ([]envoy.Node, error) {
//	nn := make([]envoy.Node, 0)
//
//	for _, m := range mm {
//		nn = append(nn, &envoy.ComposeModuleNode{
//			module: m.res,
//			Ns:     ns,
//		})
//
//		// @todo nested resources; should there be any?
//	}
//
//	return nn, nil
//}
//
//func (y *decoder) convertRecords(rr ComposeRecordSet, m *types.Module) ([]envoy.Node, error) {
//	// Iterator function for providing records to be imported.
//	// This doesn't do any validation; that should be handled by other layers.
//	f := func(f func(record *types.Record) error) error {
//		//for _, r := range rr {
//		//err := f(r.res)
//		//if err != nil {
//		//	return err
//		//}
//		//}
//
//		return nil
//	}
//
//	return []envoy.Node{&envoy.ComposeRecordNode{Mod: m, Walk: f}}, nil
//}
