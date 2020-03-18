package types

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/schollz/progressbar/v2"
)

type (
	// graph node
	Node struct {
		// unique node name
		Name string

		// keep note of parent refs, so we don't need to inverse it ;)
		Parents []*Node
		// keep note of Children, as they are our dependencies
		Children []*Node
		// mapping from migrated IDs to Corteza IDs
		mapping map[string]Map
		// determines if node is in current path; used for loop detection
		inPath bool
		// determines if node was spliced; used to break the loop
		spliced  bool
		original *Node
		spl      *Node

		// records are applicable in the case of spliced nodes
		records []*types.Record

		// some refs
		Module    *types.Module
		Namespace *types.Namespace
		Reader    *csv.Reader

		// meta
		Header  []string
		Visited bool

		Lock *sync.Mutex
	}

	// map between migrated ID and Corteza ID
	Map map[string]string

	PostProc struct {
		Leafs []*Node
		Err   error
	}
)

// helper, to determine if the two nodes are equal
func (n *Node) Compare(to *Node) bool {
	return n.Name == to.Name && n.spliced == to.spliced
}

// helper to stringify the node
func (n *Node) Stringify() string {
	return fmt.Sprintf("NODE > n: %s; spliced: %t; inPath: %t;", n.Name, n.spliced, n.inPath)
}

// adds a new map to the given node
func (n *Node) addMap(key string, m Map) {
	n.Lock.Lock()
	if n.mapping == nil {
		n.mapping = map[string]Map{}
	}

	n.mapping[key] = m
	n.Lock.Unlock()
}

// does the actual data migration for the given node
func (n *Node) Migrate(repoRecord repository.RecordRepository, users map[string]uint64, wg *sync.WaitGroup, ch chan PostProc, bar *progressbar.ProgressBar) {
	defer wg.Done()
	defer bar.Add(1)

	var err error

	mapping := make(Map)
	if n.Reader != nil {
		// if records exist (from spliced node); correct refs
		if !n.spliced && n.records != nil && len(n.records) > 0 {
			// we can just reuse the mapping object, since it will remain the same
			mapping = n.mapping[fmt.Sprint(n.Module.ID)]

			err := updateRefs(n, repoRecord)
			if err != nil {
				ch <- PostProc{
					Leafs: nil,
					Err:   err,
				}
				return
			}
		} else {
			mapping, err = importNodeSource(n, users, repoRecord)
			if err != nil {
				ch <- PostProc{
					Leafs: nil,
					Err:   err,
				}
				return
			}
		}
	}

	var rtr []*Node

	var pps []*Node
	for _, pp := range n.Parents {
		pps = append(pps, pp)
	}

	// update node refs
	for _, p := range pps {
		rtr = append(rtr, p)

		// pass mapping object to the node's parend so it can migrate it's data
		p.addMap(fmt.Sprint(n.Module.ID), mapping)
		p.LinkRemove(n)
	}

	ch <- PostProc{
		Leafs: rtr,
		Err:   nil,
	}
}

// determines if node is Satisfied and can be imported
// it is Satisfied, when all of it's dependencies have been imported ie. no
// more child refs
func (n *Node) Satisfied() bool {
	return !n.HasChildren()
}

func (n *Node) HasChildren() bool {
	return n.Children != nil && len(n.Children) > 0
}

// partially Merge the two nodes
func (n *Node) Merge(nn *Node) {
	if nn.Module != nil {
		n.Module = nn.Module
	}
	if nn.Reader != nil {
		n.Reader = nn.Reader
	}
	if nn.Header != nil {
		n.Header = nn.Header
	}
}

// link the two nodes
func (n *Node) LinkAdd(to *Node) {
	n.addChild(to)
	to.addParent(n)
}

// remove the link between the two nodes
func (n *Node) LinkRemove(from *Node) {
	n.Lock.Lock()
	n.Children = n.removeIfPresent(from, n.Children)
	from.Parents = from.removeIfPresent(n, from.Parents)
	n.Lock.Unlock()
}

// adds a parent node to the given node
func (n *Node) addParent(add *Node) {
	n.Parents = n.addIfMissing(add, n.Parents)
}

// adds a child node to the given node
func (n *Node) addChild(add *Node) {
	n.Children = n.addIfMissing(add, n.Children)
}

// adds a node, if it doesn't yet exist
func (n *Node) addIfMissing(add *Node, list []*Node) []*Node {
	var fn *Node

	for _, nn := range list {
		if add.Compare(nn) {
			fn = nn
		}
	}

	if fn != nil {
		fn.Merge(add)
		return list
	}
	return append(list, add)
}

// removes the node, if it exists
func (n *Node) removeIfPresent(rem *Node, list []*Node) []*Node {
	for i, nn := range list {
		if rem.Compare(nn) {
			// https://stackoverflow.com/a/37335777
			list[len(list)-1], list[i] = list[i], list[len(list)-1]
			return list[:len(list)-1]
		}
	}

	return list
}

// traverses the graph and notifies us of any cycles
func (n *Node) Traverse(cycle func(n *Node, to *Node)) {
	n.inPath = true
	n.Visited = true

	var cc []*Node
	for _, nn := range n.Children {
		cc = append(cc, nn)
	}

	for _, nn := range cc {
		if n.Name == "client" {
		}

		if nn.inPath {
			cycle(n, nn)
		} else {
			nn.Traverse(cycle)
		}
	}

	n.inPath = false
}

