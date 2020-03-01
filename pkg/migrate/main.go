package migrate

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	cct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/migrate/types"
	sysRepo "github.com/cortezaproject/corteza-server/system/repository"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

var (
	userModHandle = "User"
)

type (
	Migrator struct {
		// a set of nodes included in the migration
		nodes []*types.Node

		// list of leaf nodes, that we might be able to migrate
		Leafs []*types.Node
	}
)

func Migrate(mg []types.Migrateable, ns *cct.Namespace, ctx context.Context) error {
	mig := &Migrator{}
	svcMod := service.DefaultModule.With(ctx)

	// 1. migrate all the users, so we can reference then accross the entire system
	var mgUsr types.Migrateable
	for _, m := range mg {
		if m.Name == userModHandle {
			mgUsr = m
			break
		}
	}

	uMap, err := migrateUsers(mgUsr, ns, ctx)
	if err != nil {
		return err
	}

	// 2. prepare and link migration nodes
	for _, mgR := range mg {
		ss, err := splitStream(mgR)
		if err != nil {
			return err
		}

		for _, m := range ss {
			fmt.Printf("mg.processing > %s\n", m.Name)

			// 2.1 load module
			mod, err := svcMod.FindByHandle(ns.ID, m.Name)
			if err != nil {
				return err
			}

			// 2.2 get header fields
			r := csv.NewReader(m.Source)
			var header []string
			if m.Header != nil {
				header = *m.Header
			} else {
				header, err = r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}

			// 2.3 create migration node
			n := &types.Node{
				Name:      m.Name,
				Module:    mod,
				Namespace: ns,
				Reader:    r,
				Header:    header,
				Lock:      &sync.Mutex{},
			}
			n = mig.AddNode(n)

			// 2.4 prepare additional migration nodes, to provide dep. constraints
			for _, f := range mod.Fields {
				if f.Kind == "Record" {
					refMod := f.Options["moduleID"]
					if refMod == nil {
						return errors.New("moduleField.record.missingRef")
					}

					modID, ok := refMod.(string)
					if !ok {
						return errors.New("moduleField.record.invalidRefFormat")
					}
					fmt.Printf("mg.node.link > %s [%s]\n", f.Name, modID)

					vv, err := strconv.ParseUint(modID, 10, 64)
					if err != nil {
						return err
					}

					mm, err := svcMod.FindByID(ns.ID, vv)
					if err != nil {
						return err
					}

					nn := &types.Node{
						Name:      mm.Handle,
						Module:    mm,
						Namespace: ns,
						Lock:      &sync.Mutex{},
					}

					nn = mig.AddNode(nn)
					n.LinkAdd(nn)
				}
			}

			fmt.Printf("mg.processed > %s\n\n\n", m.Name)
		}
	}

	fmt.Printf("graph.remove.cycles\n")
	mig.MakeAcyclic()

	for _, n := range mig.nodes {
		// keep track of leaf nodes for later importing
		if !n.HasChildren() {
			mig.Leafs = append(mig.Leafs, n)
		}
	}

	fmt.Printf("migration.prepared\n")
	fmt.Printf("no. of nodes %d\n", len(mig.nodes))
	fmt.Printf("no. of entry points %d\n", len(mig.Leafs))

	fmt.Printf("\n\nmigrator.migrating\n")
	err = mig.Migrate(ctx, uMap)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nmigrator.migrating.finished\n")

	return nil
}

// if function resolves an existing node, it will merge with the provided node
// and return the new reference
func (m *Migrator) AddNode(n *types.Node) *types.Node {
	var fn *types.Node
	for _, nn := range m.nodes {
		if nn.Compare(n) {
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

// it converts the graph from a cyclic (unsafe) graph to an acyclic (safe) graph
// that can be processed with a single algorithm
func (m *Migrator) MakeAcyclic() {
	// splices the node from the cycle and thus preventing the cycle
	splice := func(n *types.Node, from *types.Node) {
		spl := n.Splice(from)
		m.AddNode(spl)
	}

	for _, n := range m.nodes {
		if !n.Visited {
			n.Traverse(splice)
		}
	}
}

// processess migration nodes and migrates the data from the provided source files
func (m *Migrator) Migrate(ctx context.Context, users map[string]uint64) error {
	db := repository.DB(ctx)
	repoRecord := repository.Record(ctx, db)

	return db.Transaction(func() (err error) {
		for len(m.Leafs) > 0 {
			var wg sync.WaitGroup

			ch := make(chan types.PostProc, len(m.Leafs))
			for _, n := range m.Leafs {
				wg.Add(1)

				// migrate & update leaf nodes
				go n.Migrate(repoRecord, users, &wg, ch)
			}

			wg.Wait()

			var nl []*types.Node
			for len(ch) > 0 {
				pp := <-ch
				if pp.Err != nil {
					return pp.Err
				}

				if pp.Leafs != nil {
					for _, n := range pp.Leafs {
						for _, l := range nl {
							if n.Compare(l) {
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

// migrates provided users
// this should be a pre-requisite to any further migration, as user information is required
func migrateUsers(mg types.Migrateable, ns *cct.Namespace, ctx context.Context) (map[string]uint64, error) {
	db := repository.DB(ctx)
	repoUser := sysRepo.User(ctx, db)
	// this provides a map between SF ID -> CortezaID
	mapping := make(map[string]uint64)

	// get fields
	var srcBuf bytes.Buffer
	tee := io.TeeReader(mg.Source, &srcBuf)
	r := csv.NewReader(tee)
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// create users
	for {
	looper:
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		u := &sysTypes.User{}
		for i, h := range header {
			val := record[i]

			// when creating users we only care about a handfull of values.
			// the rest are included in the module
			switch h {
			case "Username":
				u.Username = record[i]
				break

			case "Email":
				u.Email = record[i]
				break

			case "FirstName":
				u.Name = record[i]
				break

			case "LastName":
				u.Name = u.Name + " " + record[i]
				break

			case "CreatedDate":
				if val != "" {
					u.CreatedAt, err = time.Parse(types.SfDateTime, val)
					if err != nil {
						return nil, err
					}
				}
				break

			case "LastModifiedDate":
				if val != "" {
					tt, err := time.Parse(types.SfDateTime, val)
					u.UpdatedAt = &tt
					if err != nil {
						return nil, err
					}
				}
				break

				// ignore deleted values, as SF provides minimal info about those
			case "IsDeleted":
				if val != "" {
					goto looper
				}
			}
		}

		// this allows us to reuse existing users
		uu, err := repoUser.FindByEmail(u.Email)
		if err == nil {
			u = uu
		} else {
			u, err = repoUser.Create(u)
			if err != nil {
				return nil, err
			}
		}

		mapping[record[0]] = u.ID
	}

	mg.Source = &srcBuf
	return mapping, nil
}
