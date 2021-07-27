package yaml

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// cleanMap helper removes any empty k:v nodes from the mapping node
//
// The value is empty when the tag is !!null OR when the value and the content are empty
func cleanMap(n *yaml.Node) *yaml.Node {
	cc := make([]*yaml.Node, 0, len(n.Content))

	for i := 0; i < len(n.Content)-1; i += 2 {
		k := n.Content[i]
		v := n.Content[i+1]

		if v.Tag == "!!null" || v.Value == "" && (len(v.Content) == 0 || v.Content == nil) {
			continue
		}

		cc = append(cc, k, v)
	}

	n.Content = cc
	return n
}

// encodeResource encodes the provided struct v as the resource r
//
// The helper has built-in sequence -> map transformation
func encodeResource(n *yaml.Node, r string, v interface{}, mapped bool, k string) (*yaml.Node, error) {
	var err error

	if mapped {
		m, err := mapResource(v, k)
		if err != nil {
			return nil, err
		}

		n, err = addMap(n,
			r, m,
		)
		if err != nil {
			return nil, err
		}
	} else {
		n, err = addMap(n,
			r, v,
		)
		if err != nil {
			return nil, err
		}
	}

	return n, nil
}

// mapResource encodes the resource as a map instead of a sequence
func mapResource(v interface{}, key string) (*yaml.Node, error) {
	var m *yaml.Node
	m, err := encodeNode(v)
	if err != nil {
		return nil, err
	}

	// It's already a mapping node; don't touch it
	if m.Kind == yaml.MappingNode {
		return m, nil
	}

	if m.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("expecting mapping or sequence node (%s provided)", m.Tag)
	}

	m, err = seqToMap(m, key)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// encodeNode shortens raw node encoding
func encodeNode(v interface{}) (*yaml.Node, error) {
	n := &yaml.Node{}
	err := n.Encode(v)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// seqToMap converts the given sequence node into a mapping node
//
// The provided value defined by field k is used as the map key.
// The used value is removed from the map value.
// If the field is not found, an error is returned.
func seqToMap(ss *yaml.Node, k string) (*yaml.Node, error) {
	if ss.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("expecting sequence node (%s provided)", ss.Tag)
	}
	if k == "" {
		return nil, fmt.Errorf("key field not defined")
	}

	var err error
	mm, _ := makeMap()

	for _, s := range ss.Content {
		if s.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("sequence may only contain mapping nodes (%s found)", s.Tag)
		}

		// Find the key value; remove from map value
		var kn *yaml.Node
		for i := 0; i < len(s.Content)-1; i += 2 {
			if s.Content[i].Value == k {
				kn = s.Content[i+1]
				s.Content = append(s.Content[:i], s.Content[i+2:]...)
				break
			}
		}

		if kn == nil {
			return nil, errors.New("key field not defined")
		}

		mm, err = addMap(mm, kn.Value, s)
		if err != nil {
			return nil, err
		}
	}

	return mm, nil
}

// inlineContent places mm mapping node content into the n mapping node content
func inlineContent(n *yaml.Node, mm ...*yaml.Node) (*yaml.Node, error) {
	if n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expecting mapping node (destination node)")
	}

	for _, m := range mm {
		if m.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("expecting mapping node (source node)")
		}
		n.Content = append(n.Content, m.Content...)
	}

	return n, nil
}

// makeMap creates a new mapping node based on the provided k, v items
//
// pp is a set of k, v items; where k's lie at i, and v's lie at i+1.
// non-string values (required by YAML nodes) are processed further.
// eg.: ["k1", "v1", "k2", "v2"]
func makeMap(pp ...interface{}) (*yaml.Node, error) {
	return addMap(&yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, pp...)
}

// addMap adds a new item to the provided mapping node
//
// pp is a set of k, v items; where k's lie at i, and v's lie at i+1.
// non-string values (required by YAML nodes) are processed further.
// eg.: ["k1", "v1", "k2", "v2"]
func addMap(n *yaml.Node, pp ...interface{}) (*yaml.Node, error) {
	if len(pp) == 0 {
		return n, nil
	}

	if len(pp)%2 == 1 {
		return nil, fmt.Errorf("uneven number of elements provided (%d): %v", len(pp), pp)
	}

	var err error

	for i := 0; i < len(pp); i += 2 {
		kRaw := pp[i]
		vRaw := pp[i+1]

		k, ok := kRaw.(string)
		if !ok {
			return nil, fmt.Errorf("keys must be of type string: %v", kRaw)
		}

		var vn *yaml.Node
		switch v := vRaw.(type) {
		case bool:
			if !v {
				continue
			}
			vn, err = encodeNode(v)

		case string:
			if v == "" {
				continue
			}
			vn, err = encodeNode(v)

			if v == "\n" {
				vn.Style = yaml.DoubleQuotedStyle
			}

		default:
			vn, err = encodeNode(vRaw)
		}

		if err != nil {
			return nil, err
		}

		// Cleanup the content
		if vn.Kind == yaml.MappingNode {
			vn = cleanMap(vn)
		}

		// Discard any null nodes and empty nodes
		if vn.Tag == "!!null" || vn.Value == "" && (len(vn.Content) == 0 || vn.Content == nil) {
			continue
		}

		n.Content = append(n.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: k},
			vn,
		)
	}
	return n, nil
}

// makeSeq creates a new sequence node based on the provided items
func makeSeq(vv ...interface{}) (*yaml.Node, error) {
	return addSeq(&yaml.Node{Kind: yaml.SequenceNode}, vv...)
}

// addSeq adds new items to the sequence node
func addSeq(n *yaml.Node, vv ...interface{}) (*yaml.Node, error) {
	var err error
	var vn *yaml.Node

	for _, vRaw := range vv {
		vn, err = encodeNode(vRaw)
		if err != nil {
			return nil, err
		}

		// Discard any null nodes and empty nodes
		if vn.Tag == "!!null" || vn.Value == "" && (len(vn.Content) == 0 || vn.Content == nil) {
			continue
		}

		n.Content = append(n.Content, vn)
	}

	return n, nil
}
