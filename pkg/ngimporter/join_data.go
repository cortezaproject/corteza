package ngimporter

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/ngimporter/types"
)

type (
	// mapLink helps us keep track between base nodes, joined nodes and the fields used
	// for creating the link
	mapLink struct {
		jn *types.JoinNode
		// field from the base node used in the op.
		baseFields []string
		// alias to use for the base field; allows us to use the same field multiple times
		baseFieldAlias string
		// field from the joined node to use in the opp.
		joinFields []string
	}

	// temporary node for the join op.
	node struct {
		is *types.ImportSource
		// maps { alias: mapLink }
		mapper map[string]mapLink
		// maps { alisa: [field] }
		aliasMap map[string][]string
	}

	// defines params that should be used in the given join opp.
	joinEval struct {
		baseFields     []string
		baseFieldAlias string
		joinModule     string
		joinFields     []string
	}
)

// it defines join context for each ImportSource that defines the join operation.
// it returns a new set of ImportSource nodes, excluding the joined ones.
// algorighem outline:
//   * determine all nodes that define the join operation (base node)
//   * take note of all nodes, that will be used in join operations (join node)
//   * load data for each join node
//     * index each row based on the specified alias and it's specified ID field
func joinData(iss []types.ImportSource) ([]types.ImportSource, error) {
	var rr []*node
	joinedNodes := make(map[string]*types.JoinNode)

	// determine base & join nodes
	for _, mg := range iss {
		// this helps us avoid pesky pointer issues :)
		ww := mg
		nd := &node{is: &ww}
		rr = append(rr, nd)

		if mg.SourceJoin == nil {
			continue
		}

		// defer this initialization, so we can do simple nil checks
		nd.mapper = make(map[string]mapLink)
		nd.aliasMap = make(map[string][]string)

		// joinDef describes how nodes should be joined
		var joinDef map[string]string
		src, _ := ioutil.ReadAll(mg.SourceJoin)
		err := json.Unmarshal(src, &joinDef)
		if err != nil {
			return nil, err
		}

		// determine join nodes for the given base node
		for base, condition := range joinDef {
			expr := splitExpr(base, condition)

			if _, ok := nd.aliasMap[expr.baseFieldAlias]; ok {
				return nil, errors.New("alias.duplicated " + nd.is.Name + " " + expr.baseFieldAlias)
			}
			nd.aliasMap[expr.baseFieldAlias] = expr.baseFields

			for _, is := range iss {
				if is.Name == expr.joinModule {
					if _, ok := joinedNodes[expr.joinModule]; !ok {
						ww := is
						joinedNodes[expr.joinModule] = &types.JoinNode{
							Mg:   &ww,
							Name: ww.Name,
						}
					}

					// create a link between the base and joined node
					jn := joinedNodes[expr.joinModule]
					nd.mapper[expr.baseFieldAlias] = mapLink{
						jn:             jn,
						baseFields:     expr.baseFields,
						baseFieldAlias: expr.baseFieldAlias,
						joinFields:     expr.joinFields,
					}
					break
				}
			}
		}
	}

	// load join node's data
	for _, jn := range joinedNodes {
		jn.Records = make([]map[string]string, 0)
		reader := csv.NewReader(jn.Mg.Source)

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
			jn.Records = append(jn.Records, row)
			for i, c := range record {
				row[header[i]] = c
			}
		}
	}

	// generate new SourceNodes & build join node indexes
	out := make([]types.ImportSource, 0)
	for _, nd := range rr {
		found := false
		for _, jn := range joinedNodes {
			if nd.is.Name == jn.Name {
				found = true
				break
			}
		}

		// skip joined nodes
		if found {
			continue
		}

		oIs := nd.is
		// skip nodes with no mappings
		if nd.mapper == nil {
			out = append(out, *oIs)
			continue
		}

		oIs.AliasMap = nd.aliasMap

		// create `alias.ID`: []entry mappings, that will be used when importing
		for alias, link := range nd.mapper {
			if oIs.FieldMap == nil {
				oIs.FieldMap = make(map[string]types.JoinNodeRecords)
			}

			for _, e := range link.jn.Records {
				jj := []string{}
				for _, jf := range link.joinFields {
					jj = append(jj, e[jf])
				}

				index := alias + "." + strings.Join(jj[:], ".")
				if _, ok := oIs.FieldMap[index]; !ok {
					oIs.FieldMap[index] = make(types.JoinNodeRecords, 0)
				}

				oIs.FieldMap[index] = append(oIs.FieldMap[index], e)
			}
		}

		out = append(out, *oIs)
	}

	return out, nil
}

// helper to split the join expression
func splitExpr(base, joined string) joinEval {
	rr := joinEval{}

	// original node
	rx := regexp.MustCompile(`\[?(?P<bf>[\w,]+)\]?->(?P<bfa>\w+)`)
	mx := rx.FindStringSubmatch(base)
	rr.baseFields = strings.Split(mx[1], ",")
	rr.baseFieldAlias = mx[2]

	// join node
	rx = regexp.MustCompile(`(?P<jm>\w+)\.\[?(?P<jmf>[\w,]+)\]?`)
	mx = rx.FindStringSubmatch(joined)
	rr.joinModule = mx[1]
	rr.joinFields = strings.Split(mx[2], ",")

	return rr
}