func (n *Node) DFS() {
	n.inPath = true

	for _, nn := range n.Children {
		if !nn.inPath {
			nn.DFS()
		}
	}

	n.inPath = false
}

// clones the given node
func (n *Node) clone() *Node {
	return &Node{
		Name:      n.Name,
		Parents:   n.Parents,
		Children:  n.Children,
		mapping:   n.mapping,
		inPath:    n.inPath,
		spliced:   n.spliced,
		original:  n.original,
		records:   n.records,
		Visited:   n.Visited,
		Module:    n.Module,
		Namespace: n.Namespace,
		Reader:    n.Reader,
		Header:    n.Header,
	}
}

// splices the node from the original graph and removes the cycle
func (n *Node) Splice(from *Node) *Node {
	splicedN := from.spl

	if splicedN == nil {
		splicedN = from.clone()
		splicedN.spliced = true
		splicedN.Parents = nil
		splicedN.Children = nil
		splicedN.inPath = false

		splicedN.original = from
		from.spl = splicedN

		from.LinkAdd(splicedN)
	}

	n.LinkRemove(from)
	n.LinkAdd(splicedN)

	return splicedN
}

func sysField(f string) bool {
	switch f {
	case "OwnerId",
		"IsDeleted",
		"CreatedDate",
		"CreatedById",
		"LastModifiedDate",
		"LastModifiedById":
		return true
	}
	return false
}

func updateRefs(n *Node, repo repository.RecordRepository) error {
	// correct references
	for _, r := range n.records {
		for _, v := range r.Values {
			var f *types.ModuleField
			// find the applicable module field
			for _, ff := range n.Module.Fields {
				if ff.Name == v.Name {
					f = ff
					break
				}
			}

			if f == nil {
				continue
			}

			val := v.Value
			// determine value based on the provided map
			if f.Options["moduleID"] != nil {
				ref, ok := f.Options["moduleID"].(string)
				if !ok {
					return errors.New("moduleField.record.invalidRefFormat")
				}

				if mod, ok := n.mapping[ref]; ok {
					if vv, ok := mod[val]; ok {
						v.Value = vv
					} else {
						v.Value = ""
						continue
					}
				} else {
					v.Value = ""
					continue
				}
			}
		}

		// update values
		nv := types.RecordValueSet{}
		for _, v := range r.Values {
			if v.Value != "" {
				nv = append(nv, v)
			}
		}

		r.Values = nv
		err := repo.UpdateValues(r.ID, r.Values)
		if err != nil {
			return err
		}
	}
	return nil
}

func importNodeSource(n *Node, users map[string]uint64, repo repository.RecordRepository) (Map, error) {
	mapping := make(Map)

	fixUtf := func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}

	for {
	looper:
		record, err := n.Reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		rr := &types.Record{
			ModuleID:    n.Module.ID,
			NamespaceID: n.Namespace.ID,
			CreatedAt:   time.Now(),
		}

		vals := types.RecordValueSet{}
		for i, h := range n.Header {
			val := record[i]

			if sysField(h) {
				switch h {
				case "OwnerId":
					rr.OwnedBy = users[val]
					break

					// ignore deleted values, as SF provides minimal info about those
				case "IsDeleted":
					if val == "1" {
						goto looper
					}
					break

				case "CreatedDate":
					if val != "" {
						rr.CreatedAt, err = time.Parse(SfDateTime, val)
						if err != nil {
							return nil, err
						}
					}
					break

				case "CreatedById":
					rr.CreatedBy = users[val]
					break

				case "LastModifiedById":
					rr.UpdatedBy = users[val]
					break

				case "LastModifiedDate":
					if val != "" {
						tt, err := time.Parse(SfDateTime, val)
						rr.UpdatedAt = &tt
						if err != nil {
							return nil, err
						}
					}
					break
				}
			} else {
				var f *types.ModuleField
				for _, ff := range n.Module.Fields {
					if ff.Name == h {
						f = ff
						break
					}
				}

				if f == nil {
					continue
				}

				if f.Options["moduleID"] != nil {
					// spliced nodes should NOT manage their references
					if !n.spliced {
						ref, ok := f.Options["moduleID"].(string)
						if !ok {
							return nil, errors.New("moduleField.record.invalidRefFormat")
						}

						if mod, ok := n.mapping[ref]; ok {
							if v, ok := mod[val]; ok {
								val = v
							} else {
								continue
							}
						} else {
							continue
						}
					}
				} else if f.Kind == "User" {
					if u, ok := users[val]; ok {
						val = fmt.Sprint(u)
					} else {
						continue
					}
				} else {
					val = strings.Map(fixUtf, val)

					if val == "" {
						continue
					}
				}

				vals = append(vals, &types.RecordValue{
					Name:  h,
					Value: val,
				})
			}
		}

		// create record
		r, err := repo.Create(rr)
		if err != nil {
			return nil, err
		}

		// update record values with recordID
		for _, v := range vals {
			v.RecordID = r.ID
		}
		err = repo.UpdateValues(r.ID, vals)
		if err != nil {
			return nil, err
		}

		// spliced nodes should preserve their records for later ref processing
		if n.spliced {
			rr.Values = vals
			n.original.records = append(n.original.records, rr)
		}

		mapping[record[0]] = fmt.Sprint(rr.ID)
	}

	return mapping, nil
}
