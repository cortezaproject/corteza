package types

import (
	"errors"

	compTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
	"gopkg.in/yaml.v3"
)

type (
	ComposeRecordValue struct {
		compTypes.RecordValue `yaml:",inline"`
	}
	ComposeRecordValueSet []*ComposeRecordValue

	ComposeRecord struct {
		compTypes.Record `yaml:",inline"`
		Values           ComposeRecordValueSet `yaml:"values"`
	}
	ComposeRecordSet []*ComposeRecord
)

var (
	ErrRecordDecoderMapNotSupported = errors.New("record decoder: maps not supported")
)

func (rr *ComposeRecordSet) UnmarshalYAML(n *yaml.Node) error {
	crs := ComposeRecordSet{}
	if n.Kind != yaml.SequenceNode {
		return ErrRecordDecoderMapNotSupported
	}

	err := util.YamlIterator(n, func(n, m *yaml.Node) error {
		var r *ComposeRecord
		err := m.Decode(&r)
		if err != nil {
			return err
		}

		crs = append(crs, r)
		return nil
	})

	if err != nil {
		return err
	}

	*rr = crs
	return nil
}

func (vv *ComposeRecordValueSet) UnmarshalYAML(n *yaml.Node) error {
	vvs := ComposeRecordValueSet{}

	err := util.YamlIterator(n, func(n, m *yaml.Node) error {
		v := &ComposeRecordValue{}

		v.Name = n.Value
		v.Value = m.Value
		v.Updated = true
		vvs = append(vvs, v)
		return nil
	})

	if err != nil {
		return err
	}

	*vv = vvs
	return nil
}
