package types

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/compose/repository"
	cv "github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/schollz/progressbar/v2"
)

type (
	// ImportNode helps us perform the actual import.
	// Multiple ImportNodes define a graph, which helps us with dependency resolution
	// and determination of proper import order
	ImportNode struct {
		// Name is used for unique node identification. It should match the target resource name.
		Name string

		// Parents represents the node's parents and the nodes that depend on this node.
		Parents []*ImportNode
		// Children represents the node's children and this node's dependencies.
		Children []*ImportNode

		// used for idMap between import source's IDs into CortezaIDs
		idMap map[string]Map

		// determines if node is in current path; used for cycle detection
		inPath bool
		// determines if this node was spliced from the original path in order to break the cycle
		isSpliced bool
		// points to the original node (from spliced)
		original *ImportNode
		// points to the spliced node (from original)
		spliced *ImportNode
		// defines the records that were created by the spliced node.
		// They are later used to insert missing dependencies.
		records []*types.Record

		// some refs...
		Module    *types.Module
		Namespace *types.Namespace
		Reader    *csv.Reader

		// some meta...
		Header  []string
		Visited bool
		Lock    *sync.Mutex

		// FieldMap stores records from the joined import source.
		// Records are indexed by {alias: [record]}
		FieldMap map[string]JoinNodeRecords

		// Value Map allows us to map specific values from the given import source into
		// a specified value used by Corteza.
		ValueMap map[string]map[string]string
	}

	// Map maps between import sourceID -> CortezaID
	Map map[string]string
)

// CompareTo compares the two nodes. It uses the name and it's variant
func (n *ImportNode) CompareTo(to *ImportNode) bool {
	return n.Name == to.Name && n.isSpliced == to.isSpliced
}

// Stringify stringifies the given node; usefull for debugging
func (n *ImportNode) Stringify() string {
	return fmt.Sprintf("NODE > n: %s; spliced: %t", n.Name, n.isSpliced)
}

// adds a new ID map to the given node's existing ID map
func (n *ImportNode) addMap(key string, m Map) {
	n.Lock.Lock()
	defer n.Lock.Unlock()

	if n.idMap == nil {
		n.idMap = map[string]Map{}
	}

	n.idMap[key] = m
}

// Import performs the actual data import.
// The algoritem defines two steps:
//   * source import,
//   * reference correction.
// For details refer to the README.
func (n *ImportNode) Import(repoRecord repository.RecordRepository, users map[string]uint64, ch chan PostProc, bar *progressbar.ProgressBar) {
	defer bar.Add(1)

	var err error

	mapping := make(Map)
	if n.Reader != nil {
		// when importing a node that defined a spliced node, we should only correct it's refs
		if !n.isSpliced && n.records != nil && len(n.records) > 0 {
			// we can just reuse the mapping object, since it will remain the same
			mapping = n.idMap[fmt.Sprint(n.Module.ID)]

			err := n.correctRecordRefs(repoRecord)
			if err != nil {
				ch <- PostProc{
					Leafs: nil,
					Err:   err,
					Node:  n,
				}
				return
			}
		} else {
			// when importing a spliced node or a node that did not define a spliced node, we should
			// import it's data
			mapping, err = n.importNodeSource(users, repoRecord)
			if err != nil {
				ch <- PostProc{
					Leafs: nil,
					Err:   err,
					Node:  n,
				}
				return
			}
		}
	}

	var rtr []*ImportNode

	var pps []*ImportNode
	for _, pp := range n.Parents {
		pps = append(pps, pp)
	}

	// update node refs
	for _, p := range pps {
		rtr = append(rtr, p)

		// pass mapping object to the node's parent so it can map handle dependency refs
		p.addMap(fmt.Sprint(n.Module.ID), mapping)
		p.LinkRemove(n)
	}

	ch <- PostProc{
		Leafs: rtr,
		Err:   nil,
	}
}

func (n *ImportNode) fetchRemoteRef(ref, refMod string, repo repository.RecordRepository) (string, error) {
	refModU, err := strconv.ParseUint(refMod, 10, 64)
	if err != nil {
		return "", err
	}

	fl := types.RecordFilter{
		ModuleID:    refModU,
		NamespaceID: n.Namespace.ID,
		Deleted:     rh.FilterStateInclusive,
		Query:       fmt.Sprintf("%s='%s'", LegacyRefIDField, ref),
		PageFilter: rh.PageFilter{
			Page:    1,
			PerPage: 1,
		},
	}

	var refModM *types.Module
	if ModulesGlobal != nil {
		refModM = ModulesGlobal.FindByID(refModU)
	}
	if refModM != nil {
		rr, _, err := repo.Find(refModM, fl)
		if err != nil {
			return "", err
		}
		if len(rr) < 1 {
			return "", errors.New(fmt.Sprintf("[error] referenced record %s not found on node %s for module %s", ref, n.Name, refModM.Name))
		}
		return strconv.FormatUint(rr[0].ID, 10), nil
	}

	return "", nil
}

