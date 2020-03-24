package migrate

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/migrate/types"
)

type (
	// mapLink helps us keep track between base nodes, joined nodes and the fields used
	// for creating the link
	mapLink struct {
		jn *types.JoinedNode
		// field from the base node used in the op.
		baseField []string
		// alias to use for the base field; allows us to use the same field multiple times
		baseFieldAlias string
		// field from the joined node to use in the opp.
		joinField []string
	}

	// temporary node for the join op.
	node struct {
		mg *types.Migrateable
		// temporary migration node mapper based on aliases
		mapper   map[string]mapLink
		aliasMap map[string][]string
	}

	exprEval struct {
		baseFields     []string
		baseFieldAlias string
		joinModule     string
		joinFields     []string
	}
)

// Creates JoinNodes for each Migrateable node included in a source join process
// See readme for more
func sourceJoin(mm []types.Migrateable) ([]types.Migrateable, error) {
	var rr []*node
	joinedNodes := make(map[string]*types.JoinedNode)

	// Algorithm outline:
	// 1. determine all migration nodes that will be used as joined nodes
	// 2. load entries for each join node
	// 3. construct new output migration nodes

	// 1. determination
	for _, mg := range mm {
		// this helps us avoid pesky pointer issues :)
		ww := mg
		nd := &node{mg: &ww}
		rr = append(rr, nd)

		if mg.Join == nil {
			continue
		}

		// defer this, so we can do simple nil checks
		nd.mapper = make(map[string]mapLink)
		nd.aliasMap = make(map[string][]string)

		// join definition map defines how two sources are joined
		var joinDef map[string]string
		src, _ := ioutil.ReadAll(mg.Join)
		err := json.Unmarshal(src, &joinDef)
		if err != nil {
			return nil, err
		}

		// find all joined nodes for the given base node
		for base, condition := range joinDef {
			expr := splitExpr(base, condition)

			if _, ok := nd.aliasMap[expr.baseFieldAlias]; ok {
				return nil, errors.New("alias.used " + nd.mg.Name + " " + expr.baseFieldAlias)
			}
			nd.aliasMap[expr.baseFieldAlias] = expr.baseFields

			// register migration node as join node
			for _, m := range mm {
				if m.Name == expr.joinModule {
					if _, ok := joinedNodes[expr.joinModule]; !ok {
						ww := m
						joinedNodes[expr.joinModule] = &types.JoinedNode{
							Mg:   &ww,
							Name: ww.Name,
						}
					}

					// create a link between the base and joined node
					jn := joinedNodes[expr.joinModule]
					nd.mapper[expr.baseFieldAlias] = mapLink{
						jn:             jn,
						baseField:      expr.baseFields,
						baseFieldAlias: expr.baseFieldAlias,
						joinField:      expr.joinFields,
					}
					break
				}
			}
		}
	}

	// 2. load entries
	for _, jn := range joinedNodes {
		jn.Entries = make([]map[string]string, 0)
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

			row := make(map[string]string)
			jn.Entries = append(jn.Entries, row)
			for i, c := range record {
				row[header[i]] = c
			}
		}
	}

	// 3. output
	out := make([]types.Migrateable, 0)
	for _, nd := range rr {
		fnd := false
		for _, jn := range joinedNodes {
			if nd.mg.Name == jn.Name {
				fnd = true
				break
			}
		}

		// skip joined nodes
		if fnd {
			continue
		}

		o := nd.mg
		// skip nodes with no mappings
		if nd.mapper == nil {
			out = append(out, *o)
			continue
		}

		o.AliasMap = nd.aliasMap

		// create `alias.ID`: []entry mappings, that will be used when migrating
		for alias, link := range nd.mapper {
			if o.FieldMap == nil {
				o.FieldMap = make(map[string]types.JoinedNodeRecords)
			}

			for _, e := range link.jn.Entries {
				jj := []string{}
				for _, jf := range link.joinField {
					jj = append(jj, e[jf])
				}

				kk := alias + "." + strings.Join(jj[:], ".")
				if _, ok := o.FieldMap[kk]; !ok {
					o.FieldMap[kk] = make(types.JoinedNodeRecords, 0)
				}

				o.FieldMap[kk] = append(o.FieldMap[kk], e)
			}
		}

		out = append(out, *o)
	}

	return out, nil
}

// helper to split the join expression
func splitExpr(base, joined string) exprEval {
	rr := exprEval{}

	// original node
	rx := regexp.MustCompile(`\[?(?P<bf>[\w,]+)\]?->(?P<bfa>\w+)`)
	mx := rx.FindStringSubmatch(base)
	rr.baseFields = strings.Split(mx[1], ",")
	rr.baseFieldAlias = mx[2]

	// joined node
	rx = regexp.MustCompile(`(?P<jm>\w+)\.\[?(?P<jmf>[\w,]+)\]?`)
	mx = rx.FindStringSubmatch(joined)
	rr.joinModule = mx[1]
	rr.joinFields = strings.Split(mx[2], ",")

	return rr
}
