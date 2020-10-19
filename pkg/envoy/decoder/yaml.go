package decoder

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/types"
	"gopkg.in/yaml.v3"
)

type (
	// YamlDecoder is a wrapper struct for yaml related methods
	YamlDecoder struct{}

	// Document defines the supported yaml structure
	Document struct {
		Namespace  string
		Namespaces types.ComposeNamespaceSet
		Modules    types.ComposeModuleSet
		Records    map[string]types.ComposeRecordSet
	}
)

var (
	ErrorCannotResolveNamespace = errors.New("yaml decoder: cannot resolve namespace")
)

func NewYamlDecoder() *YamlDecoder {
	return &YamlDecoder{}
}

func (y *YamlDecoder) unmarshalDocument(r io.Reader) (*Document, error) {
	var c *Document

	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(buf.String()), &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// convert converts the decoded document into a set of envoy nodes
func (y *YamlDecoder) convert(c *Document) ([]types.Node, error) {
	nn := make([]types.Node, 0, 100)

	// In case of namespaces...
	if c.Namespaces != nil {
		nodes, err := y.convertNamespaces(c.Namespaces)
		if err != nil {
			return nil, err
		}
		nn = append(nn, nodes...)
	}

	ns := &types.ComposeNamespace{}
	if c.Namespace != "" {
		// In case of a namespace to provide dependencies
		ns.Slug = c.Namespace
		ns.Name = c.Namespace
	} else if len(nn) > 0 {
		// Try to fall back to a namespace node
		ns = ((nn[0]).(*types.ComposeNamespaceNode)).Ns
	} else {
		// No good; we can't link with a namespace.
		// @note This should be checked when converting Compose resources only.
		//			 Some resources don't belong to a namespace.
		return nil, ErrorCannotResolveNamespace
	}

	// In case of modules...
	if c.Modules != nil {
		nodes, err := y.convertModules(c.Modules, ns)
		if err != nil {
			return nil, err
		}
		nn = append(nn, nodes...)
	}

	if c.Records != nil {
		for modRef, rr := range c.Records {
			// We can define a basic module representation as it will be updated later
			// during validation/runtime
			mod := &types.ComposeModule{}
			mod.Handle = modRef
			mod.Name = modRef

			nodes, err := y.convertRecords(rr, mod)
			if err != nil {
				return nil, err
			}
			nn = append(nn, nodes...)
		}
	}

	return nn, nil
}

func (y *YamlDecoder) convertNamespaces(nss types.ComposeNamespaceSet) ([]types.Node, error) {
	nn := make([]types.Node, 0, 2)

	for _, ns := range nss {
		nn = append(nn, &types.ComposeNamespaceNode{Ns: ns})

		// Nested modules
		if ns.Modules != nil {
			mm, err := y.convertModules(ns.Modules, ns)
			if err != nil {
				return nil, err
			}
			nn = append(nn, mm...)
		}

		// @todo nested RBAC
	}

	return nn, nil
}

func (y *YamlDecoder) convertModules(mm types.ComposeModuleSet, ns *types.ComposeNamespace) ([]types.Node, error) {
	nn := make([]types.Node, 0)

	for _, m := range mm {
		nn = append(nn, &types.ComposeModuleNode{
			Mod: m,
			Ns:  ns,
		})

		// @todo nested resources; should there be any?
	}

	return nn, nil
}

func (y *YamlDecoder) convertRecords(rr types.ComposeRecordSet, m *types.ComposeModule) ([]types.Node, error) {
	// Iterator function for providing records to be imported.
	// This doesn't do any validation; that should be handled by other layers.
	f := func(f func(*types.ComposeRecord) error) error {
		for _, r := range rr {
			err := f(r)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return []types.Node{&types.ComposeRecordNode{Mod: m, Walk: f}}, nil
}

func (y *YamlDecoder) Decode(ctx context.Context, r io.Reader) ([]types.Node, error) {
	d, err := y.unmarshalDocument(r)
	if err != nil {
		return nil, err
	}

	return y.convert(d)
}
