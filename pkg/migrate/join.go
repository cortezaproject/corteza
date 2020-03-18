package migrate

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/migrate/types"
)

type (
	// temporary node structure
	node struct {
		mg     *types.Migrateable
		Joined []*types.JoinedNode
	}
)

// Creates JoinNodes for each Migrateable node included in a source join process
// See readme for more
func sourceJoin(mm []types.Migrateable) ([]types.Migrateable, error) {
	var rr []*node
	joinedNodes := make(map[string]*types.JoinedNode)

	for _, mg := range mm {
		ww := mg
		node := &node{mg: &ww}
		rr = append(rr, node)

		if mg.Join == nil {
			continue
		}

		// join definition map defines how two sources are joined
		var joinDef map[string]string
		src, _ := ioutil.ReadAll(mg.Join)
		err := json.Unmarshal(src, &joinDef)
		if err != nil {
			return nil, err
		}

		// find all joined nodes for the given migration node
		for baseField, condition := range joinDef {
			pts := strings.Split(condition, ".")
			joinedModule := pts[0]
			joinedByField := pts[1]

			for _, m := range mm {
				if m.Name == joinedModule {
					if _, ok := joinedNodes[joinedModule]; !ok {
						ww := m
						joinedNodes[joinedModule] = &types.JoinedNode{
							Mg:        &ww,
							Name:      ww.Name,
							BaseField: baseField,
							JoinField: joinedByField,
						}
					}

					jn := joinedNodes[joinedModule]
					node.Joined = append(node.Joined, jn)
					break
				}
			}
		}
	}

	// load each joined node's entries
	for _, jn := range joinedNodes {
		jn.Entries = make(map[string][]*types.JoinEntry)
		reader := csv.NewReader(jn.Mg.Source)

		// header
		header, err := reader.Read()
		if err == io.EOF {
			break
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}

			ee := make(types.JoinEntry)
			var eId string
			for i, c := range record {
				ee[header[i]] = c

				if header[i] == jn.JoinField {
					eId = c
				}
			}

			if jn.Entries[eId] == nil {
				jn.Entries[eId] = make([]*types.JoinEntry, 0)
			}
			jn.Entries[eId] = append(jn.Entries[eId], &ee)
		}
	}

	// construct new migration nodes
	// they should include context for the join operation
	out := make([]types.Migrateable, 0)
	for _, r := range rr {
		// include only base migration nodes; exclude joined nodes
		if _, ok := joinedNodes[r.mg.Name]; !ok {
			jn := r.Joined
			for _, s := range jn {
				// no need for it further
				s.Mg = nil
			}
			mgg := types.Migrateable{
				Name:   r.mg.Name,
				Header: r.mg.Header,
				Path:   r.mg.Path,
				Source: r.mg.Source,
				Map:    r.mg.Map,
				Joins:  jn,
			}
			out = append(out, mgg)
		}
	}

	return out, nil
}