func (n *ImportNode) AssureLegacyID(repoRecord repository.RecordRepository, toTimestamp string) error {
	limit := uint(10000)
	pager := func(page uint) (types.RecordSet, *types.RecordFilter, error) {
		// fetch all records, ordered by the ID for this module before the specified timestamp (if provided)
		f := types.RecordFilter{
			Sort:        "id ASC",
			Deleted:     rh.FilterStateInclusive,
			ModuleID:    n.Module.ID,
			NamespaceID: n.Namespace.ID,
			PageFilter: rh.PageFilter{
				Page:    page,
				PerPage: limit,
			},
		}
		if toTimestamp != "" {
			f.Query = fmt.Sprintf("createdAt <= '%s'", toTimestamp)
		}

		rr, ff, err := repoRecord.Find(n.Module, f)
		rvs, err := repoRecord.LoadValues(n.Module.Fields.Names(), rr.IDs())
		if err != nil {
			return nil, nil, err
		}

		err = rr.Walk(func(r *types.Record) error {
			r.Values = rvs.FilterByRecordID(r.ID)
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
		return rr, &ff, nil
	}

	// loop through the csv entries and provide the legacy ref id field value
	i := uint(0)
	page := uint(1)
	var rr types.RecordSet
	var f *types.RecordFilter
	var err error

	for {
		// <= because i is 0-based (array indexes)
		if f == nil || i >= f.Page*f.PerPage {
			rr, f, err = pager(page)
			if err != nil {
				return err
			}
			page++
		}

		// this only happenes when there is no source for the module; ie. some imported source
		// references a module that was not there initially.
		// such cases can be skipped.
		if n.Reader == nil {
			return nil
		}

		record, err := n.Reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// since the importer skips these, these should also be ignored here
		if record[0] == "" {
			continue
		}

		if i >= uint(f.Count) {
			return errors.New(fmt.Sprintf("[error] the number of csv entries exceeded record count: %d for node: %s", f.Count, n.Name))
		}
		r := rr[i-((f.Page-1)*f.PerPage)]
		rvs := r.Values
		rv := rvs.FilterByName(LegacyRefIDField)
		if rv == nil {
			rvs = append(rvs, &types.RecordValue{
				RecordID: r.ID,
				Name:     LegacyRefIDField,
				Place:    0,
				Value:    record[0],
				Updated:  true,
			})

			err := repoRecord.UpdateValues(r.ID, rvs)
			if err != nil {
				return err
			}
		}

		i++
	}

	// final sanity checks
	// - check that the counters match up
	if f.Count != i {
		return errors.New(fmt.Sprintf("[error] the number of records and csv entries don't match; records: %d, csv: %d, node: %s", f.Count, i, n.Name))
	}

	return nil
}

// determines if node is Satisfied and can be imported
// it is Satisfied, when all of it's dependencies have been imported ie. no
// more child refs
func (n *ImportNode) Satisfied() bool {
	return !n.HasChildren()
}

func (n *ImportNode) HasChildren() bool {
	return n.Children != nil && len(n.Children) > 0
}

// partially Merge the two nodes
func (n *ImportNode) Merge(nn *ImportNode) {
	if nn.Module != nil {
		n.Module = nn.Module
	}
	if nn.Reader != nil {
		n.Reader = nn.Reader
	}
	if nn.Header != nil {
		n.Header = nn.Header
	}
	if nn.FieldMap != nil {
		n.FieldMap = nn.FieldMap
	}
	if nn.ValueMap != nil {
		n.ValueMap = nn.ValueMap
	}
}

// link the two nodes
func (n *ImportNode) LinkAdd(to *ImportNode) {
	n.addChild(to)
	to.addParent(n)
}

// remove the link between the two nodes
func (n *ImportNode) LinkRemove(from *ImportNode) {
	n.Lock.Lock()
	n.Children = n.removeIfPresent(from, n.Children)
	from.Parents = from.removeIfPresent(n, from.Parents)
	n.Lock.Unlock()
}

// adds a parent node to the given node
func (n *ImportNode) addParent(add *ImportNode) {
	n.Parents = n.addIfMissing(add, n.Parents)
}

// adds a child node to the given node
func (n *ImportNode) addChild(add *ImportNode) {
	n.Children = n.addIfMissing(add, n.Children)
}

// adds a node, if it doesn't yet exist
func (n *ImportNode) addIfMissing(add *ImportNode, list []*ImportNode) []*ImportNode {
	var fn *ImportNode

	for _, nn := range list {
		if add.CompareTo(nn) {
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
func (n *ImportNode) removeIfPresent(rem *ImportNode, list []*ImportNode) []*ImportNode {
	for i, nn := range list {
		if rem.CompareTo(nn) {
			// https://stackoverflow.com/a/37335777
			list[len(list)-1], list[i] = list[i], list[len(list)-1]
			return list[:len(list)-1]
		}
	}

	return list
}

// SeekCycles finds cycles & calls the given function
func (n *ImportNode) SeekCycles(cycle func(n *ImportNode, to *ImportNode)) {
	n.inPath = true
	n.Visited = true

	var cc []*ImportNode
	for _, nn := range n.Children {
		cc = append(cc, nn)
	}

	for _, nn := range cc {
		if nn.inPath {
			cycle(n, nn)
		} else {
			nn.SeekCycles(cycle)
		}
	}

	n.inPath = false
}

// clones the given node
func (n *ImportNode) clone() *ImportNode {
	return &ImportNode{
		Name:      n.Name,
		Parents:   n.Parents,
		Children:  n.Children,
		idMap:     n.idMap,
		inPath:    n.inPath,
		isSpliced: n.isSpliced,
		original:  n.original,
		records:   n.records,
		Visited:   n.Visited,
		Module:    n.Module,
		Namespace: n.Namespace,
		Reader:    n.Reader,
		Header:    n.Header,
		FieldMap:  n.FieldMap,
		ValueMap:  n.ValueMap,
	}
}

// splices the node from the original graph and removes the cycle
func (n *ImportNode) Splice(from *ImportNode) *ImportNode {
	splicedN := from.spliced

	if splicedN == nil {
		splicedN = from.clone()
		splicedN.isSpliced = true
		splicedN.Parents = nil
		splicedN.Children = nil
		splicedN.inPath = false

		splicedN.original = from
		from.spliced = splicedN

		from.LinkAdd(splicedN)
	}

	n.LinkRemove(from)
	n.LinkAdd(splicedN)

	return splicedN
}

// helper to determine if this is a system field
func isSysField(f string) bool {
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

// updates the given node's record values that depend on another record
func (n *ImportNode) correctRecordRefs(repo repository.RecordRepository) error {
	s := cv.Sanitizer()

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

				fetch := false

				// in case of a missing ref, make sure to remove the reference.
				// otherwise this will cause internal errors when trying to resolve CortezaID.
				if mod, ok := n.idMap[ref]; !ok {
					fetch = true
				} else if vv, ok := mod[val]; !ok {
					fetch = true
				} else {
					v.Value = vv
					v.Updated = true
				}

				if fetch {
					val, err := n.fetchRemoteRef(val, ref, repo)
					if err != nil {
						continue
					}
					v.Value = val
					v.Updated = true
				}
			}
		}

		// update values; skip out empty values
		nv := types.RecordValueSet{}
		for _, v := range r.Values {
			if v.Value != "" {
				nv = append(nv, v)
			}
		}

		nv = s.Run(n.Module, nv)

		r.Values = nv
		err := repo.UpdateValues(r.ID, r.Values)
		if err != nil {
			log.Printf("[issue] db.UpdateValues | %d | %s | %s \n", r.ID, r.Values.String(), err.Error())
			// return err
		}
	}
	return nil
}

// imports the given node's source
func (n *ImportNode) importNodeSource(users map[string]uint64, repo repository.RecordRepository) (Map, error) {
	mapping := make(Map)
	s := cv.Sanitizer()

	for {
	looper:
		record, err := n.Reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if record[0] == "" {
			continue
		}

		rr := &types.Record{
			ModuleID:    n.Module.ID,
			NamespaceID: n.Namespace.ID,
			CreatedAt:   time.Now(),
		}

		recordValues := types.RecordValueSet{}

		// assure a valid legacy reference
		recordValues = append(recordValues, &types.RecordValue{
			Name:    LegacyRefIDField,
			Value:   record[0],
			Place:   0,
			Updated: true,
		})

		// convert the given row into a { field: value } map; this will be used
		// for expression evaluation
		row := map[string]string{}
		for i, h := range n.Header {
			row[h] = record[i]
		}

		for i, h := range n.Header {
			// will contain string values for the given field
			var values []string
			val := record[i]

			// system values should be kept on the record's root level
			if isSysField(h) {
				switch strings.ToLower(h) {
				case "ownerid":
					rr.OwnedBy = users[val]
					break

					// ignore deleted values, as SF provides minimal info about those
				case "isdeleted":
					if val == "1" || strings.ToLower(val) == "true" {
						goto looper
					}
					break

				case "createddate":
					if val != "" {
						rr.CreatedAt, err = time.Parse(SfDateTimeLayout, val)
						if err != nil {
							return nil, err
						}
					}
					break

				case "createdbyid":
					rr.CreatedBy = users[val]
					break

				case "lastmodifiedbyid":
					rr.UpdatedBy = users[val]
					break

				case "lastmodifieddate":
					if val != "" {
						tt, err := time.Parse(SfDateTimeLayout, val)
						rr.UpdatedAt = &tt
						if err != nil {
							return nil, err
						}
					}
					break
				}
			} else {
				// other user defined values should be kept inside `values`
				joined := ""
				if strings.Contains(h, ":") {
					pts := strings.Split(h, ":")
					h = pts[0]
					joined = pts[1]
				}

				// find corresponding field
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

				// temp set of raw values that should be processed further.
				// this gives us support for multi value fields when joining a sources
				rawValues := make([]string, 0)
				if joined != "" {
					tmp := n.FieldMap[val]
					for _, e := range tmp {
						rawValues = append(rawValues, e[joined])
					}
				} else {
					rawValues = []string{val}
				}

				for _, val := range rawValues {
					// handle references. Spliced nodes should not perform this step, since
					// they can't rely on any dependency. This is corrected with `correctRecordRefs`
					if f.Options["moduleID"] != nil {
						if !n.isSpliced {
							ref, ok := f.Options["moduleID"].(string)
							if !ok {
								return nil, errors.New("moduleField.record.invalidRefFormat")
							}

							fetch := false

							if val == "" {
								continue
							}

							if mod, ok := n.idMap[ref]; !ok {
								fetch = true
							} else if v, ok := mod[val]; !ok || v == "" {
								fetch = true
							} else {
								val = v
							}

							if fetch {
								val, err = n.fetchRemoteRef(val, ref, repo)
								if err != nil {
									continue
								}
							}

							if val == "" {
								continue
							}
						}
						values = append(values, val)
					} else if f.Kind == "User" {
						// handle user references
						if u, ok := users[val]; ok {
							val = fmt.Sprint(u)
						} else {
							continue
						}
						values = append(values, val)
					} else {
						// generic value handling
						val = strings.Map(fixUtf, val)
						if val == "" {
							continue
						}

						values = append(values, val)
					}
				}

				// value post-proc & record value creation
				for i, v := range values {
					v, err = n.mapValue(h, v, row)
					if err != nil {
						return nil, err
					}
					rv := &types.RecordValue{
						Name:    h,
						Value:   v,
						Place:   uint(i),
						Updated: true,
					}
					// ref values of spliced nodes should get updated later
					if n.isSpliced && f.IsRef() {
						rv.Updated = false
					}
					recordValues = append(recordValues, rv)
				}

				recordValues = s.Run(n.Module, recordValues)
			}
		}

		// create record
		r, err := repo.Create(rr)
		if err != nil {
			return nil, err
		}

		// update record values with recordID
		for _, v := range recordValues {
			v.RecordID = r.ID
		}

		if !n.isSpliced {
			err = repo.UpdateValues(r.ID, recordValues)
			if err != nil {
				log.Printf("[issue] db.UpdateValues | %d | %s | %s \n", r.ID, recordValues.String(), err.Error())
				// return nil, err
			}
		}

		// spliced nodes should preserve their records for later ref processing
		if n.isSpliced {
			rr.Values = recordValues
			n.original.records = append(n.original.records, rr)
		}

		// update mapping map
		mapping[record[0]] = fmt.Sprint(rr.ID)
	}

	return mapping, nil
}

func (n *ImportNode) mapValue(field, val string, row map[string]string) (string, error) {
	if fmp, ok := n.ValueMap[field]; ok {
		nvl := ""
		if mpv, ok := fmp[val]; ok {
			nvl = mpv
		} else if mpv, ok := fmp["*"]; ok {
			nvl = mpv
		}

		// expression evaluation
		if nvl != "" && strings.HasPrefix(nvl, EvalPrefix) {
			opp := nvl[len(EvalPrefix):len(nvl)]
			ev, err := ExprLang.NewEvaluable(opp)
			if err != nil {
				return "", err
			}

			val, err = ev.EvalString(context.Background(), map[string]interface{}{"cell": val, "row": row})
			if err != nil {
				return "", err
			}
		} else if nvl != "" {
			val = nvl
		}
	}
	return val, nil
}

// helper to assure correct date time formatting
func assureDateFormat(val string, opt types.ModuleFieldOptions) (string, error) {
	pvl, err := time.Parse(SfDateTimeLayout, val)
	if err != nil {
		return "", err
	}

	if opt.Bool("onlyDate") {
		val = pvl.Format(DateOnlyLayout)
	} else if opt.Bool("onlyTime") {
		val = pvl.Format(TimeOnlyLayout)
	} else {
		val = pvl.Format(time.RFC3339)
	}
	return val, nil
}
