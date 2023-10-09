package envoyx

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	Provider interface {
		Next(ctx context.Context, out map[string]string) (more bool, err error)
		Reset(ctx context.Context) error
		SetIdent(string)
		Ident() string
	}

	Datasource interface {
		Next(ctx context.Context, out map[string]string) (ident []string, more bool, err error)
		Reset(ctx context.Context) error
		SetProvider(Provider) bool
	}

	MapEntry struct {
		Column string
		Field  string
		Skip   bool
	}

	FieldMapping struct {
		// @note This had to be like so to simplify decoding
		Map map[string]MapEntry
	}

	DatasourceMapping struct {
		SourceIdent string   `yaml:"source"`
		KeyField    []string `yaml:"key"`
		References  map[string]string
		Scope       map[string]string

		// Defaultable indicates wether the mapping should keep the values where
		// the ident is not explicitly mapped.
		//
		// When true, the value is assigned to the given identifier.
		Defaultable bool `yaml:"defaultable"`
		Mapping     FieldMapping
	}
)

func SetDecoderSources(nn NodeSet, dd ...Provider) {
	for _, n := range nn {
		if n.Datasource == nil {
			continue
		}

		for _, d := range dd {
			if n.Datasource.SetProvider(d) {
				break
			}
		}
	}
}

// UnmarshalYAML is used to get the yaml parsed into a series of nodes so
// we can easily pass it down
func (d *FieldMapping) UnmarshalYAML(n *yaml.Node) (err error) {
	d.Map = make(map[string]MapEntry)
	if y7s.IsSeq(n) {
		err = y7s.EachSeq(n, func(n *yaml.Node) error {
			a, err := d.unmarshalMappingNode(n)
			d.Map[a.Column] = a
			return err
		})
	} else {
		err = y7s.EachMap(n, func(k, n *yaml.Node) error {
			a, err := d.unmarshalMappingNode(n)
			if a.Column == "" {
				err = y7s.DecodeScalar(k, "fieldMapping column", &a.Column)
				if err != nil {
					return err
				}
			}

			d.Map[a.Column] = a
			return err
		})
	}

	return
}

func (d *FieldMapping) unmarshalMappingNode(n *yaml.Node) (out MapEntry, err error) {
	if y7s.IsKind(n, yaml.ScalarNode) {
		err = y7s.DecodeScalar(n, "Column", &out.Column)
		if err != nil {
			return
		}
		err = y7s.DecodeScalar(n, "Field", &out.Field)
		return
	}

	// @todo we're omitting errors because there will be a bunch due to invalid
	//       resource field types. This might be a bit unstable as other errors may
	//       also get ignored.
	//
	//       A potential fix would be to firstly unmarshal into an any, check errors
	//       and then unmarshal into the resource while omitting errors.
	n.Decode(&out)

	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch strings.ToLower(k.Value) {
		case "skip":
			if v.Value == "/" {
				out.Skip = true
			}
		}
		return nil
	})

	return
}
