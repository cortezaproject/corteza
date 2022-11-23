package reporting

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/modern-go/reflect2"
)

type (
	// reportFrameBuilder assist in frame construction from iterators
	//
	// Primarily simplifies mapping the correct iterator attributes to correct
	// frame row columns and having them encoded properly.
	reportFrameBuilder struct {
		def    *FrameDefinition
		frame  *Frame
		ownCol string

		attrMapping map[string]int
		attrMvDel   map[string]string
	}
)

// newReportFrameBuilder initializes a new report frame builder
func newReportFrameBuilder(def *FrameDefinition) *reportFrameBuilder {
	// Index cols for easier lookups
	attrMap := make(map[string]int)
	mvDelMap := make(map[string]string)
	for i, c := range def.Columns {
		attrMap[c.Name] = i

		if c.Multivalue {
			mvDelMap[c.Name] = c.MultivalueDelimiter
			if mvDelMap[c.Name] == "" {
				mvDelMap[c.Name] = "\n"
			}
		}
	}

	out := &reportFrameBuilder{
		def:         def,
		attrMapping: attrMap,
	}

	// Init output frame
	out.freshFrame()
	return out
}

// linked includes additional metadata required by the link step
func (b *reportFrameBuilder) linked(ownCol, relCol, relSrc string) {
	b.ownCol = ownCol
	b.frame.RelColumn = relCol
	b.frame.RelSource = relSrc
}

// addRow adds a new dal.Row to the frame
func (b *reportFrameBuilder) addRow(r *dal.Row) {
	rrow := make(FrameRow, len(b.def.Columns))

	for k, cc := range r.CountValues() {
		ix, ok := b.attrMapping[k]
		if !ok {
			continue
		}

		auxRow := make(FrameRow, cc)
		for i := uint(0); i < cc; i++ {
			val, _ := r.GetValue(k, i)
			auxRow[i] = b.stringifyVal(k, val)
		}
		rrow[ix] = b.joinMultiVal(k, auxRow)
	}

	b.frame.Rows = append(b.frame.Rows, rrow)

	if b.ownCol != "" {
		v, _ := r.GetValue(b.ownCol, 0)
		b.frame.RefValue = b.stringifyVal(b.ownCol, v)
	}
}

// done returns the constructed frame and prepares a new frame with the same
// metadata as the original one
func (b *reportFrameBuilder) done() *Frame {
	out := b.frame
	b.freshFrame()

	return out
}

func (b *reportFrameBuilder) stringifyVal(col string, val any) string {
	// Edgecase for when values are nil since the stringer uses <nil> for those
	if reflect2.IsNil(val) {
		return ""
	}

	// @todo nicer formatting and such? V1 didn't do much different
	return fmt.Sprintf("%v", val)
}

func (b *reportFrameBuilder) joinMultiVal(col string, vals []string) string {
	// @todo b.attrMvDel[col] is never actually set so it should be ALONG with this default
	del, ok := b.attrMvDel[col]
	if !ok {
		del = "\n"
	}
	return strings.Join(vals, del)
}

func (b *reportFrameBuilder) freshFrame() {
	// reuse the old frame metadata and clears out the rows
	if b.frame != nil {
		aux := *b.frame
		b.frame = &aux
		b.frame.Rows = []FrameRow{}
		return
	}

	b.frame = &Frame{
		Name:   b.def.Name,
		Source: b.def.Source,
		Ref:    b.def.Ref,

		Rows: []FrameRow{},

		Columns: b.def.Columns,
		Sort:    b.def.Sort,
		Filter:  b.def.Filter,
	}

	if b.def.Paging != nil {
		b.frame.Paging = &filter.Paging{
			Limit:      b.def.Paging.Limit,
			PageCursor: b.def.Paging.PageCursor,
		}
	}
}
