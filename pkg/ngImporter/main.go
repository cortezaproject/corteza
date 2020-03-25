package migrate

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	cct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/ngImporter/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/schollz/progressbar/v2"
)

type (
	// Importer contains the context of the entire importing operation
	Importer struct {
		// a set of import nodes that define a graph
		nodes []*types.ImportNode
		// a set of leaf import nodes, that can be imported in the next cycle
		Leafs []*types.ImportNode
	}
)

// Import initializes the import process for the given set of ImportSource nodes
// algorithm outline:
//   * import all users used within the given import sources
//   * handle source join operations
//   * handle data mapping operations
//   * build graph from ImportNodes based on the provided ImportSource nodes
//   * remove cycles from the given graph
//   * import data based on node dependencies
func Import(iss []types.ImportSource, ns *cct.Namespace, ctx context.Context) error {
	// contains warnings raised by the pre process steps
	var preProcW []string
	imp := &Importer{}
	svcMod := service.DefaultModule.With(ctx)
	var err error

	// import users
	var usrSrc *types.ImportSource
	for _, m := range iss {
		if m.Name == types.UserModHandle {
			usrSrc = &m
			break
		}
	}

	// maps sourceUserID to CortezaID
	var uMap map[string]uint64
	if usrSrc != nil {
		um, mgu, err := importUsers(ctx, usrSrc, ns)
		if err != nil {
			return err
		}
		uMap = um

		// replace the old source node with the new one (updated data stream)
		found := false
		for i, m := range iss {
			if m.Name == mgu.Name {
				iss[i] = *mgu
				found = true
				break
			}
		}

		if !found {
			iss = append(iss, *mgu)
		}
	}

	iss, err = joinData(iss)
	if err != nil {
		return err
	}

	// data mapping & graph construction
	for _, is := range iss {
		nIss, err := mapData(is)
		if err != nil {
			return err
		}

		for _, nIs := range nIss {
			// preload module
			mod, err := svcMod.FindByHandle(ns.ID, nIs.Name)
			if err != nil {
				preProcW = append(preProcW, err.Error()+" "+nIs.Name)
				continue
			}

			// define headers
			r := csv.NewReader(nIs.Source)
			var header []string
			if nIs.Header != nil {
				header = *nIs.Header
			} else {
				header, err = r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}

			// create node & push to graph
			n := &types.ImportNode{
				Name:      nIs.Name,
				Module:    mod,
				Namespace: ns,
				Reader:    r,
				Header:    header,
				Lock:      &sync.Mutex{},
				FieldMap:  nIs.FieldMap,
				ValueMap:  nIs.ValueMap,
			}
			n = imp.AddNode(n)

			// prepare additional import nodes based on it's record fields
			for _, f := range mod.Fields {
				if f.Kind == "Record" {
					refMod := f.Options["moduleID"]
					if refMod == nil {
						preProcW = append(preProcW, "moduleField.record.missingRef"+" "+nIs.Name+" "+f.Name)
						continue
					}

					modID, ok := refMod.(string)
					if !ok {
						preProcW = append(preProcW, "moduleField.record.invalidRefFormat"+" "+nIs.Name+" "+f.Name)
						continue
					}

					vv, err := strconv.ParseUint(modID, 10, 64)
					if err != nil {
						preProcW = append(preProcW, err.Error())
						continue
					}

					mm, err := svcMod.FindByID(ns.ID, vv)
					if err != nil {
						preProcW = append(preProcW, err.Error()+" "+nIs.Name+" "+f.Name+" "+modID)
						continue
					}

					nn := &types.ImportNode{
						Name:      mm.Handle,
						Module:    mm,
						Namespace: ns,
						Lock:      &sync.Mutex{},
					}

					nn = imp.AddNode(nn)
					n.LinkAdd(nn)
				}
			}
		}
	}

	spew.Dump("PRE-PROC WARNINGS", preProcW)

	imp.RemoveCycles()

	// take note of leaf nodes that can be imported right away
	for _, n := range imp.nodes {
		if !n.HasChildren() {
			imp.Leafs = append(imp.Leafs, n)
		}
	}

	fmt.Printf("import.prepared\n")
	fmt.Printf("no. of nodes %d\n", len(imp.nodes))
	fmt.Printf("no. of entry points %d\n", len(imp.Leafs))

	err = imp.Import(ctx, uMap)
	if err != nil {
		return err
	}

	return nil
}

// AddNode attempts to add the given node into the graph. If the node can already be
// identified, the two nodes are merged.
func (m *Importer) AddNode(n *types.ImportNode) *types.ImportNode {
	var fn *types.ImportNode
	for _, nn := range m.nodes {
		if nn.CompareTo(n) {
			fn = nn
			break
		}
	}

	if fn == nil {
		m.nodes = append(m.nodes, n)
		return n
	}

	fn.Merge(n)
	return fn
}

// RemoveCycles removes all cycles in the given graph, by restructuring/recreating the nodes
// and their dependencies.
func (m *Importer) RemoveCycles() {
	splice := func(n *types.ImportNode, from *types.ImportNode) {
		spl := n.Splice(from)
		m.AddNode(spl)
	}

	for _, n := range m.nodes {
		if !n.Visited {
			n.SeekCycles(splice)
		}
	}
}

// Import runs the import over each ImportNode in the given graph
func (m *Importer) Import(ctx context.Context, users map[string]uint64) error {
	fmt.Println("(•_•)")
	fmt.Println("(-•_•)>⌐■-■")
	fmt.Println("(⌐■_■)")
	fmt.Print("\n\n\n")

	db := repository.DB(ctx)
	repoRecord := repository.Record(ctx, db)
	bar := progressbar.New(len(m.nodes))

	return db.Transaction(func() (err error) {
		for len(m.Leafs) > 0 {
			var wg sync.WaitGroup

			ch := make(chan types.PostProc, len(m.Leafs))
			for _, n := range m.Leafs {
				wg.Add(1)
				go n.Import(repoRecord, users, &wg, ch, bar)
			}

			wg.Wait()

			var nl []*types.ImportNode
			for len(ch) > 0 {
				pp := <-ch
				if pp.Err != nil {
					spew.Dump("ERR", pp.Err, pp.Node.Stringify())
					return pp.Err
				}

				// update the set of available leaf nodes
				if pp.Leafs != nil {
					for _, n := range pp.Leafs {
						for _, l := range nl {
							if n.CompareTo(l) {
								goto skip
							}
						}
						if n.Satisfied() {
							nl = append(nl, n)
						}

					skip:
					}
				}
			}
			m.Leafs = nl
		}

		fmt.Print("\n\n\n")
		fmt.Println("(⌐■_■)")
		fmt.Println("(-•_•)>⌐■-■")
		fmt.Println("(•_•)")

		return nil
	})
}
