package ngimporter

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/cortezaproject/corteza-server/compose/repository"
	cct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/ngimporter/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
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
func Import(ctx context.Context, iss []types.ImportSource, ns *cct.Namespace, cfg *types.Config) error {
	// contains warnings raised by the pre process steps
	var preProcW []string
	imp := &Importer{}
	db := repository.DB(ctx)
	modRepo := repository.Module(ctx, db)
	recRepo := repository.Record(ctx, db)
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
	uMap := make(map[string]uint64)
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
			mod, err := findModuleByHandle(modRepo, ns.ID, nIs.Name)
			if mod != nil {
				types.ModulesGlobal = append(types.ModulesGlobal, mod)
			}
			if err != nil {
				preProcW = append(preProcW, err.Error()+" "+nIs.Name)
				continue
			}

			mod, err = assureLegacyFields(modRepo, mod, cfg)
			if err != nil {
				// this is a fatal error, we shouldn't continue if this fails
				return err
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

					mm, err := findModuleByID(modRepo, ns.ID, vv)
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
					if mm != nil {
						types.ModulesGlobal = append(types.ModulesGlobal, mm)
					}

					nn = imp.AddNode(nn)
					n.LinkAdd(nn)
				}
			}
		}
	}

	log.Println("PRE-PROC WARNINGS")
	for _, w := range preProcW {
		log.Printf("[warning] %s\n", w)
	}

	if cfg.RefFixup {
		err = imp.AssureLegacyID(ctx, cfg)
		if err != nil {
			log.Println("[importer] failed")
			return err
		}
	} else {
		// populate with existing users
		uMod, err := findModuleByHandle(modRepo, ns.ID, "user")
		if err != nil {
			return err
		}
		rr, _, err := recRepo.Find(uMod, cct.RecordFilter{
			ModuleID:    uMod.ID,
			Deleted:     rh.FilterStateInclusive,
			NamespaceID: ns.ID,
			Query:       "sys_legacy_ref_id IS NOT NULL",
			PageFilter: rh.PageFilter{
				Page:    1,
				PerPage: 0,
			},
		})
		if err != nil {
			return err
		}

		rvs, err := recRepo.LoadValues(uMod.Fields.Names(), rr.IDs())
		if err != nil {
			return err
		}

		err = rr.Walk(func(r *cct.Record) error {
			r.Values = rvs.FilterByRecordID(r.ID)
			return nil
		})
		if err != nil {
			return err
		}

		rr.Walk(func(r *cct.Record) error {
			vr := r.Values.Get("sys_legacy_ref_id", 0)
			vu := r.Values.Get("OwnerID", 0)
			u, err := strconv.ParseUint(vu.Value, 10, 64)
			if err != nil {
				return err
			}
			uMap[vr.Value] = u
			return nil
		})

		imp.RemoveCycles()

		// take note of leaf nodes that can be imported right away
		for _, n := range imp.nodes {
			if !n.HasChildren() {
				imp.Leafs = append(imp.Leafs, n)
			}
		}

		log.Printf("[importer] prepared\n")
		log.Printf("[importer] node count: %d\n", len(imp.nodes))
		log.Printf("[importer] leaf count: %d\n", len(imp.Leafs))

		log.Println("[importer] started")
		err = imp.Import(ctx, uMap)
		if err != nil {
			log.Println("[importer] failed")
			return err
		}
		log.Println("[importer] finished")
	}

	return nil
}

// AddNode attempts to add the given node into the graph. If the node can already be
// identified, the two nodes are merged.
func (imp *Importer) AddNode(n *types.ImportNode) *types.ImportNode {
	var fn *types.ImportNode
	for _, nn := range imp.nodes {
		if nn.CompareTo(n) {
			fn = nn
			break
		}
	}

	if fn == nil {
		imp.nodes = append(imp.nodes, n)
		return n
	}

	fn.Merge(n)
	return fn
}

// RemoveCycles removes all cycles in the given graph, by restructuring/recreating the nodes
// and their dependencies.
func (imp *Importer) RemoveCycles() {
	splice := func(n *types.ImportNode, from *types.ImportNode) {
		spl := n.Splice(from)
		imp.AddNode(spl)
	}

	for _, n := range imp.nodes {
		if !n.Visited {
			n.SeekCycles(splice)
		}
	}
}

func (m *Importer) AssureLegacyID(ctx context.Context, cfg *types.Config) error {
	db := repository.DB(ctx)
	repoRecord := repository.Record(ctx, db)
	bar := progressbar.New(len(m.nodes))

	return db.Transaction(func() (err error) {
		// since this is a ott ment to be ran after the data is already there, there is no
		// need to worry about references.
		for _, n := range m.nodes {
			ts := ""
			if cfg != nil {
				ts = cfg.ToTimestamp
			}
			err := n.AssureLegacyID(repoRecord, ts)
			if err != nil {
				return err
			}
			bar.Add(1)
		}
		return nil
	})
}

// Import runs the import over each ImportNode in the given graph
func (m *Importer) Import(ctx context.Context, users map[string]uint64) error {
	db := repository.DB(ctx)
	repoRecord := repository.Record(ctx, db)
	bar := progressbar.New(len(m.nodes))

	return db.Transaction(func() (err error) {
		for len(m.Leafs) > 0 {

			ch := make(chan types.PostProc, len(m.Leafs))
			for _, n := range m.Leafs {
				n.Import(repoRecord.With(ctx, db), users, ch, bar)
			}

			var nl []*types.ImportNode
			for len(ch) > 0 {
				pp := <-ch
				if pp.Err != nil {
					log.Printf("[importer] node %s failed with %s\n", pp.Err.Error(), pp.Node.Stringify())
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

		return nil
	})
}

func findModuleByID(repo repository.ModuleRepository, namespaceID, moduleID uint64) (*cct.Module, error) {
	var err error
	mod, err := repo.FindByID(namespaceID, moduleID)
	if err != nil {
		return nil, err
	}
	mod.Fields, err = repo.FindFields(mod.ID)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func findModuleByHandle(repo repository.ModuleRepository, namespaceID uint64, handle string) (*cct.Module, error) {
	var err error
	mod, err := repo.FindByHandle(namespaceID, handle)
	if err != nil {
		return nil, err
	}
	mod.Fields, err = repo.FindFields(mod.ID)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func assureLegacyFields(repo repository.ModuleRepository, mod *cct.Module, cfg *types.Config) (*cct.Module, error) {
	dirty := false

	// make a copy of the original fields, so we don't mess with it
	ff := make(cct.ModuleFieldSet, 0)
	mod.Fields.Walk(func(f *cct.ModuleField) error {
		ff = append(ff, f)
		return nil
	})

	// assure the legacy id reference
	f := mod.Fields.FindByName(types.LegacyRefIDField)
	if f == nil {
		dirty = true
		ff = append(ff, &cct.ModuleField{
			ModuleID: mod.ID,
			Kind:     "String",
			Name:     types.LegacyRefIDField,
		})
	}

	if dirty {
		// we are simply adding the given field, there is no harm in skipping the records checking
		err := repo.UpdateFields(mod.ID, ff, true)
		if err != nil {
			return nil, err
		}
		mod.Fields = ff
	}

	return mod, nil
}
